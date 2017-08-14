CREATE TABLE group_transactions (
    id integer NOT NULL,
    amount integer,
    senders integer ARRAY NOT NULL,
    recipients integer ARRAY NOT NULL,
    senderSplits integer ARRAY NOT NULL,
    recipientSplits integer ARRAY NOT NULL,
    group_id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,

    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE SEQUENCE group_transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE group_transactions_id_seq OWNED BY group_transactions.id;

SELECT pg_catalog.setval('group_transactions_id_seq', 1, false);

ALTER TABLE ONLY group_transactions ALTER COLUMN id SET DEFAULT nextval('group_transactions_id_seq'::regclass);

ALTER TABLE ONLY group_transactions
    ADD CONSTRAINT group_transactions_pkey PRIMARY KEY (id);
