--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: friendship_data; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

CREATE TABLE friendship_data (
    id integer NOT NULL,
    balance integer,
    positive_user_id integer
);


ALTER TABLE friendship_data OWNER TO maazali;

--
-- Name: friendship_data_id_seq; Type: SEQUENCE; Schema: public; Owner: maazali
--

CREATE SEQUENCE friendship_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE friendship_data_id_seq OWNER TO maazali;

--
-- Name: friendship_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maazali
--

ALTER SEQUENCE friendship_data_id_seq OWNED BY friendship_data.id;


--
-- Name: friendships; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

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


ALTER TABLE friendships OWNER TO maazali;

--
-- Name: friendship_id_seq; Type: SEQUENCE; Schema: public; Owner: maazali
--

CREATE SEQUENCE friendship_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE friendship_id_seq OWNER TO maazali;

--
-- Name: friendship_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maazali
--

ALTER SEQUENCE friendship_id_seq OWNED BY friendships.id;


--
-- Name: group_users; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

CREATE TABLE group_users (
    group_id integer,
    user_id integer
);


ALTER TABLE group_users OWNER TO maazali;

--
-- Name: groups; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

CREATE TABLE groups (
    id integer NOT NULL,
    name text,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    hash_id text
);


ALTER TABLE groups OWNER TO maazali;

--
-- Name: groups_id_seq; Type: SEQUENCE; Schema: public; Owner: maazali
--

CREATE SEQUENCE groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE groups_id_seq OWNER TO maazali;

--
-- Name: groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maazali
--

ALTER SEQUENCE groups_id_seq OWNED BY groups.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

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


ALTER TABLE tokens OWNER TO maazali;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: maazali
--

CREATE SEQUENCE tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE tokens_id_seq OWNER TO maazali;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maazali
--

ALTER SEQUENCE tokens_id_seq OWNED BY tokens.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

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


ALTER TABLE transactions OWNER TO maazali;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: maazali
--

CREATE SEQUENCE transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE transactions_id_seq OWNER TO maazali;

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maazali
--

ALTER SEQUENCE transactions_id_seq OWNED BY transactions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: maazali; Tablespace:
--

CREATE TABLE users (
    id integer NOT NULL,
    name text,
    email text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    external_id text,
    gender text,
    avatar_url text
);


ALTER TABLE users OWNER TO maazali;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: maazali
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_id_seq OWNER TO maazali;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maazali
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: maazali
--

ALTER TABLE ONLY friendship_data ALTER COLUMN id SET DEFAULT nextval('friendship_data_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: maazali
--

ALTER TABLE ONLY friendships ALTER COLUMN id SET DEFAULT nextval('friendship_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: maazali
--

ALTER TABLE ONLY groups ALTER COLUMN id SET DEFAULT nextval('groups_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: maazali
--

ALTER TABLE ONLY tokens ALTER COLUMN id SET DEFAULT nextval('tokens_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: maazali
--

ALTER TABLE ONLY transactions ALTER COLUMN id SET DEFAULT nextval('transactions_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: maazali
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Data for Name: friendship_data; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY friendship_data (id, balance, positive_user_id) FROM stdin;
\.


--
-- Name: friendship_data_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maazali
--

SELECT pg_catalog.setval('friendship_data_id_seq', 1, false);


--
-- Name: friendship_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maazali
--

SELECT pg_catalog.setval('friendship_id_seq', 1, false);


--
-- Data for Name: friendships; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY friendships (id, user_id, friend_id, created_at, updated_at, deleted_at, hash_id, friendship_data_id) FROM stdin;
\.


--
-- Data for Name: group_users; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY group_users (group_id, user_id) FROM stdin;
\.


--
-- Data for Name: groups; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY groups (id, name, description, created_at, updated_at, deleted_at, hash_id) FROM stdin;
\.


--
-- Name: groups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maazali
--

SELECT pg_catalog.setval('groups_id_seq', 1, false);


--
-- Data for Name: tokens; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY tokens (id, category, access_token, refresh_token, refresh_expiration, auth_expiration, expired, user_id, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Name: tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maazali
--

SELECT pg_catalog.setval('tokens_id_seq', 1, false);


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY transactions (id, amount, memo, sender_id, creator_id, related_object_id, created_at, updated_at, deleted_at, type, related_object_type, recipient_id) FROM stdin;
\.


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maazali
--

SELECT pg_catalog.setval('transactions_id_seq', 1, false);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: maazali
--

COPY users (id, name, email, created_at, updated_at, deleted_at, external_id, gender, avatar_url) FROM stdin;
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: maazali
--

SELECT pg_catalog.setval('users_id_seq', 1, false);


--
-- Name: friendship_data_pkey; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY friendship_data
    ADD CONSTRAINT friendship_data_pkey PRIMARY KEY (id);


--
-- Name: friendship_pkey; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY friendships
    ADD CONSTRAINT friendship_pkey PRIMARY KEY (id);


--
-- Name: groups_hash_id_key; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY groups
    ADD CONSTRAINT groups_hash_id_key UNIQUE (hash_id);


--
-- Name: groups_pkey; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (id);


--
-- Name: tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users_external_id_key; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_external_id_key UNIQUE (external_id);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: maazali; Tablespace:
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: public; Type: ACL; Schema: -; Owner: maazali
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM maazali;
GRANT ALL ON SCHEMA public TO maazali;
GRANT ALL ON SCHEMA public TO PUBLIC;
