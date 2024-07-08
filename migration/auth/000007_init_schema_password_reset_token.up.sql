CREATE TABLE IF NOT EXISTS password_reset_token (
  password_reset_token_id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  token_value VARCHAR(500) NOT NULL,
  expiry_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE password_reset_token ADD FOREIGN KEY (user_id) REFERENCES "user" (user_id);
COMMENT ON COLUMN password_reset_token.password_reset_token_id IS 'Unique identifier for each token (Primary Key).';
COMMENT ON COLUMN password_reset_token.user_id IS 'Refers to the user table.';
COMMENT ON COLUMN password_reset_token.token_value IS 'The actual token sent to user emails.';
COMMENT ON COLUMN password_reset_token.expiry_date IS 'Token expiration time.';