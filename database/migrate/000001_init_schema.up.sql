-- 创建状态表
CREATE TABLE IF NOT EXISTS project_status (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  sort_order INT NOT NULL
);

CREATE TABLE process_status (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  sort_order INT NOT NULL
);

CREATE TABLE todo_status (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  sort_order INT NOT NULL
);

CREATE TABLE roles (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT ''
);

INSERT INTO project_status (name, sort_order) VALUES
('进行', 10),
('排单', 20),
('暂停', 30),
('完成', 40),
('终止', 100);

INSERT INTO process_status (name, sort_order) VALUES
('进行', 10),
('准备', 20),
('等待', 30),
('暂停', 40),
('返工', 50),
('完成', 60),
('终止', 100);

INSERT INTO todo_status (name, sort_order) VALUES
('进行', 10),
('待办', 20),
('完成', 30),
('终止', 40);

INSERT INTO roles (name, description) VALUES
('admin', '管理员'),
('manager', '项目经理'),
('employee', '员工'),
('viewer', '用户');

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
  name VARCHAR(20) NOT NULL,
  username VARCHAR(20) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  team_id INT REFERENCES teams(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE employee_roles (
  employee_id INT REFERENCES employees(id) ON DELETE CASCADE,
  role_id INT REFERENCES roles(id) ON DELETE CASCADE,
  PRIMARY KEY (employee_id, role_id)
);

-- 项目表
CREATE TABLE projects (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  remark TEXT NOT NULL DEFAULT '',
  status_id INT NOT NULL REFERENCES project_status(id),
  manager_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  priority INT NOT NULL DEFAULT 1,
  started_at TIMESTAMPTZ,
  duration_days INT,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT check_duration_days_positive CHECK (duration_days >= 0),
  CONSTRAINT check_projects_priority_range CHECK (priority BETWEEN 1 AND 5)
);

-- todos表
CREATE TABLE todos (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  remark TEXT NOT NULL DEFAULT '',
  priority INT NOT NULL DEFAULT 1,
  status_id INT NOT NULL REFERENCES todo_status(id),
  reporter_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  assignee_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  CONSTRAINT check_todos_priority_range CHECK (priority BETWEEN 1 AND 5)
);

-- 工艺流程表
CREATE TABLE processes (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  remark TEXT NOT NULL DEFAULT '',
  status_id INT NOT NULL REFERENCES process_status(id),
  employee_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  step_order INT NOT NULL,
  spend_hours NUMERIC(10, 2),
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_project_step_order UNIQUE(project_id, step_order)
  DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT check_step_order_positive CHECK (step_order > 0),
  CONSTRAINT check_spend_hours_positive CHECK (spend_hours >= 0)
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
  ON projects(manager_id, status_id, started_at);
CREATE INDEX idx_projects_priority_started
  ON projects(priority, started_at);

-- processes
CREATE INDEX idx_processes_project_step_order
  ON processes(project_id, step_order);
CREATE INDEX idx_processes_employee_status
  ON processes(employee_id, status_id);

-- todos
CREATE INDEX todos_reporter_id_status
  ON todos(reporter_id, status_id);
CREATE INDEX todos_assignee_id_status
  ON todos(assignee_id, status_id);

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

-- 应用到todos表
CREATE TRIGGER trigger_update_updated_at_todos
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
