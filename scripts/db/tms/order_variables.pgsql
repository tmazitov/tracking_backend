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