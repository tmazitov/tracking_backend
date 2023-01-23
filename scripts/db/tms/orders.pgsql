DROP TABLE IF EXISTS orders CASCADE; 

CREATE TABLE orders (
    id              BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    
    startAt         TIMESTAMP NOT NULL,
    endAt           TIMESTAMP,

    points          BIGINT[] NOT NULL,
    helpers         INT NOT NULL DEFAULT 0,
    comment_message VARCHAR(256) NOT NULL DEFAULT '', 
    is_fragile_cargo BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY(id)
);