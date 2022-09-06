package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrMessageRes struct {
	Message string `json:"message"`
	RawErrorMessage string `json:"raw err message"`
}

type ResMessage struct {
	Message string `json:"message"`
	Count int64 `json:"count"`
}

func ParseBody(r *http.Request, x interface{}) error{
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &x)
	if err != nil {
		return err
	}
	return nil
}

func TimeParser(s interface{}) (*time.Time, error){
	val := reflect.ValueOf(s).Elem()
	n := val.FieldByName("DateOfBirth").Interface().(string)
	isoLayout := "2006-1-2"
	t, err := time.Parse(isoLayout,n)
	if err != nil {
		return nil, errors.New("failed to parse date field:" + err.Error())
	}
	return &t, nil
}

func ValidateOwner(ownerId string, w http.ResponseWriter, ctx context.Context, dogOwnerCol *mongo.Collection) error {
	ownerIdHex, err := primitive.ObjectIDFromHex(ownerId)
	if err != nil {
		ErrorHandlerDogs(w, err, "Error Can Not Convert")
		return err
	}
	filter := bson.D{{Key: "_id", Value: ownerIdHex}}
	ownerIdCount, err := dogOwnerCol.CountDocuments(ctx, filter)
	if err != nil || ownerIdCount == 0 {
		ErrorHandlerDogs(w, err, "Error Un Known User")
		return err
	}
	return nil
}

func ErrorHandlerDogs(w http.ResponseWriter, err error, message string) {
	if err != nil {
		errMessage := ErrMessageRes {
			Message: message,
			RawErrorMessage: err.Error(),
		}
		errMes, _:= json.Marshal(errMessage)

		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMes)
	}
	return
}

