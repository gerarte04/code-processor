CREATE TABLE tasks (
	id uuid PRIMARY KEY NOT NULL,
	finished boolean NOT NULL DEFAULT false,

	output text DEFAULT '',
	status_code int8 DEFAULT -1,
	
	translator varchar(50) NOT NULL,
	code text
);

CREATE TABLE users (
	id uuid PRIMARY KEY NOT NULL,
	login varchar(50) NOT NULL UNIQUE,
	password varchar(50) NOT NULL
);

CREATE INDEX ON users (login);
