--
-- PostgreSQL database dump
--

-- Dumped from database version 15.7
-- Dumped by pg_dump version 15.7

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
-- Name: update_person_ranking(uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_person_ranking(pid uuid) RETURNS integer
    LANGUAGE plpgsql
    AS $$
declare
    current_search_count int;
begin
    select search_count
    into current_search_count
    from persons_ranking
    where person_id = pid
    for update skip locked;

    if not found then
        insert into persons_ranking(person_id, search_count)
        values (pid, 1)
        on conflict (person_id) do update
        set search_count = persons_ranking.search_count + 1
        returning search_count into current_search_count;
    else
        update persons_ranking
        set search_count = current_search_count + 1
        where person_id = pid
        returning search_count into current_search_count;
    end if;

    return current_search_count;
end;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: persons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.persons (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    full_name character varying NOT NULL,
    email character varying NOT NULL,
    avatar character varying NOT NULL,
    location character varying NOT NULL,
    metadata jsonb,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL
);


--
-- Name: persons_ranking; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.persons_ranking (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    person_id uuid NOT NULL,
    search_count integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: persons persons_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.persons
    ADD CONSTRAINT persons_email_key UNIQUE (email);


--
-- Name: persons_ranking persons_ranking_person_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.persons_ranking
    ADD CONSTRAINT persons_ranking_person_id_key UNIQUE (person_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- PostgreSQL database dump complete
--

