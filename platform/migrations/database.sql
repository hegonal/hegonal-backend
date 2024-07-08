-- Create TimescaleDB extension if not exists
CREATE EXTENSION IF NOT EXISTS timescaledb;
-- Table: server_locations
CREATE TABLE server_locations (
    server_id varchar(32) NOT NULL,
    server_display_name varchar(32) NOT NULL,
    country char(3),
    CONSTRAINT SERVER_LOCATIONS_PK_1 PRIMARY KEY (server_id)
);
-- Table: users
CREATE TABLE users (
    "user_id" bigint NOT NULL,
    name varchar(32) NOT NULL,
    password char(60) NOT NULL,
    email varchar(255) NOT NULL,
    avatar varchar(255),
    role smallint NOT NULL,
    two_factor_auth varchar(64),
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT USERS_PK_1 PRIMARY KEY ("user_id")
);
-- Table: sessions
CREATE TABLE sessions (
    user_id bigint NOT NULL,
    session char(128) NOT NULL,
    expiry_time timestamp NOT NULL,
    ip varchar(64) NOT NULL,
    device varchar(64) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT SESSIONS_PK_1 PRIMARY KEY ("user_id", session)
);
-- Table: teams
CREATE TABLE teams (
    "team_id" bigint NOT NULL,
    name varchar(64) NOT NULL,
    description varchar(128) NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT TEAMS_PK_1 PRIMARY KEY ("team_id")
);
-- Table: team_members
CREATE TABLE team_members (
    member_id bigint NOT NULL,
    team_id bigint NOT NULL,
    role smallint NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT TEAM_MEMBERS_PK_1 PRIMARY KEY (member_id, team_id)
);
-- Table: team_invites
CREATE TABLE team_invites (
    invite_id bigint NOT NULL,
    team_id bigint NOT NULL,
    user_id bigint NOT NULL,
    "role" smallint NOT NULL,
    "expiry_date" timestamp,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT team_invites_pk PRIMARY KEY ("invite_id")
);
-- Table: http_monitors
CREATE TABLE http_monitors (
    "http_monitor_id" bigint NOT NULL,
    team_id bigint NOT NULL,
    status smallint NOT NULL,
    url varchar(255) NOT NULL,
    interval integer NOT NULL,
    retries smallint NOT NULL,
    retry_interval smallint NOT NULL,
    request_timeout smallint NOT NULL,
    resend_notification smallint NOT NULL,
    follow_redirections boolean NOT NULL,
    max_redirects smallint NOT NULL,
    check_ssl_error boolean NOT NULL,
    ssl_expiry_reminders smallint NOT NULL,
    domain_expiry_reminders smallint NOT NULL,
    http_status_codes char(3) [] NOT NULL,
    http_method smallint NOT NULL,
    body_encoding smallint NULL,
    request_body text NULL,
    request_headers text NULL,
    "group" bigint NULL,
    proxy bigint NULL,
    send_to_oncall boolean NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT HTTP_MONITORS_PK_1 PRIMARY KEY ("http_monitor_id")
);
-- Table: http_monitor_notifications
CREATE TABLE http_monitor_notifications (
    http_monitor_id bigint NOT NULL,
    notification_id bigint NOT NULL,
    CONSTRAINT HTTP_MONITOR_NOTIFICATIONS_PK PRIMARY KEY (http_monitor_id, notification_id)
);
-- Table: monitor_groups
CREATE TABLE monitor_groups (
    "monitor_group_id" bigint NOT NULL,
    team_id bigint NOT NULL,
    name varchar(64) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT MONITOR_GROUPS_PK_1 PRIMARY KEY ("monitor_group_id")
);
-- Table: notifications
CREATE TABLE notifications (
    notification_id bigint NOT NULL,
    user_id bigint,
    team_id bigint,
    notification_type smallint NOT NULL,
    notification_config jsonb NOT NULL,
    created_at timestamp NULL,
    updated_at timestamp NULL,
    CONSTRAINT NOTIFICATION_PK_1 PRIMARY KEY (notification_id)
);
-- Table: proxies
CREATE TABLE proxies (
    "proxy_id" bigint NOT NULL,
    CONSTRAINT PROXIES_PK_1 PRIMARY KEY ("proxy_id")
);
-- Table: incidents
CREATE TABLE incidents (
    incident_id bigint NOT NULL,
    team_id bigint NOT NULL,
    http_monitor_id bigint NULL,
    expiry_date timestamp NULL,
    confirm_location varchar(32) [] NULL,
    recover_location varchar(32) [] NULL,
    http_status_code char(3) NULL,
    incident_type smallint NOT NULL,
    incident_status smallint NOT NULL,
    incident_message text NOT NULL,
    notifications bool NOT NULL,
    incident_end timestamp,
    incident_start timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT INCIDENT_PK_1 PRIMARY KEY (incident_id)
);
-- Table: incident_timelines
CREATE TABLE incident_timelines (
    incident_timeline_id bigint NOT NULL,
    incident_id bigint NOT NULL,
    status_type smallint NOT NULL,
    message text NOT NULL,
    created_by bigint,
    server_id varchar(32),
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT INCIDENT_TIMELINES_PK_1 PRIMARY KEY (incident_timeline_id, incident_id)
);
-- Table: ping_data
CREATE TABLE ping_data (
    time TIMESTAMPTZ NOT NULL,
    http_monitor_id bigint NOT NULL,
    ping int NOT NULL,
    server_id varchar(32) NOT NULL
);
-- Create hypertable for ping_data table
SELECT create_hypertable('ping_data', 'time');
-- Add retention policy for ping_data table
SELECT add_retention_policy('ping_data', INTERVAL '26 hours');
-- Create materialized view for 15-minute aggregates
CREATE MATERIALIZED VIEW ping_data_15m WITH (timescaledb.continuous) AS
SELECT time_bucket('15 minutes', time) AS bucket,
    http_monitor_id,
    server_id,
    AVG(ping) AS avg_ping,
    MAX(ping) AS max_ping,
    MIN(ping) AS min_ping
