package db

import (
	"context"
	"hotel_reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	collection *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client: client,
		collection: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (h *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	resp, err := h.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (h *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := h.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)
	// Update the hotel with this room id
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := h.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}
