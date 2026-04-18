--
-- PostgreSQL database dump
--

\restrict B3Kz3SCyTXQdEseSKXvAws9ofAJDd6yE05VvZtSdOkusgSUNS14tpli3dNOQIwE

-- Dumped from database version 16.13 (Debian 16.13-1.pgdg13+1)
-- Dumped by pg_dump version 16.13 (Debian 16.13-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: account_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.account_status AS ENUM (
    'active',
    'inactive',
    'blocked'
);


ALTER TYPE public.account_status OWNER TO postgres;

--
-- Name: transaction_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.transaction_type AS ENUM (
    'deposit',
    'withdraw',
    'transfer_in',
    'transfer_out'
);


ALTER TYPE public.transaction_type OWNER TO postgres;

--
-- Name: prevent_account_transactions_mutation(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.prevent_account_transactions_mutation() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    RAISE EXCEPTION 'account_transactions is immutable';
END;
$$;


ALTER FUNCTION public.prevent_account_transactions_mutation() OWNER TO postgres;

--
-- Name: account_number_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.account_number_seq
    START WITH 10000000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.account_number_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: account_transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_transactions (
    id uuid NOT NULL,
    account_id uuid NOT NULL,
    type public.transaction_type NOT NULL,
    amount bigint NOT NULL,
    balance_after bigint NOT NULL,
    reference_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    related_account_id uuid,
    idempotency_key character varying(100)
);


ALTER TABLE public.account_transactions OWNER TO postgres;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.accounts (
    id uuid NOT NULL,
    customer_id uuid NOT NULL,
    number character varying(20) NOT NULL,
    branch character varying(10) NOT NULL,
    balance bigint DEFAULT 0 NOT NULL,
    status public.account_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- Name: customers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customers (
    id uuid NOT NULL,
    name character varying(120) NOT NULL,
    cpf character varying(11) NOT NULL,
    email character varying(120),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_cpf_format CHECK (((cpf)::text ~ '^\d{11}$'::text))
);


ALTER TABLE public.customers OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: user_sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    token_hash character(64) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.user_sessions OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(120) NOT NULL,
    password_hash text NOT NULL,
    role character varying(20) NOT NULL,
    customer_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    CONSTRAINT chk_users_customer_role_consistency CHECK (((((role)::text = 'customer'::text) AND (customer_id IS NOT NULL)) OR ((role)::text <> 'customer'::text)))
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: account_transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.account_transactions (id, account_id, type, amount, balance_after, reference_id, created_at, related_account_id, idempotency_key) FROM stdin;
\.


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.accounts (id, customer_id, number, branch, balance, status, created_at) FROM stdin;
\.


--
-- Data for Name: customers; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.customers (id, name, cpf, email, created_at) FROM stdin;
3f824547-96e5-4df6-a670-1017593c7f15	Rudson Alves	86264184772	\N	2026-04-17 17:10:45.755699+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
9	f
\.


--
-- Data for Name: user_sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_sessions (id, user_id, token_hash, expires_at, revoked_at, created_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password_hash, role, customer_id, created_at, updated_at, status) FROM stdin;
c541fbf8-798f-42f8-9ffd-b74633585973	rudsonalves67@gmail.com	$2a$10$T/00miLcxMGQfLzS9zhoSOgTJMK6fCm2X0rO8BMV5MhrfAXtPRkCG	admin	3f824547-96e5-4df6-a670-1017593c7f15	2026-04-17 17:10:45.755699	2026-04-17 17:10:45.755699	pending
\.


--
-- Name: account_number_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.account_number_seq', 10000000, false);


--
-- Name: account_transactions account_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_transactions
    ADD CONSTRAINT account_transactions_pkey PRIMARY KEY (id);


--
-- Name: accounts accounts_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_number_key UNIQUE (number);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: customers customers_cpf_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_cpf_key UNIQUE (cpf);


--
-- Name: customers customers_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_email_key UNIQUE (email);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_sessions user_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_token_hash_key UNIQUE (token_hash);


--
-- Name: users users_customer_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_customer_id_key UNIQUE (customer_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_account_transactions_account_created; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_account_transactions_account_created ON public.account_transactions USING btree (account_id, created_at DESC);


--
-- Name: idx_account_transactions_account_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_account_transactions_account_id ON public.account_transactions USING btree (account_id);


--
-- Name: idx_account_transactions_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_account_transactions_created_at ON public.account_transactions USING btree (created_at DESC);


--
-- Name: idx_account_transactions_reference_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_account_transactions_reference_id ON public.account_transactions USING btree (reference_id);


--
-- Name: idx_account_transactions_reference_type; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_account_transactions_reference_type ON public.account_transactions USING btree (reference_id, type);


--
-- Name: idx_accounts_customer_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_accounts_customer_id ON public.accounts USING btree (customer_id);


--
-- Name: idx_user_sessions_expires_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_sessions_expires_at ON public.user_sessions USING btree (expires_at);


--
-- Name: idx_user_sessions_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_user_sessions_user_id ON public.user_sessions USING btree (user_id);


--
-- Name: ux_account_transactions_idempotency; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_account_transactions_idempotency ON public.account_transactions USING btree (account_id, idempotency_key) WHERE (idempotency_key IS NOT NULL);


--
-- Name: ux_account_transactions_transfer_pair; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ux_account_transactions_transfer_pair ON public.account_transactions USING btree (reference_id, type) WHERE ((reference_id IS NOT NULL) AND (type = ANY (ARRAY['transfer_in'::public.transaction_type, 'transfer_out'::public.transaction_type])));


--
-- Name: account_transactions trg_account_transactions_no_mutation; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trg_account_transactions_no_mutation BEFORE DELETE OR UPDATE ON public.account_transactions FOR EACH ROW EXECUTE FUNCTION public.prevent_account_transactions_mutation();


--
-- Name: account_transactions account_transactions_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_transactions
    ADD CONSTRAINT account_transactions_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id);


--
-- Name: accounts accounts_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE RESTRICT;


--
-- Name: users fk_users_customer_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_customer_id FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE SET NULL;


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict B3Kz3SCyTXQdEseSKXvAws9ofAJDd6yE05VvZtSdOkusgSUNS14tpli3dNOQIwE

