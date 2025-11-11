DROP TRIGGER IF EXISTS trigger_update_updated_at_projects ON projects;
DROP TRIGGER IF EXISTS trigger_update_updated_at_processes ON processes;

DROP INDEX IF EXISTS idx_processes_project_id;

DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS processes;

DROP TYPE IF EXISTS project_status;
DROP TYPE IF EXISTS process_status;