CREATE TABLE tokens (
    id integer NOT NULL,
    category text,
    access_token text,
    refresh_token text,
    refresh_expiration timestamp with time zone,
    auth_expiration timestamp with time zone,
    expired boolean,
    user_id integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE SEQUENCE tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tokens_id_seq OWNED BY tokens.id;

SELECT pg_catalog.setval('tokens_id_seq', 1, false);

ALTER TABLE ONLY tokens ALTER COLUMN id SET DEFAULT nextval('tokens_id_seq'::regclass);

ALTER TABLE ONLY tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);
