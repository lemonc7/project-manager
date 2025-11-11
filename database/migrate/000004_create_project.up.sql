-- 创建枚举类型
CREATE TYPE project_status AS ENUM (
  '排单',
  '进行',
  '暂停',
  '完成', 
  '终止',
  '归档'
);

CREATE TYPE process_status as ENUM (
  '等待',
  '进行',
  '完成',
  '返工',
  '终止',
  '取消',
  '暂停'
);

-- 项目表
CREATE TABLE projects (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  status project_status NOT NULL DEFAULT '排单',
  manager_id INT REFERENCES employees(id) ON DELETE SET NULL,
  started_at TIMESTAMPTZ,
  duration_days INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 工艺流程表
CREATE TABLE processes (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  status process_status NOT NULL DEFAULT '等待',
  employee_id INT REFERENCES employees(id) ON DELETE SET NULL,
  sequence INT NOT NULL,
  start_time TIMESTAMPTZ,
  end_time TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_project_sequence UNIQUE(project_id, sequence)
);

-- 应用updated_at触发器
CREATE TRIGGER trigger_update_updated_at_projects
BEFORE UPDATE ON projects
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TRIGGER trigger_update_updated_at_processes
BEFORE UPDATE ON processes
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

-- 建立索引
CREATE INDEX IF NOT EXISTS idx_processes_project_id ON processes(project_id);