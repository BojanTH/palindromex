CREATE TABLE IF NOT EXISTS _user (
    id SERIAL PRIMARY KEY,
    name character varying(255),
    email character varying(255) NOT NULL,
    password character varying(255),
    enabled boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);