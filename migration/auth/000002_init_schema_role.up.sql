CREATE TABLE IF NOT EXISTS role (
  role_id SERIAL PRIMARY KEY,
  role_name VARCHAR(200) NOT NULL,
  description VARCHAR(500) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  status VARCHAR(50) DEFAULT 'NEW'
);
COMMENT ON COLUMN role.role_id IS 'A unique identifier for each role (Primary Key).';
COMMENT ON COLUMN role.role_name IS 'Name of the role.';
COMMENT ON COLUMN role.description IS 'A brief about the role’s functionalities.';
COMMENT ON COLUMN role.created_at IS 'When the role is created.';
COMMENT ON COLUMN role.updated_at IS 'When the role is updated.';
COMMENT ON COLUMN role.status IS 'Indicates whether the role is active.';