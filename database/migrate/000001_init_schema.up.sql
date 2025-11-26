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
  started_at TIMESTAMPTZ,
  duration INT,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_duration_positive CHECK (duration > 0)
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
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_project_sequence UNIQUE(project_id, sequence)
  DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT check_sequence_positive CHECK (sequence > 0)
);

-- 系统设置表
CREATE TABLE system_settings (
  key TEXT PRIMARY KEY,
  value TEXT,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
-- projects
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_manager_id ON projects(manager_id);
CREATE INDEX idx_projects_started_at ON projects(started_at);

-- processes
CREATE INDEX idx_processes_project_id ON processes(project_id);
CREATE INDEX idx_processes_status ON processes(status);
CREATE INDEX idx_processes_employee_id ON processes(employee_id);

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