DROP TABLE IF EXISTS users_holidays CASCADE; 

CREATE TABLE users_holidays (
	worker_id			BIGINT		NOT NULL,
	author_id			BIGINT		NOT NULL,
	start_at			TIMESTAMP	NOT NULL,
	CONSTRAINT fk_worker FOREIGN KEY (worker_id) REFERENCES users(id),
	CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id)
);