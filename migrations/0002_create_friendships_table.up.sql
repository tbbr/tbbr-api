CREATE TABLE friendships (
    id integer NOT NULL,
    user_id integer,
    friend_id integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    hash_id text,
    friendship_data_id integer
);

CREATE SEQUENCE friendship_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE friendship_id_seq OWNED BY friendships.id;

SELECT pg_catalog.setval('friendship_id_seq', 1, false);

ALTER TABLE ONLY friendships ALTER COLUMN id SET DEFAULT nextval('friendship_id_seq'::regclass);

ALTER TABLE ONLY friendships
    ADD CONSTRAINT friendship_pkey PRIMARY KEY (id);
