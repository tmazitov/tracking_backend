DROP TABLE IF EXISTS users CASCADE; 

CREATE TABLE users (
    id			BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    
	role		INT			NOT NULL DEFAULT 0,
    email		VARCHAR(64)	NOT NULL,
	short_name	VARCHAR(20),
	telephone	VARCHAR(12),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    edited_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(id)
);

DROP TABLE IF EXISTS users_holidays CASCADE; 

CREATE TABLE users_holidays (
	worker_id			BIGINT		NOT NULL,
	author_id			BIGINT		NOT NULL,
	start_at			TIMESTAMP	NOT NULL,
	CONSTRAINT fk_worker FOREIGN KEY (worker_id) REFERENCES users(id),
	CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id)
);

DROP TABLE IF EXISTS users_job CASCADE; 

CREATE TABLE users_job (
	id					BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
	user_id				BIGINT		NOT NULL,
	job_type			INT8		NOT NULL,
	job_experience		INT8		NOT NULL,
	job_mail			VARCHAR(512) NOT NULL,	

	created_at			TIMESTAMP	NOT NULL DEFAULT NOW(),
	accepted_at			TIMESTAMP			 DEFAULT NULL,
	finished_at			TIMESTAMP			 DEFAULT NULL,
	canceled_at			TIMESTAMP			 DEFAULT NULL,
	CONSTRAINT fk_worker FOREIGN KEY (user_id) REFERENCES users(id)
);