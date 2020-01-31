CREATE KEYSPACE IF NOT EXISTS compass
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE compass;

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
    timestamp timeuuid,
    -- session text,
    event text,
    item_id text,
    properties map<text, text>,
    PRIMARY KEY (user_id, timestamp, event, item_id)
) WITH CLUSTERING ORDER BY (timestamp DESC);
