CREATE TABLE transactions (
    id integer NOT NULL,
    amount integer,
    memo text,
    sender_id integer,
    creator_id integer,
    related_object_id integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type text,
    related_object_type text,
    recipient_id integer
);

CREATE SEQUENCE transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE transactions_id_seq OWNED BY transactions.id;

SELECT pg_catalog.setval('transactions_id_seq', 1, false);

ALTER TABLE ONLY transactions ALTER COLUMN id SET DEFAULT nextval('transactions_id_seq'::regclass);

ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
