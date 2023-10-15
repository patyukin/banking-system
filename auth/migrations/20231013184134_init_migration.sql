-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles
(
	uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL
);

CREATE TABLE users
(
	uuid          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name          VARCHAR(255) NOT NULL,
	email         VARCHAR(255) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	created_at    TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP,
	updated_at    TIMESTAMPTZ      DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_roles
(
	user_uuid UUID REFERENCES users (uuid),
	role_uuid UUID REFERENCES roles (uuid),
	PRIMARY KEY (user_uuid, role_uuid)
);

-- +goose Down
drop table users_roles;
drop table users;
drop table roles;
