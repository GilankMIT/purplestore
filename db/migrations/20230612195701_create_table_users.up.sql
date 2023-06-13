create table users
(
	id BIGSERIAL PRIMARY KEY NOT NULL,
	email VARCHAR(255),
	password VARCHAR(255),
	created_at TIMESTAMPTZ
);