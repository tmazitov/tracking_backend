DROP TABLE IF EXISTS default_prices CASCADE;

CREATE TABLE default_prices (
	id          	INT  	GENERATED ALWAYS AS IDENTITY UNIQUE,
	name           	VARCHAR(32) NOT NULL, 
	value 			INT			NOT NULL
);

INSERT INTO default_prices ( name, value) VALUES 
	('big_car_price', 500),
	('helper_price', 400),
	('fragile_price', 300)
;