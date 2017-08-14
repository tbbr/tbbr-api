CREATE TABLE group_singular_transactions (
    id integer NOT NULL,
    amount integer,
    sender_id integer NOT NULL,
    recipient_id integer NOT NULL,
    group_friendship_id integer NOT NULL,
    group_transaction_id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,

    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    FOREIGN KEY (group_friendship_id) REFERENCES group_friendships(id),
    FOREIGN KEY (group_transaction_id) REFERENCES group_transactions(id)
);

CREATE SEQUENCE group_singular_transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE group_singular_transactions_id_seq OWNED BY group_singular_transactions.id;

SELECT pg_catalog.setval('group_singular_transactions_id_seq', 1, false);

ALTER TABLE ONLY group_singular_transactions ALTER COLUMN id SET DEFAULT nextval('group_singular_transactions_id_seq'::regclass);

ALTER TABLE ONLY group_singular_transactions
    ADD CONSTRAINT group_singular_transactions_pkey PRIMARY KEY (id);
