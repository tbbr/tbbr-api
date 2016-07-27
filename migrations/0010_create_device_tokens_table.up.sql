CREATE TABLE device_tokens (
    id integer NOT NULL,
    token text,
    user_id integer,
    type text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE SEQUENCE device_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE device_tokens_id_seq OWNED BY device_tokens.id;

SELECT pg_catalog.setval('device_tokens_id_seq', 1, false);

ALTER TABLE ONLY device_tokens ALTER COLUMN id SET DEFAULT nextval('device_tokens_id_seq'::regclass);

ALTER TABLE ONLY device_tokens
    ADD CONSTRAINT device_tokens_pkey PRIMARY KEY (id);

ALTER TABLE device_tokens ADD CONSTRAINT token_must_be_unique UNIQUE (token);
