-- 删除视图
DROP VIEW IF EXISTS projects_view;

-- 删除触发器
DROP TRIGGER IF EXISTS trigger_update_updated_at_todos ON todos;
DROP TRIGGER IF EXISTS trigger_update_updated_at_system_settings ON system_settings;
DROP TRIGGER IF EXISTS trigger_update_updated_at_processes ON processes;
DROP TRIGGER IF EXISTS trigger_update_updated_at_projects ON projects;
DROP TRIGGER IF EXISTS trigger_update_updated_at_employees ON employees;
DROP TRIGGER IF EXISTS trigger_update_updated_at_teams ON teams;

-- 删除函数
DROP FUNCTION IF EXISTS update_updated_at();

-- 删除索引
DROP INDEX IF EXISTS todos_status_id;
DROP INDEX IF EXISTS todos_assignee_id_status;
DROP INDEX IF EXISTS todos_reporter_id_status;
DROP INDEX IF EXISTS idx_processes_status_id;
DROP INDEX IF EXISTS idx_processes_employee_status;
DROP INDEX IF EXISTS idx_processes_project_step_order;
DROP INDEX IF EXISTS idx_projects_status_id;
DROP INDEX IF EXISTS idx_projects_priority_started;
DROP INDEX IF EXISTS idx_projects_manager_status_started;
DROP INDEX IF EXISTS idx_projects_manager_id;

-- 删除表（按依赖顺序）
DROP TABLE IF EXISTS processes;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS employee_roles;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS todo_status;
DROP TABLE IF EXISTS process_status;
DROP TABLE IF EXISTS project_status;

-- 删除枚举类型
DROP TYPE IF EXISTS status_category;
