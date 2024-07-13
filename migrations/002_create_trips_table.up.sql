CREATE TABLE IF NOT EXISTS trips
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255)             NOT NULL,
    description TEXT,
    start_date  DATE                     NOT NULL,
    end_date    DATE                     NOT NULL,
    created_by  INTEGER                  NOT NULL REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_trips_created_by ON trips (created_by);