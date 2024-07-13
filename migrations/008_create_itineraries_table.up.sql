CREATE TABLE IF NOT EXISTS itineraries
(
    id          SERIAL PRIMARY KEY,
    trip_id     INTEGER                  NOT NULL REFERENCES trips (id) ON DELETE CASCADE,
    title       VARCHAR(100)             NOT NULL,
    description TEXT,
    date        DATE                     NOT NULL,
    created_by  INTEGER                  NOT NULL REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_itineraries_trip_id ON itineraries (trip_id);