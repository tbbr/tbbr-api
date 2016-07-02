CREATE TABLE groups (
    id integer NOT NULL,
    name text,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    hash_id text
);

CREATE SEQUENCE groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE groups_id_seq OWNED BY groups.id;

ALTER TABLE ONLY groups ALTER COLUMN id SET DEFAULT nextval('groups_id_seq'::regclass);

SELECT pg_catalog.setval('groups_id_seq', 1, false);

ALTER TABLE ONLY groups
    ADD CONSTRAINT groups_hash_id_key UNIQUE (hash_id);

ALTER TABLE ONLY groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (id);
