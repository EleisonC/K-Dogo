package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EleisonC/K-Dogo/configs"
	"github.com/EleisonC/K-Dogo/models"
	"github.com/EleisonC/K-Dogo/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var newDog models.DogRequest
type resMessage struct {
	Message string `json:"message"`
	Count int64 `json:"count"`
}
var dogCollection *mongo.Collection = configs.GetCollection(configs.DB, "dogs")

var validate = validator.New()

func CreateDog(w http.ResponseWriter, r *http.Request) {
	//create context that this function will use
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// parse the dog data from the request to the newDog location.
	err := utils.ParseBody(r, &newDog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Print(err)
		return
	}
	if validationErr := validate.Struct((&newDog)); validationErr != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validationErr.Error()))
		return
	}

	// add the data to the mongo 
	newDogTime, err := utils.TimeParser(&newDog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Print(err)
		return
	}
	createDog := models.Dog{
		Name: newDog.Name,
		Breed: newDog.Breed,
		DateOfBirth: *newDogTime,
		Sex: newDog.Sex,
	}

	resultDog, err := dogCollection.InsertOne(ctx, createDog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err:=json.Marshal(resultDog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetDogById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context .WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := mux.Vars(r)
	dogId := params["dogId"]
	var dog models.Dog

	objId, err := primitive.ObjectIDFromHex(dogId)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	findErr := dogCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&dog)
	if findErr != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err:=json.Marshal(dog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetAllDogs(w http.ResponseWriter, r * http.Request) {
	ctx, cancel := context .WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var allDogs[] models.Dog
	results, err := dogCollection.Find(ctx, bson.M{})
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var dog models.Dog
		if err = results.Decode(&dog); err != nil {
			w.Header().Set("Content-Type", "pkglication/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return	
		}
		allDogs = append(allDogs, dog)
	}

	res, err := json.Marshal(allDogs)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteDog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context .WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := mux.Vars(r)

	dogId := params["dogId"]
	objId, err := primitive.ObjectIDFromHex(dogId)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	delResult, err := dogCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if delResult.DeletedCount == 0 {
		message := "The dog with ID " + dogId + " has not been deleted from the DB or does not exist"
		count := delResult.DeletedCount
		res := resMessage{
			Message: message,
			Count: count,
		}
		finalRes, err := json.Marshal(res)
		if err != nil {
			w.Header().Set("Content-Type", "pkglication/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(finalRes)
		return
	}
	
	message := "The dog with ID " + dogId + " has been deleted from the DB"
	delCont := delResult.DeletedCount

	res := resMessage{
		Message: message,
		Count: delCont,
	}
	finalRes, err := json.Marshal(res)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalRes)
}

func UpdateDog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := mux.Vars(r)
	dogId := params["dogId"]
	objId, err := primitive.ObjectIDFromHex(dogId)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = utils.ParseBody(r, &newDog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err = validate.Struct((&newDog)); err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	newDogTime, err := utils.TimeParser(&newDog)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Print(err)
		return
	}

	filter := bson.D{{Key: "_id", Value: objId}}
	update := bson.D{{Key:"$set", Value: bson.D{{Key: "name", Value: newDog.Name}, 
	{Key: "breed", Value: newDog.Breed}, {Key: "dateofbirth", Value: *newDogTime}, 
	{Key: "sex", Value: newDog.Sex}}}}


	updDog, err := dogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if updDog.ModifiedCount == 0 {
		message := "The dog with ID " + dogId + " has not updated in the DB or does not exist"
		count := updDog.ModifiedCount
		res := resMessage{
			Message: message,
			Count: count,
		}
		finalRes, err := json.Marshal(res)
		if err != nil {
			w.Header().Set("Content-Type", "pkglication/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(finalRes)
		return
	}

	message := "The dog with ID " + dogId + " has been updated in the DB"
	count := updDog.ModifiedCount

	res := resMessage{
		Message: message,
		Count: count,
	}

	finalRes, err := json.Marshal(res)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalRes)
}








