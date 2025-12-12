-- 删除status trigger
DROP TRIGGER IF EXISTS trigger_update_processes_status_timestamp ON processes;
DROP TRIGGER IF EXISTS trigger_update_projects_status_timestamp ON projects;
-- 删除函数
DROP FUNCTION IF EXISTS trigger_processes_status_timestamp();
DROP FUNCTION IF EXISTS trigger_projects_status_timestamp();