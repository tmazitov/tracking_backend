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

	km 					INT,
	price 				INT 	NOT NULL,
	price_at_fact		INT					DEFAULT NULL,
    is_fragile_cargo 	BOOLEAN NOT NULL 	DEFAULT FALSE,
	PRIMARY KEY(id),
	FOREIGN KEY(order_id) 	REFERENCES orders(id)
);