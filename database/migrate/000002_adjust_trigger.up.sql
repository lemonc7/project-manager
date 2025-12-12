-- 删除status trigger
DROP TRIGGER IF EXISTS trigger_update_next_processes_status ON processes;
DROP TRIGGER IF EXISTS trigger_start_project_processes ON projects;
-- 删除函数
DROP FUNCTION IF EXISTS update_next_processes_status();
DROP FUNCTION IF EXISTS start_project_processes();

-- 创建processes状态更新时间触发器
CREATE OR REPLACE FUNCTION trigger_processes_status_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  -- status第一次变成"进行"，设置started_at
  IF NEW.status = '进行'
    AND OLD.status <> '进行' -- NOT NULL，不用判断NULL
    AND NEW.started_at IS NULL THEN
    NEW.started_at = CURRENT_TIMESTAMP;
  END IF;

  -- status变成"完成"或"终止"，设置finished_at
  IF NEW.status IN ('完成', '终止')
    AND OLD.status NOT IN ('完成', '终止') THEN
    NEW.finished_at = CURRENT_TIMESTAMP;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建projects状态更新时间触发器
CREATE OR REPLACE FUNCTION trigger_projects_status_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  -- status第一次变成"进行",设置started_at
  IF NEW.status = '进行'
    AND OLD.status <> '进行' -- NOT NULL，不用判断NULL
    AND NEW.started_at IS NULL THEN
    NEW.started_at = CURRENT_TIMESTAMP;
  END IF;

  -- status变成"完成"或"终止",设置finished_at
  IF NEW.status IN ('完成', '终止')
    AND OLD.status NOT IN ('完成', '终止') THEN
    NEW.finished_at = CURRENT_TIMESTAMP;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 应用到processes表
CREATE TRIGGER trigger_update_processes_status_timestamp
BEFORE UPDATE ON processes
FOR EACH ROW
EXECUTE FUNCTION trigger_processes_status_timestamp();

-- 应用到projects表
CREATE TRIGGER trigger_update_projects_status_timestamp
BEFORE UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION trigger_projects_status_timestamp();



