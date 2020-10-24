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
    updated_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    content text NOT NULL,
    palindrome boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
