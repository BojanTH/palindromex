CREATE TABLE IF NOT EXISTS _user (
    id SERIAL PRIMARY KEY,
    name character varying(255),
    email character varying(255) NOT NULL,
    password character varying(255),
    enabled boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS api_key (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    key character varying(255),
    enabled boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES _user(id)
);

CREATE INDEX user_key_idx ON api_key (user_id, key);

CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    content bytea NOT NULL,
    palindrome boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES _user(id)
);

CREATE INDEX user_idx ON message (user_id);
