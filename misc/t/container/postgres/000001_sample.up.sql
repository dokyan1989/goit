CREATE TABLE IF NOT EXISTS foos (
  foo_id SERIAL PRIMARY KEY,
  foo_1 VARCHAR(100) NOT NULL,
  foo_2 TEXT ARRAY NULL,
  foo_3 JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bars (
  bar_id SERIAL PRIMARY KEY,
  foo_id BIGINT NOT NULL,
  bar_1 VARCHAR(50) NULL,
  bar_2 POINT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE bars 
ADD CONSTRAINT fk_bars_to_foos 
FOREIGN KEY (foo_id) 
REFERENCES foos (foo_id) ON DELETE CASCADE;