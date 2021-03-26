package Middleware

import (
	"clone_project/Models"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Email string
	Password string
}

var collection *mongo.Collection
const connectionString = "mongodb://localhost:27017/"
const dataBase = "test"
const collectionName = "sales"
var ctx = context.TODO()

func init(){
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DataBase Connected!!")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("test").Collection("user")
}

func GetLandingPage( w http.ResponseWriter , r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	str:= "Server is up and Running!!!"
	json.NewEncoder(w).Encode(str)
}

func GetAllUsers( w http.ResponseWriter , r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type","application/json")
	findOptions := options.Find()
	var data []*Models.UserModel
	response, err := collection.Find(context.Background() , bson.D{{}} , findOptions)
	if err != nil{
		log.Fatal(err)
	}
	for response.Next(context.TODO()) {
		var elem Models.UserModel
		err := response.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, &elem)
	}
	if err := response.Err(); err != nil {
		log.Fatal(err)
	}

	response.Close(context.TODO())
	json.NewEncoder(w).Encode(data)
}

func CreateUser(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var user Models.UserModel
	reqBody, err := ioutil.ReadAll(r.Body)
	if err!= nil {
		log.Fatal(err)
	}
	json.Unmarshal(reqBody , &user)
	var password string
	password = user.Password
	hash := HashPassword(password)
	if err!= nil {
		log.Fatal(err)
	}
	var new_user Models.UserModel
	new_user.UserName = user.Email
	new_user.Email = user.Email
	new_user.Password = hash
	user_obj , err := collection.InsertOne(context.Background(), new_user)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(user_obj)
}

func GetUser( w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	var user User
	reqBody , err := ioutil.ReadAll(r.Body)
	if err!= nil {
		log.Fatal(err)
	}
	json.Unmarshal(reqBody ,&user)
	var check_user Models.UserModel
	filter := bson.D{{"email", user.Email }}
	var error  = collection.FindOne(context.TODO(), filter).Decode(&check_user)
	if error!= nil {
		log.Fatal(error)
	}
	match := CheckPasswordHash(user.Password , check_user.Password)

	if match == true {
		json.NewEncoder(w).Encode(check_user)
	}
}

func HashPassword( password string) (string){
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string ) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetNearPoints(w http.ResponseWriter , r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var user User
	var check_user Models.UserModel
	reqBody , err := ioutil.ReadAll(r.Body)
	if err!= nil {
		log.Fatal(err)
	}
	json.Unmarshal(reqBody, &user)
	user_filter := bson.D{{"email" ,  user.Email }}
	err1 := collection.FindOne(context.Background(), user_filter).Decode(&check_user)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}
}