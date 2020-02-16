CREATE KEYSPACE IF NOT EXISTS farseer
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE farseer;

CREATE TABLE IF NOT EXISTS events (
    user_id text,
    item_id text,
    event_type text,
    event_value int,
    version_timestamp timeuuid,
    event_timestamp timestamp,
    properties map<text, text>,
    PRIMARY KEY (user_id, version_timestamp, event_type, item_id)
) WITH CLUSTERING ORDER BY (version_timestamp DESC);

CREATE TABLE IF NOT EXISTS items (
    id text,
    item_type text,
    properties map<text, text>,
    created_time timestamp,
    updated_time timestamp,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
    id text,
    user_type text,
    properties map<text, text>,
    created_time timestamp,
    updated_time timestamp,
    PRIMARY KEY (id)
);
