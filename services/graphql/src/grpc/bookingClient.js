import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

const packageDef = protoLoader.loadSync(
  '../../proto/booking/v1/booking.proto',
  { keepCase: true }
);

const bookingProto = grpc.loadPackageDefinition(packageDef).booking.v1;

const client = new bookingProto.BookingService(
  process.env.BOOKING_GRPC_ADDR || 'localhost:50052',
  grpc.credentials.createInsecure()
);

export function createBooking({ eventId, seatIds, userId, idempotencyKey }) {
  return new Promise((resolve, reject) => {
    client.CreateBooking(
      { event_id: eventId, seat_ids: seatIds, user_id: userId, idempotency_key: idempotencyKey },
      (err, response) => {
        if (err) return reject(err);
        resolve(response);
      }
    );
  });
}

export function confirmBooking({ bookingId }) {
  return new Promise((resolve, reject) => {
    client.ConfirmBooking(
      { booking_id: bookingId },
      (err) => {
        if (err) return reject(err);
        resolve(true);
      }
    );
  });
}
