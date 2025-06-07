SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: bot_user; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bot_user (
    id bigint NOT NULL,
    tg_id text NOT NULL,
    person_id bigint
);


--
-- Name: bot_user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bot_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bot_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bot_user_id_seq OWNED BY public.bot_user.id;


--
-- Name: game_player; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.game_player (
    player_id bigint NOT NULL,
    game_id bigint NOT NULL,
    team bigint NOT NULL
);


--
-- Name: guest_person; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.guest_person (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    uid_string text NOT NULL
);


--
-- Name: guest_person_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.guest_person_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: guest_person_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.guest_person_id_seq OWNED BY public.guest_person.id;


--
-- Name: hokm4_game; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hokm4_game (
    id bigint NOT NULL,
    team_one_tricks_score bigint DEFAULT 0 NOT NULL,
    team_two_tricks_score bigint DEFAULT 0 NOT NULL,
    bet_amount bigint DEFAULT 0 NOT NULL,
    created_stamp timestamp with time zone DEFAULT now(),
    end_stamp timestamp with time zone
);


--
-- Name: hokm4_game_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.hokm4_game_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: hokm4_game_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.hokm4_game_id_seq OWNED BY public.hokm4_game.id;


--
-- Name: person; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.person (
    id bigint NOT NULL,
    display_name text,
    bio text,
    coin bigint NOT NULL,
    created timestamp with time zone DEFAULT now()
);


--
-- Name: person_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.person_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: person_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.person_id_seq OWNED BY public.person.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: trick; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.trick (
    trick_id bigint NOT NULL,
    game_id bigint DEFAULT 0 NOT NULL,
    hokm bigint,
    hakem_index bigint DEFAULT 0 NOT NULL,
    team_one_turn_score bigint DEFAULT 0 NOT NULL,
    team_two_turn_score bigint DEFAULT 0 NOT NULL
);


--
-- Name: trick_trick_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.trick_trick_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: trick_trick_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.trick_trick_id_seq OWNED BY public.trick.trick_id;


--
-- Name: turn; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.turn (
    turn_id bigint NOT NULL,
    trick_id bigint DEFAULT 0 NOT NULL,
    moves text DEFAULT 0 NOT NULL
);


--
-- Name: turn_turn_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.turn_turn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: turn_turn_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.turn_turn_id_seq OWNED BY public.turn.turn_id;


--
-- Name: user_statistic; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_statistic (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    win bigint DEFAULT 0 NOT NULL,
    lose bigint DEFAULT 0 NOT NULL,
    total_tricks_won bigint DEFAULT 0 NOT NULL,
    total_tricks_lost bigint DEFAULT 0 NOT NULL,
    total_turns_won bigint DEFAULT 0 NOT NULL,
    total_turns_lost bigint DEFAULT 0 NOT NULL,
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: user_statistic_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_statistic_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_statistic_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_statistic_id_seq OWNED BY public.user_statistic.id;


--
-- Name: bot_user id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_user ALTER COLUMN id SET DEFAULT nextval('public.bot_user_id_seq'::regclass);


--
-- Name: guest_person id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guest_person ALTER COLUMN id SET DEFAULT nextval('public.guest_person_id_seq'::regclass);


--
-- Name: hokm4_game id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hokm4_game ALTER COLUMN id SET DEFAULT nextval('public.hokm4_game_id_seq'::regclass);


--
-- Name: person id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person ALTER COLUMN id SET DEFAULT nextval('public.person_id_seq'::regclass);


--
-- Name: trick trick_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.trick ALTER COLUMN trick_id SET DEFAULT nextval('public.trick_trick_id_seq'::regclass);


--
-- Name: turn turn_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.turn ALTER COLUMN turn_id SET DEFAULT nextval('public.turn_turn_id_seq'::regclass);


--
-- Name: user_statistic id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_statistic ALTER COLUMN id SET DEFAULT nextval('public.user_statistic_id_seq'::regclass);


--
-- Name: bot_user bot_user_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_user
    ADD CONSTRAINT bot_user_pkey PRIMARY KEY (id);


--
-- Name: bot_user bot_user_tg_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_user
    ADD CONSTRAINT bot_user_tg_id_key UNIQUE (tg_id);


--
-- Name: game_player game_player_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.game_player
    ADD CONSTRAINT game_player_pkey PRIMARY KEY (player_id, game_id);


--
-- Name: guest_person guest_person_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guest_person
    ADD CONSTRAINT guest_person_pkey PRIMARY KEY (id);


--
-- Name: hokm4_game hokm4_game_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hokm4_game
    ADD CONSTRAINT hokm4_game_pkey PRIMARY KEY (id);


--
-- Name: person person_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.person
    ADD CONSTRAINT person_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: trick trick_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.trick
    ADD CONSTRAINT trick_pkey PRIMARY KEY (trick_id);


--
-- Name: turn turn_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.turn
    ADD CONSTRAINT turn_pkey PRIMARY KEY (turn_id);


--
-- Name: user_statistic user_statistic_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_statistic
    ADD CONSTRAINT user_statistic_pkey PRIMARY KEY (id);


--
-- Name: trick fk_game; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.trick
    ADD CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES public.hokm4_game(id);


--
-- Name: game_player fk_game; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.game_player
    ADD CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES public.hokm4_game(id);


--
-- Name: bot_user fk_person; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_user
    ADD CONSTRAINT fk_person FOREIGN KEY (person_id) REFERENCES public.person(id);


--
-- Name: game_player fk_player; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.game_player
    ADD CONSTRAINT fk_player FOREIGN KEY (player_id) REFERENCES public.person(id);


--
-- Name: turn fk_trick; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.turn
    ADD CONSTRAINT fk_trick FOREIGN KEY (trick_id) REFERENCES public.trick(trick_id);


--
-- Name: guest_person guest_person_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.guest_person
    ADD CONSTRAINT guest_person_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.person(id);


--
-- Name: user_statistic user_statistic_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_statistic
    ADD CONSTRAINT user_statistic_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.person(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250417173613'),
    ('20250417173930'),
    ('20250417174013'),
    ('20250417174043'),
    ('20250417174107'),
    ('20250417174129'),
    ('20250504104006'),
    ('20250601070008'),
    ('20250601075712'),
    ('20250601081430'),
    ('20250601082252'),
    ('20250607081626');
