package db

import (
	"context"
	"hotel_reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		collection: client.Database(DBNAME).Collection("hotels"),
	}
}

func (h *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	var hotel *types.Hotel
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := h.collection.FindOne(ctx, bson.M{"_id": oid})
	if err := res.Decode(&hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (h *MongoHotelStore) GetHotels(ctx context.Context, filter Map, pag *Pagination) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip(int64((pag.Page - 1) * pag.Limit))
	opts.SetLimit(int64(pag.Limit))
	resp, err := h.collection.Find(ctx, filter, &opts)
	if err != nil {
		return nil,  err
	}
	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (h *MongoHotelStore) Update(ctx context.Context, filter Map, update Map) error {
	_, err := h.collection.UpdateOne(ctx, filter, update)
	return err
}

func (h *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := h.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
