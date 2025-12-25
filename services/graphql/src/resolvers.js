import { createBooking, confirmBooking } from './grpc/bookingclient.js';

export const resolvers = {
  Query: {
    health: () => 'ok',
  },

  Mutation: {
    createBooking: async (_, args) => {
      const res = await createBooking(args);
      return {
        bookingId: res.booking_id,
        status: 'CREATED',
      };
    },

    confirmBooking: async (_, { bookingId }) => {
      await confirmBooking({ bookingId });
      return true;
    },
  },
};
