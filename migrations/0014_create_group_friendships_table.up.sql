CREATE TABLE group_friendships (
    id integer NOT NULL,
    user_id integer NOT NULL,
    friend_id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    group_friendship_data_id integer NOT NULL,
    friendship_id integer,
    group_id integer NOT NULL,

    FOREIGN KEY (group_friendship_data_id) REFERENCES group_friendship_data(id),
    FOREIGN KEY (friendship_id) REFERENCES friendships(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE SEQUENCE group_friendship_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE group_friendship_id_seq OWNED BY group_friendships.id;

SELECT pg_catalog.setval('group_friendship_id_seq', 1, false);

ALTER TABLE ONLY group_friendships ALTER COLUMN id SET DEFAULT nextval('group_friendship_id_seq'::regclass);

ALTER TABLE ONLY group_friendships
    ADD CONSTRAINT group_friendship_pkey PRIMARY KEY (id);
