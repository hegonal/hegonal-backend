CREATE TABLE server_locations (
    server_id varchar(32) NOT NULL,
    server_display_name varchar(32) NOT NULL,
    country char(3),
    CONSTRAINT SERVER_LOCATIONS_PK_1 PRIMARY KEY (server_id)
);

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

ALTER TABLE
    sessions
ADD
    CONSTRAINT SESSIONS_FK_1 FOREIGN KEY (user_id) REFERENCES users(user_id);

CREATE TABLE teams (
    "team_id" bigint NOT NULL,
    name varchar(64) NOT NULL,
    description varchar(128) NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT TEAMS_PK_1 PRIMARY KEY ("team_id")
);

CREATE TABLE team_members (
    member_id bigint NOT NULL,
    team_id bigint NOT NULL,
    role smallint NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT TEAM_MEMBERS_PK_1 PRIMARY KEY (member_id, team_id)
);

ALTER TABLE
    team_members
ADD
    CONSTRAINT TEAM_MEMBERS_FK_1 FOREIGN KEY (member_id) REFERENCES users(user_id),
ADD
    CONSTRAINT TEAM_MEMBERS_FK_2 FOREIGN KEY (team_id) REFERENCES teams(team_id);

CREATE TABLE http_monitors (
    "http_monitor_id" bigint NOT NULL,
    team_id bigint NOT NULL,
    status smallint NOT NULL,
    url varchar(255) NOT NULL,
    "interval" integer NOT NULL,
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
    notification bigint NULL,
    proxy bigint NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT HTTP_MONITORS__PK_1 PRIMARY KEY ("http_monitor_id")
);

ALTER TABLE
    http_monitors
ADD
    CONSTRAINT HTTP_MONITORS_FK_1 FOREIGN KEY (team_id) REFERENCES teams(team_id),
ADD
    CONSTRAINT HTTP_MONITORS_FK_2 FOREIGN KEY ("group") REFERENCES monitor_groups(group_id),
ADD
    CONSTRAINT HTTP_MONITORS_FK_3 FOREIGN KEY (notification) REFERENCES notifications(notification_id),
ADD
    CONSTRAINT HTTP_MONITORS_FK_4 FOREIGN KEY (proxy) REFERENCES proxies(proxy_id);

CREATE TABLE monitor_groups (
    "monitor_group_id" bigint NOT NULL,
    team_id bigint NOT NULL,
    name varchar(64) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT MONITOR_GROUPS_PK_1 PRIMARY KEY ("monitor_group_id")
);

ALTER TABLE
    monitor_groups
ADD
    CONSTRAINT MONITOR_GROUPS_FK_1 FOREIGN KEY (team_id) REFERENCES teams(team_id);

CREATE TABLE notifications (
    notification_id bigint NOT NULL,
    notification_type smallint NOT NULL,
    notification_config jsonb NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    CONSTRAINT NOTIFICATION_PK_1 PRIMARY KEY (notification_id)
);

CREATE TABLE proxies (
    "proxy_id" bigint NOT NULL,
    CONSTRAINT PROXIES_PK_1 PRIMARY KEY ("proxy_id")
);

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

ALTER TABLE
    incidents
ADD
    CONSTRAINT INCIDENTS_FK_1 FOREIGN KEY (http_monitor_id) REFERENCES http_monitors(http_monitor_id),
ADD
    CONSTRAINT INCIDENTS_FK_2 FOREIGN KEY (team_id) REFERENCES teams(team_id);

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

ALTER TABLE
    incident_timelines
ADD
    CONSTRAINT INCIDENT_TIMELINES_FK_1 FOREIGN KEY (incident_id) REFERENCES incidents(incident_id),
ADD
    CONSTRAINT INCIDENT_TIMELINES_FK_2 FOREIGN KEY (server_id) REFERENCES server_locations(server_id),
ADD
    CONSTRAINT INCIDENT_TIMELINES_FK_3 FOREIGN KEY (created_by) REFERENCES users(user_id);

CREATE TABLE http_monitor_long_history (
    http_monitor_id bigint NOT NULL,
    created_at timestamp NOT NULL,
    ping smallint NOT NULL,
    CONSTRAINT HTTP_MONITOR_LONG_HISTORY_PK_1 PRIMARY KEY (monitor_id, created_at)
);

ALTER TABLE
    http_monitor_long_history
ADD
    CONSTRAINT HTTP_MONITOR_LONG_HISTORY_FK_1 FOREIGN KEY (monitor_id) REFERENCES http_monitors(http_monitor_id);

CREATE TABLE http_monitor_short_history (
    http_monitor_id bigint NOT NULL,
    created_at timestamp NOT NULL,
    ping smallint NOT NULL,
    CONSTRAINT HTTP_MONITOR_SHORT_HISTORY_PK_1 PRIMARY KEY (monitor_id, created_at)
);

ALTER TABLE
    http_monitor_short_history
ADD
    CONSTRAINT HTTP_MONITOR_SHORT_HISTORY_FK_1 FOREIGN KEY (monitor_id) REFERENCES http_monitors(http_monitor_id);