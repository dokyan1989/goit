CREATE TABLE IF NOT EXISTS role_permission (
  role_permission_id SERIAL PRIMARY KEY,
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL
);
ALTER TABLE role_permission ADD FOREIGN KEY (role_id) REFERENCES role (role_id);
ALTER TABLE role_permission ADD FOREIGN KEY (permission_id) REFERENCES permission (permission_id);
COMMENT ON COLUMN role_permission.role_permission_id IS 'A unique identifier for each role-permission pairing (Primary Key).';
COMMENT ON COLUMN role_permission.role_id IS 'Refers to the role table.';
COMMENT ON COLUMN role_permission.permission_id IS 'Refers to the permission table.';