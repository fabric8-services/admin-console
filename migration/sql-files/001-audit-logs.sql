CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- event types
CREATE TABLE event_types (
    id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
    name varchar
);

insert into event_types (id, name) values ('7aea0277-d6fa-4df9-8224-a27fa4096ec7', 'user_search');

-- index to query event type by name, which must be unique
CREATE UNIQUE INDEX uix_event_type_name ON event_types USING btree (name);

-- audit logs
CREATE TABLE audit_logs (
    id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone NOT NULL default now(),
    identity_id uuid NOT NULL,
    event_type_id uuid NOT NULL,
    event_params jsonb NOT NULL
);

-- Add a foreign key constraint to identities
ALTER TABLE audit_logs add constraint audit_logs_event_types_fk foreign key (event_type_id) REFERENCES event_types (id);

