-- todos表
CREATE TABLE todos (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  remark TEXT NOT NULL DEFAULT '',
  priority INT NOT NULL DEFAULT 1,
  status VARCHAR(10) NOT NULL DEFAULT 'todo',
  reporter_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  assignee_id INT NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  CONSTRAINT check_priority_range CHECK (priority BETWEEN 1 AND 5),
  CONSTRAINT check_status CHECK (status IN ('todo', 'doing', 'done'))
);

-- 索引
CREATE INDEX todos_reporter_id_status ON todos(reporter_id,status);
CREATE INDEX todos_assignee_id_status ON todos(assignee_id,status);

-- todos状态更新时间触发器
CREATE OR REPLACE FUNCTION trigger_todos_status_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  -- status第一次变成"doing",设置started_at
  IF NEW.status = 'doing'
    AND OLD.status <> 'doing' -- NOT NULL，不用判断NULL
    AND NEW.started_at IS NULL THEN
    NEW.started_at = CURRENT_TIMESTAMP;
  END IF;

  -- status变成"done",设置finished_at
  IF NEW.status = 'done'
    AND OLD.status <> 'done' THEN
    NEW.finished_at = CURRENT_TIMESTAMP;
  -- status从"done"回退,清空finished_at
  ELSIF OLD.status = 'done'
    AND NEW.status <> 'done' THEN
    NEW.finished_at = NULL;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 应用到todos表
CREATE TRIGGER trigger_update_todos_status_timestamp
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE FUNCTION trigger_todos_status_timestamp();

-- 自动更新updated_at
CREATE TRIGGER trigger_update_updated_at_todos
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();