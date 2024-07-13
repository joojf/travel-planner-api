CREATE TABLE IF NOT EXISTS expenses
(
    id          SERIAL PRIMARY KEY,
    trip_id     INTEGER                  NOT NULL REFERENCES trips (id),
    amount      DECIMAL(10, 2)           NOT NULL,
    description TEXT,
    category    VARCHAR(100),
    paid_by     INTEGER                  NOT NULL REFERENCES users (id),
    date        DATE                     NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_expenses_trip_id ON expenses (trip_id);
CREATE INDEX idx_expenses_paid_by ON expenses (paid_by);