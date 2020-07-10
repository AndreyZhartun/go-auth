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

func remove(w http.ResponseWriter, r *http.Request) {
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	collection := client.Database("goauthtest").Collection("users")
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

	var rtIndex int = -1
	for i, hash := range result.Rts {
		err = bcrypt.CompareHashAndPassword(hash, []byte(rtString))
		if err == nil {
			rtIndex = i
			break
		}
	}
	if rtIndex == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newRts := make([]([]byte), len(result.Rts), len(result.Rts))
	copy(newRts, result.Rts)
	newRts[rtIndex] = newRts[len(newRts)-1]
	newRts[len(newRts)-1] = []byte("")
	newRts = newRts[:len(newRts)-1]
	newUser := User{GUID: rtClaims.UserID, Rts: newRts}

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
