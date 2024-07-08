CREATE TABLE IF NOT EXISTS permission (
  permission_id SERIAL PRIMARY KEY,
  permission_name VARCHAR(200) NOT NULL,
  description VARCHAR(500) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  status VARCHAR(50) DEFAULT 'NEW'
);
COMMENT ON COLUMN permission.permission_id IS 'A unique identifier for each permission (Primary Key).';
COMMENT ON COLUMN permission.permission_name IS 'The name of the specific permission.';
COMMENT ON COLUMN permission.description IS 'Explains the permission functionality.';
COMMENT ON COLUMN permission.created_at IS 'When the permission is created.';
COMMENT ON COLUMN permission.updated_at IS 'When the permission is updated.';
COMMENT ON COLUMN permission.status IS 'Indicates whether the permission is active.';
