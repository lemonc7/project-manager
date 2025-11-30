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
  '完成',
  '暂停',
  '终止'
);

-- 小组表
CREATE TABLE teams (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  description TEXT,
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
  description TEXT,
  remark TEXT,
  status project_status NOT NULL DEFAULT '排单',
  manager_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  priority INT NOT NULL DEFAULT 5,
  started_at TIMESTAMPTZ,
  duration INT,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_duration_positive CHECK (duration > 0),
  CONSTRAINT check_priority_range CHECK (priority BETWEEN 1 AND 9)
);

-- 工艺流程表
CREATE TABLE processes (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  remark TEXT,
  status process_status NOT NULL DEFAULT '等待',
  employee_id INT REFERENCES employees(id) ON DELETE SET NULL,
  sequence INT NOT NULL,
  spend_time DOUBLE PRECISION,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_project_sequence UNIQUE(project_id, sequence)
  DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT check_sequence_positive CHECK (sequence > 0),
  CONSTRAINT check_spend_time_positive CHECK (spend_time > 0)
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