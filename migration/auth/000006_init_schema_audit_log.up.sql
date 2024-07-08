CREATE TABLE IF NOT EXISTS audit_log (
  audit_log_id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  activity_type VARCHAR(500) NOT NULL,
  description VARCHAR(500) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE audit_log ADD FOREIGN KEY (user_id) REFERENCES "user" (user_id);
COMMENT ON COLUMN audit_log.audit_log_id IS 'A unique identifier for each log (Primary Key).';
COMMENT ON COLUMN audit_log.user_id IS 'Refers to the user table.';
COMMENT ON COLUMN audit_log.activity_type IS 'Type of activity, e.g., Role Assignment, Permission Change.';
COMMENT ON COLUMN audit_log.description IS 'Detailed record of the activity.';
COMMENT ON COLUMN audit_log.created_at IS 'When the activity occurred.';