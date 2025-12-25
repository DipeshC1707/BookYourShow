export const typeDefs = `#graphql
  type Booking {
    bookingId: String!
    status: String
  }

  type Query {
    health: String!
  }

  type Mutation {
    createBooking(
      eventId: String!
      seatIds: [String!]!
      userId: String!
      idempotencyKey: String!
    ): Booking!

    confirmBooking(
      bookingId: String!
    ): Boolean!
  }
`;
