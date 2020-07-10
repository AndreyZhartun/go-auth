package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// jwt.StandardClaims для exp и jti
type сlaims struct {
	UserID string `json:"uid"`
	jwt.StandardClaims
}

func getClaims(c *http.Cookie) (*сlaims, error) {
	tknStr := c.Value
	claims := &сlaims{}

	// парс JWT в claims
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		// проверка, что токен подписан с SHA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//fmt.Printf("sign method unexpected: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("Неожиданный алгоритм подписи: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		/*if err == jwt.ErrSignatureInvalid {
			//fmt.Println("tk sign not valid")
			return nil, jwt.ErrSignatureInvalid
		}*/
		//fmt.Printf("parse err: %v", err)
		return nil, err
	}
	if !tkn.Valid {
		//fmt.Println("tk not valid")
		return claims, errors.New("Token is not valid")
	}

	return claims, nil
}

func validateWithDB(rt string, guid string) error {
	//проверка наличия refresh токена в БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		connString,
	))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			//w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	collection := client.Database("goauth").Collection("users")
	findFilter := bson.M{"guid": guid}
	var result User

	err = collection.FindOne(ctx, findFilter).Decode(&result)
	if err != nil {
		fmt.Printf("collection find err: %v\n", err)
		/*if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}*/
		//w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	var rtExists bool
	for _, hash := range result.Rts {
		err = bcrypt.CompareHashAndPassword(hash, []byte(rt))
		if err == nil {
			rtExists = true
			break
		}
	}
	if !rtExists {
		//w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Not found")
	}
	return nil
}

// `Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов`
func refresh(w http.ResponseWriter, r *http.Request) {
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

	c, err = r.Cookie("at")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	atClaims, err := getClaims(c)
	if err != nil && err.Error() != "Token is not valid" {
		fmt.Println("xd not vaild")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	/*	`Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена
		можно выполнить только тем Refresh токеном который был выдан вместе с ним.`	*/
	//TODO: генерация Id
	if rtClaims.UserID != atClaims.UserID || rtClaims.Id != atClaims.Id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//проверка наличия refresh токена в БД
	err = validateWithDB(rtString, rtClaims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	/*ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		return
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
	}*/

	//проверка токенов завершена
	atExpiration := time.Now().Add(5 * time.Minute)
	atClaims = &сlaims{
		UserID: atClaims.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
			//jti для связи токенов в паре
			Id: rtClaims.Id,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Value:    atSigned,
		Expires:  atExpiration,
		HttpOnly: true,
	})
}
