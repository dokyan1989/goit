-- https://www.linkedin.com/pulse/database-structure-user-role-management-aj-february
CREATE TABLE IF NOT EXISTS "user" (
  user_id SERIAL PRIMARY KEY,
  user_name VARCHAR(100) NOT NULL,
  password VARCHAR(80) NOT NULL,
  email VARCHAR(100) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_login_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  status VARCHAR(50) NOT NULL
);
ALTER TABLE "user" ADD UNIQUE (user_name);
COMMENT ON COLUMN "user".user_id IS 'A unique identifier for each user (Primary Key).';
COMMENT ON COLUMN "user".user_name IS 'The chosen username.';
COMMENT ON COLUMN "user".password IS 'This should always be hashed and salted, never stored in plain text for security reasons.';
COMMENT ON COLUMN "user".email IS 'User email address.';
COMMENT ON COLUMN "user".first_name IS 'User first name.';
COMMENT ON COLUMN "user".last_name IS 'User last name.';
COMMENT ON COLUMN "user".created_at IS 'When the user first joined.';
COMMENT ON COLUMN "user".updated_at IS 'When the user is updated.';
COMMENT ON COLUMN "user".last_login_at IS 'The last date and time the user accessed the platform.';
COMMENT ON COLUMN "user".status IS 'Indicates whether the user is active.';
