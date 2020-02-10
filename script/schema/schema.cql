CREATE KEYSPACE IF NOT EXISTS farseer
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE farseer;

CREATE TABLE IF NOT EXISTS users (
    id text,
    properties map<text, text>,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS items (
    id text,
    properties map<text, text>,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS events (
    user_id text,
    item_id text,
    event_type text,
    event_value int,
    timestamp timeuuid,
    -- session text,
    properties map<text, text>,
    PRIMARY KEY (user_id, timestamp, event_type, item_id)
) WITH CLUSTERING ORDER BY (timestamp DESC);
