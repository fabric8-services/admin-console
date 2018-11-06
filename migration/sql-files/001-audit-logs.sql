CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- event types
CREATE TABLE event_type (
    event_type_id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
    name varchar
);

insert into event_type (event_type_id, name) values ('7aea0277-d6fa-4df9-8224-a27fa4096ec7', 'user_search');

-- index to query event type by name, which must be unique
CREATE UNIQUE INDEX uix_event_type_name ON event_type USING btree (name);

-- audit logs
CREATE TABLE audit_log (
    audit_log_id uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone NOT NULL default now(),
    identity_id uuid NOT NULL,
    event_type_id uuid NOT NULL,
    event_params jsonb NOT NULL
);

-- Add a foreign key constraint to event_type
ALTER TABLE audit_log add constraint audit_logs_event_type_fk foreign key (event_type_id) REFERENCES event_type (event_type_id);

