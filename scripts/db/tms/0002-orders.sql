CREATE EXTENSION IF NOT EXISTS Postgis;
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
    start_at_fact       TIMESTAMP   DEFAULT NULL,
    end_at              TIMESTAMP 	NOT NULL,
    end_at_fact         TIMESTAMP   DEFAULT NULL,

    type_id             INT         NOT NULL DEFAULT 1,
    status_id			INT			NOT NULL DEFAULT 3,
    points_id           BIGINT[] 	NOT NULL,
    helpers             INT 		NOT NULL DEFAULT 0,
    comment_message     VARCHAR(256) NOT NULL DEFAULT '', 
    is_regular_customer BOOLEAN     NOT NULL DEFAULT FALSE,

    canceled_at         TIMESTAMP            DEFAULT NULL,
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


DROP TABLE IF EXISTS order_bills CASCADE;

CREATE TABLE order_bills (
	id          		BIGINT  GENERATED ALWAYS AS IDENTITY UNIQUE,
	order_id			BIGINT	NOT NULL UNIQUE,
	car_type_id     	INT     NOT NULL,
	car_price			INT		NOT NULL,
	car_hours			INT		NOT NULL,
	helper_count		INT,		
	helper_price 		INT,
	helper_hours 		INT,

	km_count 			INT,
	km_price 			INT,
	price 				INT 	NOT NULL,
	price_at_fact		INT					DEFAULT NULL,
    is_fragile_cargo 	BOOLEAN NOT NULL 	DEFAULT FALSE,
	PRIMARY KEY(id),
	FOREIGN KEY(order_id) 	REFERENCES orders(id)
);

DROP TABLE IF EXISTS order_variables CASCADE;

CREATE TABLE order_variables (
	id          	INT  		GENERATED ALWAYS AS IDENTITY UNIQUE,
	name           	VARCHAR(32) NOT NULL, 
	value 			JSON		NOT NULL
);

INSERT INTO order_variables ( name, value ) VALUES 
('prices', '{
    "bigCarPrice"	:0,
    "bigCarTime"	:0,
    "helperPrice"	:0,
    "helperTime"	:0,
    "fragilePrice"	:0,
    "kmPrice"		:0
}'),
('work_time', '{
    "startAt" : 0,
    "endAt" : 1380
}');