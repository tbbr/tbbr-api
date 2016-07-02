CREATE TABLE friendship_data (
    id integer NOT NULL,
    balance integer,
    positive_user_id integer
);

CREATE SEQUENCE friendship_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE friendship_data_id_seq OWNED BY friendship_data.id;

ALTER TABLE ONLY friendship_data ALTER COLUMN id SET DEFAULT nextval('friendship_data_id_seq'::regclass);

SELECT pg_catalog.setval('friendship_data_id_seq', 1, false);

ALTER TABLE ONLY friendship_data
    ADD CONSTRAINT friendship_data_pkey PRIMARY KEY (id);
