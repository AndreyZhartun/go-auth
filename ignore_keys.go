package main

import "os"

//var connString = "mongodb+srv://azadmin:saJgPQmeekDwE5S@az-cluster.9qxft.mongodb.net/goauth?retryWrites=true&w=majority"
var connString = string(os.Getenv("CONNSTR"))

//var jwtKey = []byte("why-you-have-to-be-mad")
var jwtKey = []byte(string(os.Getenv("SECRET")))
