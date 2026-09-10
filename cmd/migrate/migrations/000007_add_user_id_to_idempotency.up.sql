ALTER TABLE idempotency_keys ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id);
