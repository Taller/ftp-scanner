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


CREATE SEQUENCE scanner.empty_folder_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE scanner.empty_folder_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE scanner.empty_folder (
    id bigint DEFAULT nextval('scanner.empty_folder_id_seq'::regclass) NOT NULL,
    folder character varying(4096) NOT NULL,
    server_id bigint NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE scanner.empty_folder OWNER TO postgres;

CREATE SEQUENCE scanner.files_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE scanner.files_id_seq OWNER TO postgres;

CREATE TABLE scanner.file (
    id bigint DEFAULT nextval('scanner.files_id_seq'::regclass) NOT NULL,
    full_name character varying(4096) NOT NULL,
    fname character varying(1024) NOT NULL,
    path character varying(4096) NOT NULL,
    size bigint NOT NULL,
    folder_id bigint NOT NULL,
    server_id bigint NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE scanner.file OWNER TO postgres;

CREATE TABLE scanner.file_attr (
    id bigint NOT NULL,
    file_id bigint NOT NULL,
    hash character varying(1024),
    length bigint NOT NULL
);

ALTER TABLE scanner.file_attr OWNER TO postgres;

CREATE SEQUENCE scanner.file_attrs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE scanner.file_attrs_id_seq OWNER TO postgres;

ALTER SEQUENCE scanner.file_attrs_id_seq OWNED BY scanner.file_attr.id;

CREATE TABLE scanner.folder (
    id bigint NOT NULL,
    name character varying(4096) NOT NULL,
    parent_id bigint,
    server_id bigint NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE scanner.folder OWNER TO postgres;

CREATE SEQUENCE scanner.folder_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE scanner.folder_id_seq OWNER TO postgres;

ALTER SEQUENCE scanner.folder_id_seq OWNED BY scanner.folder.id;

CREATE TABLE scanner.server (
    id bigint NOT NULL,
    name character varying(1024) NOT NULL,
    created timestamp without time zone DEFAULT '2023-02-27 20:39:56.880867'::timestamp without time zone NOT NULL
);

ALTER TABLE scanner.server OWNER TO postgres;

CREATE SEQUENCE scanner.server_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE scanner.server_id_seq OWNER TO postgres;

ALTER SEQUENCE scanner.server_id_seq OWNED BY scanner.server.id;

ALTER TABLE ONLY scanner.file_attr ALTER COLUMN id SET DEFAULT nextval('scanner.file_attrs_id_seq'::regclass);

ALTER TABLE ONLY scanner.folder ALTER COLUMN id SET DEFAULT nextval('scanner.folder_id_seq'::regclass);

ALTER TABLE ONLY scanner.server ALTER COLUMN id SET DEFAULT nextval('scanner.server_id_seq'::regclass);

ALTER TABLE ONLY scanner.empty_folder
    ADD CONSTRAINT empty_folders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY scanner.file_attr
    ADD CONSTRAINT file_attrs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY scanner.file
    ADD CONSTRAINT files_pkey PRIMARY KEY (id);

ALTER TABLE ONLY scanner.folder
    ADD CONSTRAINT folder_name_unique UNIQUE (name, server_id);

ALTER TABLE ONLY scanner.empty_folder
    ADD CONSTRAINT folder_on_server_unique UNIQUE (folder, server_id);

ALTER TABLE ONLY scanner.folder
    ADD CONSTRAINT folder_pkey PRIMARY KEY (id);

ALTER TABLE ONLY scanner.file
    ADD CONSTRAINT fullname_unique UNIQUE (full_name);

ALTER TABLE ONLY scanner.server
    ADD CONSTRAINT server_name_unique UNIQUE (name);

ALTER TABLE ONLY scanner.server
    ADD CONSTRAINT server_pkey PRIMARY KEY (id);

CREATE INDEX fki_file_in_folder_fk ON scanner.file USING btree (folder_id);

CREATE INDEX fki_folder_on_server_fk ON scanner.folder USING btree (server_id);

CREATE INDEX fki_t ON scanner.file USING btree (server_id);

ALTER TABLE ONLY scanner.file
    ADD CONSTRAINT file_in_folder_fk FOREIGN KEY (folder_id) REFERENCES scanner.folder(id) NOT VALID;

ALTER TABLE ONLY scanner.file
    ADD CONSTRAINT file_to_server_fk FOREIGN KEY (server_id) REFERENCES scanner.server(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;

ALTER TABLE ONLY scanner.folder
    ADD CONSTRAINT folder_on_server_fk FOREIGN KEY (server_id) REFERENCES scanner.server(id) NOT VALID;

ALTER TABLE ONLY scanner.file_attr
    ADD CONSTRAINT to_file_fk FOREIGN KEY (file_id) REFERENCES scanner.file(id);


