--
-- PostgreSQL database dump
--

\restrict E8DxJXpOheD9fqsKuaLhbdQPywaJJWARxVEcS5CBdXPJ5eXH95MgwNrIztwGceR

-- Dumped from database version 16.15 (Debian 16.15-1.pgdg13+2)
-- Dumped by pg_dump version 16.15 (Debian 16.15-1.pgdg13+2)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.accounts (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    name character varying(50) NOT NULL,
    type character varying(10) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    currency character varying(10) NOT NULL,
    target_amount bigint,
    start_date timestamp without time zone,
    target_date timestamp without time zone,
    category_id uuid,
    note text,
    completed_at timestamp without time zone,
    cancelled_at timestamp without time zone,
    finalized_amount bigint,
    closed_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT accounts_check CHECK (((((type)::text <> 'goal'::text) AND (target_amount IS NULL) AND (start_date IS NULL) AND (target_date IS NULL) AND (category_id IS NULL) AND (note IS NULL) AND (completed_at IS NULL) AND (cancelled_at IS NULL) AND (finalized_amount IS NULL)) OR (((type)::text = 'goal'::text) AND (target_amount IS NOT NULL) AND (start_date IS NOT NULL) AND (category_id IS NOT NULL)))),
    CONSTRAINT accounts_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'closed'::character varying])::text[]))),
    CONSTRAINT accounts_type_check CHECK (((type)::text = ANY ((ARRAY['checking'::character varying, 'savings'::character varying, 'goal'::character varying])::text[])))
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- Name: expense_categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expense_categories (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    name character varying(30) NOT NULL,
    type character varying(10) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    archived_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT expense_categories_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'archived'::character varying])::text[]))),
    CONSTRAINT expense_categories_type_check CHECK (((type)::text = ANY ((ARRAY['needs'::character varying, 'wants'::character varying, 'savings'::character varying])::text[])))
);


ALTER TABLE public.expense_categories OWNER TO postgres;

--
-- Name: income_categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_categories (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    name character varying(30) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    archived_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT income_categories_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'archived'::character varying])::text[])))
);


ALTER TABLE public.income_categories OWNER TO postgres;

--
-- Name: ledger; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ledger (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    type character varying(20) NOT NULL,
    from_account_id uuid,
    to_account_id uuid,
    income_category_id uuid,
    expense_category_id uuid,
    amount bigint NOT NULL,
    description text NOT NULL,
    date timestamp without time zone NOT NULL,
    note text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT ledger_amount_check CHECK ((amount > 0)),
    CONSTRAINT ledger_check CHECK (((from_account_id IS NOT NULL) OR (to_account_id IS NOT NULL))),
    CONSTRAINT ledger_check1 CHECK (((from_account_id IS NULL) OR (to_account_id IS NULL) OR (from_account_id <> to_account_id))),
    CONSTRAINT ledger_check2 CHECK (((((type)::text = 'income'::text) AND (to_account_id IS NOT NULL) AND (from_account_id IS NULL) AND (income_category_id IS NOT NULL) AND (expense_category_id IS NULL)) OR (((type)::text = 'expense'::text) AND (from_account_id IS NOT NULL) AND (to_account_id IS NULL) AND (expense_category_id IS NOT NULL) AND (income_category_id IS NULL)) OR (((type)::text = 'transfer'::text) AND (from_account_id IS NOT NULL) AND (to_account_id IS NOT NULL) AND (income_category_id IS NULL)) OR (((type)::text = 'interest'::text) AND (from_account_id IS NULL) AND (to_account_id IS NOT NULL) AND (income_category_id IS NOT NULL) AND (expense_category_id IS NULL)) OR (((type)::text = 'goal-fulfillment'::text) AND (from_account_id IS NOT NULL) AND (to_account_id IS NULL) AND (income_category_id IS NULL) AND (expense_category_id IS NOT NULL)) OR (((type)::text = 'initial-balance'::text) AND (from_account_id IS NULL) AND (to_account_id IS NOT NULL) AND (income_category_id IS NULL) AND (expense_category_id IS NULL)))),
    CONSTRAINT ledger_type_check CHECK (((type)::text = ANY ((ARRAY['income'::character varying, 'expense'::character varying, 'transfer'::character varying, 'interest'::character varying, 'goal-fulfillment'::character varying, 'initial-balance'::character varying])::text[])))
);


