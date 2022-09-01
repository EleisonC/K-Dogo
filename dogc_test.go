package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EleisonC/K-Dogo/routes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)


func Router() *mux.Router{
	// create router
	router := mux.NewRouter()
	// register routes
	routes.RegisterDogRoutes(router)
	return router
}

func TestGetDogs(t *testing.T) {
	request, _ := http.NewRequest("GET", "/getdogs", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestCreateDogs(t *testing.T) {
	type args struct {
		Name string `json:"name,omitempty"`
		Breed string `json:"breed,omitempty"`
		Dateofbirth string `json:"dateofbirth,omitempty"`
		Sex string `json:"sex,omitempty"`
	}
	type expectedRes struct{
		message string
		expectedStatus int
	}
	tests := []struct {
		name string
		args args
		expectation expectedRes
	}{
		{
			name: "dog added successfully",
			args: args {
				Name: "Jipsy",
				Breed: "GSD",
				Dateofbirth: "2021-06-30",
				Sex: "Female",
			},
			expectation: expectedRes{
				message: "Data created successfully",
				expectedStatus: 200,
			},
		},
		{
			name: "invalid data types",
			args: args {
				Name: "Eleven",
				Breed: "GSD",
				Dateofbirth: "2021-06-30",
			},
			expectation: expectedRes{
				message: "Data created successfully",
				expectedStatus: 400,
			},

		},
	} 
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bites, err := json.Marshal(tt.args)
			
			if err != nil {
				t.Error(err)
			}
			m := bytes.NewBuffer(bites)
			request, _ := http.NewRequest("POST", "/createdog", m)
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			assert.Equal(t, tt.expectation.expectedStatus, response.Code)
		})
	 }
}

