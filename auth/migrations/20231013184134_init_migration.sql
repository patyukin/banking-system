-- +goose Up
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles
(
	id   BIGSERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL
);

CREATE TABLE users
(
	id            BIGSERIAL PRIMARY KEY DEFAULT uuid_generate_v4(),
	name          VARCHAR(255) NOT NULL,
	email         VARCHAR(255) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	created_at    TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP,
	updated_at    TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_roles
(
	user_id BIGSERIAL REFERENCES users (id),
	role_id BIGSERIAL REFERENCES roles (id),
	PRIMARY KEY (user_id, role_id)
);

CREATE TABLE refresh_tokens
(
	id         SERIAL PRIMARY KEY,
	token      UUID        DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
	user_id    BIGSERIAL REFERENCES users (id),
	expired_at TIMESTAMPTZ DEFAULT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
drop table users_roles;
drop table users;
drop table roles;
drop table refresh_tokens;
