package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// `Первый маршрут выдает пару Access, Refresh токенов
// для пользователя с идентификатором (GUID) указанным в параметре запроса`
func receive(w http.ResponseWriter, r *http.Request) {

	pars, ok := r.URL.Query()["guid"]
	if !ok || len(pars[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		//log.Println("Url Param 'key' is missing")
		return
	}
	var guid string = pars[0]
	//azadmin: saJgPQmeekDwE5S
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
	//два соединения??
	/*if err = client.Ping(ctx, readpref.Primary()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//log.Println(err)
		return
	}*/

	collection := client.Database("goauth").Collection("users")
	findFilter := bson.M{"guid": guid}
	var result User

	err = collection.FindOne(ctx, findFilter).Decode(&result)
	isNewUser := err == mongo.ErrNoDocuments
	//fmt.Printf("%v; %s -> %v\n", err, result.GUID, result.Rts)
	if err != nil && !isNewUser {
		fmt.Printf("collection find err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//`Refresh токен тип произвольный (jwt), формат передачи base64`
	rtExpiration := time.Now().Add(5 * time.Hour)
	rtClaims := &сlaims{
		UserID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpiration.Unix(),
			Id:        "test",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtSigned, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("rt err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rtHashed, err := bcrypt.GenerateFromPassword([]byte(rtSigned), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("rt bcrypt err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("inserting: %s\n", rtSigned)
	if isNewUser {
		newRts := make([]([]byte), 0, 1)
		newRts = append(newRts, rtHashed)
		newUser := User{GUID: guid, Rts: newRts}

		_, err = collection.InsertOne(ctx, newUser)
		if err != nil {
			fmt.Printf("new user insert err: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		newRts := make([]([]byte), len(result.Rts), len(result.Rts)+1)
		copy(newRts, result.Rts)
		newRts = append(newRts, rtHashed)
		newUser := User{GUID: guid, Rts: newRts}
		//TODO без создания структуры тест
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

	//`Access токен тип JWT, алгоритм SHA512.`
	atExpiration := time.Now().Add(5 * time.Minute)
	atClaims := &сlaims{
		UserID: guid,
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
