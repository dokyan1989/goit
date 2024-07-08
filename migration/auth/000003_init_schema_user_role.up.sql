CREATE TABLE IF NOT EXISTS user_role (
  user_role_id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL
);
ALTER TABLE user_role ADD FOREIGN KEY (user_id) REFERENCES "user" (user_id);
ALTER TABLE user_role ADD FOREIGN KEY (role_id) REFERENCES role (role_id);
COMMENT ON COLUMN user_role.user_role_id IS 'A unique identifier for each user-role combination (Primary Key).';
COMMENT ON COLUMN user_role.user_id IS 'Refers to the user table.';
COMMENT ON COLUMN user_role.role_id IS 'Refers to the role table.';