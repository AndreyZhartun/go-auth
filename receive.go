package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//receive  - первый маршрут, пока без обращения к БД
func receive(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//log.Println(err)
		return
	}

	collection := client.Database("goauthtest").Collection("users")

	ac := Credentials{"admin", "admin", nil}
	_, err = collection.InsertOne(context.TODO(), ac)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//log.Println(err)
		return
	}
	filter := bson.M{"username": creds.Username}
	var result Credentials

	//TODO: фикс сравнение чистых паролей
	err = collection.FindOne(ctx, filter).Decode(&result)
	fmt.Printf("%v; using %s comparing true %s and given %s", err, creds.Username, result.Password, creds.Password)
	if err != nil || result.Password != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	/*res1, err := collection.InsertOne(ctx, bson.M{"name": "user1", "pass": "pass1"})
	fmt.Println(res1.InsertedID)
	res2, err := collection.InsertOne(ctx, bson.M{"name": "user2", "pass": "pass2"})
	fmt.Println(res2.InsertedID)*/

	//`Access токен тип JWT, алгоритм SHA512.`
	atExpiration := time.Now().Add(5 * time.Minute)
	atClaims := &сlaims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
			//jti для связи токенов в паре
			Id: "test",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//`Refresh токен тип произвольный (jwt), формат передачи base64`
	rtExpiration := time.Now().Add(5 * time.Hour)
	rtClaims := &сlaims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpiration.Unix(),
			Id:        "test",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtSigned, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: решить где хранить
	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Value:    atSigned,
		Expires:  atExpiration,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Value:    rtSigned,
		Expires:  rtExpiration,
		HttpOnly: true,
	})
}
