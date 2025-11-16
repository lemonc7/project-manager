ALTER TABLE processes
DROP CONSTRAINT unique_project_sequence;

ALTER TABLE processes
ADD CONSTRAINT unique_project_sequence UNIQUE (project_id, sequence);