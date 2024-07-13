CREATE TABLE links
(
    id          SERIAL PRIMARY KEY,
    trip_id     INTEGER      NOT NULL REFERENCES trips (id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    url         TEXT         NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_trip_links ON links (trip_id);