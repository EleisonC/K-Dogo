package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
	// "reflect"
	// "time"
)

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