CREATE TABLE expenses
(
    id          SERIAL PRIMARY KEY,
    trip_id     INTEGER                  NOT NULL REFERENCES trips (id),
    category    VARCHAR(50)              NOT NULL,
    amount      DECIMAL(10, 2)           NOT NULL,
    description TEXT,
    date        DATE                     NOT NULL,
    created_by  INTEGER                  NOT NULL REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_expenses_trip_id ON expenses (trip_id);