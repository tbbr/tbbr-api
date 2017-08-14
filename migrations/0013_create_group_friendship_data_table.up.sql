CREATE TABLE group_friendship_data (
    id integer NOT NULL,
    balance integer,
    positive_user_id integer
);

CREATE SEQUENCE group_friendship_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE group_friendship_data_id_seq OWNED BY group_friendship_data.id;

ALTER TABLE ONLY group_friendship_data ALTER COLUMN id SET DEFAULT nextval('group_friendship_data_id_seq'::regclass);

SELECT pg_catalog.setval('group_friendship_data_id_seq', 1, false);

ALTER TABLE ONLY group_friendship_data
    ADD CONSTRAINT group_friendship_data_pkey PRIMARY KEY (id);
