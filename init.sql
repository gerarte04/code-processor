CREATE TABLE tasks (
	id uuid PRIMARY KEY NOT NULL UNIQUE,
	finished boolean NOT NULL,

	output text,
	status_code int8,
	
	translator varchar(50) NOT NULL,
	code text
);

CREATE TABLE users (
	id uuid PRIMARY KEY NOT NULL UNIQUE,
	login varchar(50) NOT NULL UNIQUE,
	password varchar(50) NOT NULL
);

CREATE INDEX ON users (login, password);
