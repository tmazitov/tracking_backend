DROP TABLE IF EXISTS orders CASCADE; 

CREATE TABLE orders (
    id              BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    
    startAt         DATE NOT NULL,
    endAt           DATE NOT NULL,

    points          BIGINT[] NOT NULL,
    helpers         INT NOT NULL DEFAULT 0,
    comment_message VARCHAR(256) NOT NULL DEFAULT '', 
    PRIMARY KEY(id)
);