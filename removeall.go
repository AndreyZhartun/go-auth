package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func removeAll(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("rt")
	rtString := c.Value
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	rtClaims, err := getClaims(c)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//TODO: в одну транзакцию
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		connString,
	))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	collection := client.Database("goauth").Collection("users")
	findFilter := bson.M{"guid": rtClaims.UserID}
	var result User

	err = collection.FindOne(ctx, findFilter).Decode(&result)
	if err != nil {
		fmt.Printf("collection find err: %v\n", err)
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	var rtExists bool
	for _, hash := range result.Rts {
		err = bcrypt.CompareHashAndPassword(hash, []byte(rtString))
		if err == nil {
			rtExists = true
			break
		}
	}
	if !rtExists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newUser := User{GUID: rtClaims.UserID, Rts: make([]([]byte), 0, 0)}

	updateFilter := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "rts", Value: newUser.Rts}}}}

	_, err = collection.UpdateOne(ctx, findFilter, updateFilter)
	if err != nil {
		fmt.Printf("update err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
