ALTER TABLE tasks 
    ADD COLUMN IF NOT EXISTS frequency JSONB,
    ADD COLUMN IF NOT EXISTS next_run_time TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00';


COMMENT ON COLUMN tasks.frequency IS 'JSONB поле с правилами периодичности задачи.';
COMMENT ON COLUMN tasks.next_run_time IS 'Время следующего запуска задачи';

CREATE INDEX IF NOT EXISTS idx_tasks_frequency ON tasks USING GIN(frequency);
CREATE INDEX IF NOT EXISTS idx_tasks_next_run_time ON tasks(next_run_time);