UPDATE teams SET description = '' WHERE description IS NULL;

-- 修改teams表
ALTER TABLE teams
  ALTER COLUMN description TYPE TEXT,
  ALTER COLUMN description SET NOT NULL,
  ALTER COLUMN description SET DEFAULT '',
  ALTER COLUMN name TYPE VARCHAR(20),
  ALTER COLUMN name SET NOT NULL;

ALTER TABLE teams
  ADD CONSTRAINT unique_teams_name UNIQUE (name);

-- 修改employees表
ALTER TABLE employees
  ALTER COLUMN name TYPE VARCHAR(10),
  ALTER COLUMN name SET NOT NULL;
