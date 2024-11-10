package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tensaidevX/go-projects/mongo-golang/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	//  Creates a filter to match a document that has the specified
	//  "_id" value

	result := models.User{}
	filter := bson.D{{"_id", oid}}
	opts := options.FindOne()

	coll := uc.client.Database("mongo-golang").Collection("users")
	err = coll.FindOne(context.TODO(), filter, opts).Decode(&result)

	if err != nil {
		fmt.Println(err)
	}
	result.ID = oid

	res, err := json.Marshal(result)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	u := models.User{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	coll := uc.client.Database("mongo-golang").Collection("users")

	result, err := coll.InsertOne(context.TODO(), u)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user_id": result.InsertedID,
	})

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		w.WriteHeader(404)
		return
	}
	oid, _ := primitive.ObjectIDFromHex(id)

	coll := uc.client.Database("mongo-golang").Collection("users")
	filter := bson.D{{"id", oid}}

	_, err := coll.DeleteOne(context.TODO(), filter)

	if err != nil {
		w.WriteHeader(404)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User", oid, "\n")

}
