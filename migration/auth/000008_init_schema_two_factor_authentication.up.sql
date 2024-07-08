CREATE TABLE IF NOT EXISTS two_factor_authentication (
  user_id INTEGER NOT NULL,
  method VARCHAR(50) NOT NULL,
  verification_code VARCHAR(500) NOT NULL,
  expiry_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE two_factor_authentication ADD FOREIGN KEY (user_id) REFERENCES "user" (user_id);
COMMENT ON COLUMN two_factor_authentication.user_id IS 'Refers to the user table.';
COMMENT ON COLUMN two_factor_authentication.method IS 'How the second authentication step is carried out, e.g., SMS, Authenticator App.';
COMMENT ON COLUMN two_factor_authentication.verification_code IS 'The actual verification code.';
COMMENT ON COLUMN two_factor_authentication.expiry_date IS 'When the code becomes invalid.';