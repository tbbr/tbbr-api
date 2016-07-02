CREATE TABLE users (
    id integer NOT NULL,
    name text,
    email text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    external_id text,
    gender text,
    avatar_url text
);

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE users_id_seq OWNED BY users.id;

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

SELECT pg_catalog.setval('users_id_seq', 1, false);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_external_id_key UNIQUE (external_id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
