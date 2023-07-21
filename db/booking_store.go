package db

import (
	"context"
	
	"hotel_reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, Map) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, Map) error
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("bookings"),
	}
}

func (b *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update Map) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	u := bson.M{"$set": update}
	_, err = b.collection.UpdateByID(ctx, oid, u)
	return err
}

func (b *MongoBookingStore) GetBookings(ctx context.Context, filter Map) ([]*types.Booking, error) {
	res, err := b.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := res.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (b *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := b.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (b *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := b.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}
