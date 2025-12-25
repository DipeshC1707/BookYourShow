CREATE TABLE IF NOT EXISTS booking_seats (
    booking_id TEXT NOT NULL,
    event_id   TEXT NOT NULL,
    seat_id    TEXT NOT NULL,
    status     TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_confirmed_seat
ON booking_seats (event_id, seat_id)
WHERE status = 'CONFIRMED';