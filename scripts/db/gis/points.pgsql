DROP TABLE IF EXISTS points CASCADE; 

CREATE TABLE points (
    id  BIGINT GENERATED ALWAYS AS IDENTITY UNIQUE,
    floor INT8 NOT NULL DEFAULT 1,
    title VARCHAR(36)   NOT NULL,
    point geography(POINT)  NOT NULL,
    PRIMARY KEY(id)
);