ALTER TABLE public.ledger OWNER TO postgres;

--
-- Name: monthly_budgets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.monthly_budgets (
    category_id uuid NOT NULL,
    month date NOT NULL,
    amount bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT monthly_budgets_amount_check CHECK ((amount >= 0))
);


ALTER TABLE public.monthly_budgets OWNER TO postgres;

--
-- Name: monthly_summary; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.monthly_summary (
    account_id uuid NOT NULL,
    month date NOT NULL,
    money_in bigint NOT NULL,
    money_out bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT monthly_summary_money_in_check CHECK ((money_in >= 0)),
    CONSTRAINT monthly_summary_money_out_check CHECK ((money_out >= 0)),
    CONSTRAINT monthly_summary_month_check CHECK ((month = date_trunc('month'::text, (month)::timestamp with time zone)))
);


ALTER TABLE public.monthly_summary OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying(50) NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: accounts accounts_user_id_type_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_user_id_type_name_key UNIQUE (user_id, type, name);


--
-- Name: expense_categories expense_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_pkey PRIMARY KEY (id);


--
-- Name: expense_categories expense_categories_user_id_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_user_id_name_key UNIQUE (user_id, name);


--
-- Name: income_categories income_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_pkey PRIMARY KEY (id);


--
-- Name: income_categories income_categories_user_id_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_user_id_name_key UNIQUE (user_id, name);


--
-- Name: ledger ledger_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger
    ADD CONSTRAINT ledger_pkey PRIMARY KEY (id);


--
-- Name: monthly_budgets monthly_budgets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.monthly_budgets
    ADD CONSTRAINT monthly_budgets_pkey PRIMARY KEY (category_id, month);


--
-- Name: monthly_summary monthly_summary_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.monthly_summary
    ADD CONSTRAINT monthly_summary_pkey PRIMARY KEY (account_id, month);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_name_key UNIQUE (name);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_ledger_from_account; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ledger_from_account ON public.ledger USING btree (from_account_id);


--
-- Name: idx_ledger_to_account; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ledger_to_account ON public.ledger USING btree (to_account_id);


--
-- Name: idx_ledger_user_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ledger_user_date ON public.ledger USING btree (user_id, date DESC, created_at DESC);


--
-- Name: idx_ledger_user_type_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ledger_user_type_date ON public.ledger USING btree (user_id, type, date DESC);


--
-- Name: accounts accounts_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.expense_categories(id);


--
-- Name: accounts accounts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: expense_categories expense_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: income_categories income_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: ledger ledger_expense_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger
    ADD CONSTRAINT ledger_expense_category_id_fkey FOREIGN KEY (expense_category_id) REFERENCES public.expense_categories(id);


--
-- Name: ledger ledger_from_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger
    ADD CONSTRAINT ledger_from_account_id_fkey FOREIGN KEY (from_account_id) REFERENCES public.accounts(id);


--
-- Name: ledger ledger_income_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger
    ADD CONSTRAINT ledger_income_category_id_fkey FOREIGN KEY (income_category_id) REFERENCES public.income_categories(id);


--
-- Name: ledger ledger_to_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger
    ADD CONSTRAINT ledger_to_account_id_fkey FOREIGN KEY (to_account_id) REFERENCES public.accounts(id);


--
-- Name: ledger ledger_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ledger
    ADD CONSTRAINT ledger_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: monthly_budgets monthly_budgets_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.monthly_budgets
    ADD CONSTRAINT monthly_budgets_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.expense_categories(id) ON DELETE CASCADE;


--
-- Name: monthly_summary monthly_summary_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.monthly_summary
    ADD CONSTRAINT monthly_summary_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict E8DxJXpOheD9fqsKuaLhbdQPywaJJWARxVEcS5CBdXPJ5eXH95MgwNrIztwGceR