FROM ping_data
GROUP BY bucket,
    http_monitor_id,
    server_id;
-- Add retention policy for ping_data_15m materialized view
SELECT add_retention_policy('ping_data_15m', INTERVAL '7 days');
-- Create materialized view for 2-hour aggregates
CREATE MATERIALIZED VIEW ping_data_2h WITH (timescaledb.continuous) AS
SELECT time_bucket('2 hours', time) AS bucket,
    http_monitor_id,
    server_id,
    AVG(ping) AS avg_ping,
    MAX(ping) AS max_ping,
    MIN(ping) AS min_ping
FROM ping_data
GROUP BY bucket,
    http_monitor_id,
    server_id;
-- Add retention policy for ping_data_2h materialized view
SELECT add_retention_policy('ping_data_2h', INTERVAL '90 days');
-- Create materialized view for 1-day aggregates
CREATE MATERIALIZED VIEW ping_data_1d WITH (timescaledb.continuous) AS
SELECT time_bucket('1 day', time) AS bucket,
    http_monitor_id,
    server_id,
    AVG(ping) AS avg_ping,
    MAX(ping) AS max_ping,
    MIN(ping) AS min_ping
FROM ping_data
GROUP BY bucket,
    http_monitor_id,
    server_id;
-- Add continuous aggregate policies
SELECT add_continuous_aggregate_policy(
        'ping_data_15m',
        start_offset => INTERVAL '60 minutes',
        end_offset => NULL,
        schedule_interval => INTERVAL '15 minutes'
    );
SELECT add_continuous_aggregate_policy(
        'ping_data_2h',
        start_offset => INTERVAL '6 hours',
        end_offset => NULL,
        schedule_interval => INTERVAL '2 hours'
    );
SELECT add_continuous_aggregate_policy(
        'ping_data_1d',
        start_offset => INTERVAL '2 days',
        end_offset => NULL,
        schedule_interval => INTERVAL '1 day'
    );
-- -------------------------------------------------------------
-- FOREIGN KEY
-- -------------------------------------------------------------
ALTER TABLE sessions
ADD CONSTRAINT session_fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id);
ALTER TABLE team_members
ADD CONSTRAINT team_members_fk_member_id FOREIGN KEY (member_id) REFERENCES users(user_id),
    ADD CONSTRAINT team_members_fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id);
ALTER TABLE team_invites
ADD CONSTRAINT team_invites_fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id),
    ADD CONSTRAINT team_invites_fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id);
ALTER TABLE http_monitors
ADD CONSTRAINT http_monitors_fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id),
    ADD CONSTRAINT http_monitors_fk_group FOREIGN KEY ("group") REFERENCES monitor_groups(monitor_group_id),
    ADD CONSTRAINT http_monitors_fk_proxy FOREIGN KEY (proxy) REFERENCES proxies(proxy_id);
ALTER TABLE http_monitor_notifications
ADD CONSTRAINT http_monitor_notifications_fk_http_monitor_id FOREIGN KEY (http_monitor_id) REFERENCES http_monitors(http_monitor_id),
    ADD CONSTRAINT http_monitor_notifications_fk_notification_id FOREIGN KEY (notification_id) REFERENCES notifications(notification_id);
ALTER TABLE monitor_groups
ADD CONSTRAINT monitor_groups_fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id);
ALTER TABLE notifications
ADD CONSTRAINT notifications_fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id),
    ADD CONSTRAINT notifications_fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id);
ALTER TABLE incidents
ADD CONSTRAINT incidents_fk_http_monitor_id FOREIGN KEY (http_monitor_id) REFERENCES http_monitors(http_monitor_id),
    ADD CONSTRAINT incidents_fk_team_id FOREIGN KEY (team_id) REFERENCES teams(team_id);
ALTER TABLE incident_timelines
ADD CONSTRAINT incident_timelines_fk_incident_id FOREIGN KEY (incident_id) REFERENCES incidents(incident_id),
    ADD CONSTRAINT incident_timelines_fk_server_id FOREIGN KEY (server_id) REFERENCES server_locations(server_id),
    ADD CONSTRAINT incident_timelines_fk_created_by FOREIGN KEY (created_by) REFERENCES users(user_id);
ALTER TABLE ping_data
ADD CONSTRAINT ping_data_fk_http_monitor_id FOREIGN KEY (http_monitor_id) REFERENCES http_monitors(http_monitor_id),
    ADD CONSTRAINT ping_data_fk_server_id FOREIGN KEY (server_id) REFERENCES server_locations(server_id);