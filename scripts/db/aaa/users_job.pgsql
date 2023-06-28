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
)