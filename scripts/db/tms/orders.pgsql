DROP TABLE IF EXISTS points CASCADE;

CREATE TABLE points (
    id          BIGINT  GENERATED ALWAYS AS IDENTITY UNIQUE,
    step_id     INT2                NOT NULL,
    floor       INT8                NOT NULL DEFAULT 1,
    title       VARCHAR(256)         NOT NULL,
    point       geography(POINT)    NOT NULL,
    PRIMARY KEY(id)
);


DROP TABLE IF EXISTS orders CASCADE; 

CREATE TABLE orders (
    id                  BIGINT		GENERATED ALWAYS AS IDENTITY UNIQUE,

	owner_id            BIGINT		NOT NULL,
	worker_id			BIGINT,
	manager_id			BIGINT,

    title               VARCHAR(256) NOT NULL,

    start_at            TIMESTAMP	NOT NULL,
    end_at              TIMESTAMP 	DEFAULT NULL,

    type_id             INT         NOT NULL DEFAULT 1,
    status_id			INT			NOT NULL DEFAULT 1,
    points_id           BIGINT[] 	NOT NULL,
    helpers             INT 		NOT NULL DEFAULT 0,
    comment_message     VARCHAR(256) NOT NULL DEFAULT '', 
    is_fragile_cargo    BOOLEAN 	NOT NULL DEFAULT FALSE,
    is_regular_customer BOOLEAN     NOT NULL DEFAULT FALSE,

    canceled_a          TIMESTAMP            DEFAULT NULL,
    created_at			TIMESTAMP   NOT NULL DEFAULT NOW(),
    edited_at 			TIMESTAMP	NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id),
    CONSTRAINT fk_worker FOREIGN KEY (worker_id) REFERENCES users(id),
    CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES users(id),
    PRIMARY KEY(id)
);

DROP TABLE IF EXISTS points_to_orders CASCADE;

CREATE TABLE points_to_orders (
    order_id  BIGINT NOT NULL,
    point_id  BIGINT NOT NULL,
    PRIMARY KEY (order_id, point_id),
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id),
    CONSTRAINT fk_point FOREIGN KEY(point_id) REFERENCES points(id)
);
