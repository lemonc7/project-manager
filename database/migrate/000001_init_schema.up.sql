-- 创建枚举类型
CREATE TYPE project_status AS ENUM (
  '排单',
  '进行',
  '完成', 
  '暂停',
  '终止'
);

CREATE TYPE process_status as ENUM (
  '等待',
  '准备',
  '进行',
  '返工',
  '完成',
  '暂停',
  '终止'
);

-- 小组表
CREATE TABLE teams (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 员工表
CREATE TABLE employees (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(10) NOT NULL,
  username VARCHAR(20) NOT NULL UNIQUE,
  password_hash VARCHAR(60) NOT NULL,
  roles TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
  team_id INT REFERENCES teams(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- 项目表
CREATE TABLE projects (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  remark TEXT NOT NULL DEFAULT '',
  status project_status NOT NULL DEFAULT '排单',
  manager_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  priority INT NOT NULL DEFAULT 1,
  started_at TIMESTAMPTZ,
  duration INT,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_duration_positive CHECK (duration >= 0),
  CONSTRAINT check_priority_range CHECK (priority BETWEEN 1 AND 5)
);

-- 项目视图，用来计算剩余天数
CREATE OR REPLACE VIEW projects_view AS
SELECT
  p.id,
  p.name,
  p.description,
  p.remark,
  p.status,
  p.manager_id,
  p.priority,
  p.started_at,
  p.duration,
  p.finished_at,
  p.created_at,
  p.updated_at,
  -- 剩余天数 (必须用于排序，保留在数据库计算)
  (
    CASE 
      WHEN 
        p.status NOT IN ('完成', '终止') 
        AND p.started_at IS NOT NULL 
        AND p.duration IS NOT NULL THEN
        CEIL(
            EXTRACT(EPOCH FROM ((p.started_at + make_interval(days => p.duration)) - NOW())) / 86400.0
        )::INTEGER
      ELSE NULL
    END
  ) AS remaining_days
FROM 
  projects p;

-- 工艺流程表
CREATE TABLE processes (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  remark TEXT NOT NULL DEFAULT '',
  status process_status NOT NULL DEFAULT '等待',
  employee_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  sequence INT NOT NULL,
  spend_time DOUBLE PRECISION,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_project_sequence UNIQUE(project_id, sequence)
  DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT check_sequence_positive CHECK (sequence > 0),
  CONSTRAINT check_spend_time_positive CHECK (spend_time >= 0)
);

-- 系统设置表
CREATE TABLE system_settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
-- projects
CREATE INDEX idx_projects_manager_id ON projects(manager_id);
CREATE INDEX idx_projects_manager_status_started
  ON projects(manager_id, status, started_at);
CREATE INDEX idx_projects_priority_started
  ON projects(priority, started_at);

-- processes
CREATE INDEX idx_processes_project_sequence
  ON processes(project_id, sequence);
CREATE INDEX idx_processes_employee_status
  ON processes(employee_id, status);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 应用到teams表
CREATE TRIGGER trigger_update_updated_at_teams
BEFORE UPDATE ON teams
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- 应用到employees表
CREATE TRIGGER trigger_update_updated_at_employees
BEFORE UPDATE ON employees
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- 应用到projects表
CREATE TRIGGER trigger_update_updated_at_projects
BEFORE UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- 应用到processes表
CREATE TRIGGER trigger_update_updated_at_processes
BEFORE UPDATE ON processes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- 应用到system_settings表
CREATE TRIGGER trigger_update_updated_at_system_settings
BEFORE UPDATE ON system_settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();


-- 创建工艺流程状态流转监听器
CREATE OR REPLACE FUNCTION update_next_processes_status()
RETURNS TRIGGER AS $$
DECLARE
  current_sequence INT;
  next_sequence INT;
  next_next_sequence INT;
BEGIN
  -- 确保新的状态是完成，旧的状态不是完成
  IF NEW.status = '完成' AND OLD.status <> '完成' THEN
    current_sequence := NEW.sequence;
    next_sequence := current_sequence + 1;
    next_next_sequence := current_sequence + 2;
    -- 更新下一步工艺为进行，确认下一步工艺是准备状态
    UPDATE processes
    SET 
      status = '进行',
      started_at = CURRENT_TIMESTAMP
    WHERE 
      project_id = NEW.project_id 
      AND sequence = next_sequence
      AND status = '准备';

    -- 更新下下步工艺为准备，确认工艺是等待状态
    UPDATE processes
    SET status = '准备'
    WHERE 
      project_id = NEW.project_id 
      AND sequence = next_next_sequence
      AND status = '等待';
    
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建项目开始状态监听器
CREATE OR REPLACE FUNCTION start_project_processes()
RETURNS TRIGGER AS $$
BEGIN
  -- 确保项目新的状态是进行，旧的状态不是进行
  IF NEW.status = '进行' AND OLD.status <> '进行' THEN
    -- 更新第一步工艺为进行，确认工艺是等待状态
    UPDATE processes
    SET 
      status = '进行',
      started_at = CURRENT_TIMESTAMP
    WHERE 
      project_id = NEW.id 
      AND sequence = 1
      AND status = '等待';
    
    -- 更新第二步工艺为准备，确认工艺是等待状态
    UPDATE processes
    SET status = '准备'
    WHERE 
      project_id = NEW.id 
      AND sequence = 2
      AND status = '等待';

  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 应用到processes表
CREATE TRIGGER trigger_update_next_processes_status
AFTER UPDATE ON processes
FOR EACH ROW
EXECUTE FUNCTION update_next_processes_status();

-- 应用到projects表
CREATE TRIGGER trigger_start_project_processes
AFTER UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION start_project_processes();