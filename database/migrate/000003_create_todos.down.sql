-- 删除trigger
DROP TRIGGER IF EXISTS trigger_update_todos_status_timestamp ON todos;
DROP TRIGGER IF EXISTS trigger_update_updated_at_todos ON todos;

-- 删除函数
DROP FUNCTION IF EXISTS trigger_todos_status_timestamp();

-- 删除索引
DROP INDEX IF EXISTS todos_reporter_id_status;
DROP INDEX IF EXISTS todos_assignee_id_status;

-- 删除表
DROP TABLE IF EXISTS todos;