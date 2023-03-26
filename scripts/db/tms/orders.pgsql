DROP TABLE IF EXISTS orders CASCADE; 

CREATE TABLE orders (
    id                  BIGINT		GENERATED ALWAYS AS IDENTITY UNIQUE,
    
	owner_id            BIGINT		NOT NULL,
	manager_id			BIGINT,

    start_at            TIMESTAMP	NOT NULL,
    end_at              TIMESTAMP 	DEFAULT NULL,

    status_id			INT			NOT NULL DEFAULT 1,
    points              BIGINT[] 	NOT NULL,
    helpers             INT 		NOT NULL DEFAULT 0,
    comment_message     VARCHAR(256) NOT NULL DEFAULT '', 
    is_fragile_cargo    BOOLEAN 	NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP    		NOT NULL DEFAULT NOW(),
    edited_at TIMESTAMP     		NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_owner FOREIGN KEY(owner_id) REFERENCES users(id),
    CONSTRAINT fk_manager FOREIGN KEY(manager_id) REFERENCES users(id),
    PRIMARY KEY(id)
);