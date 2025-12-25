CREATE TABLE IF NOT EXISTS bookings (
    booking_id TEXT PRIMARY KEY,

    user_id    TEXT NOT NULL,
    event_id   TEXT NOT NULL,

    seat_ids   TEXT[] NOT NULL,

    status     TEXT NOT NULL,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE bookings
ADD CONSTRAINT bookings_status_check
CHECK (status IN (
    'CREATED',
    'PENDING_PAYMENT',
    'CONFIRMED',
    'EXPIRED',
    'CANCELLED'
));
