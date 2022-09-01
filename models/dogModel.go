package models

import (
	"time"
)

// type DogVaccines struct {
// 	Name string `json:"name,omitempty" validate:"required"`
// 	VaccinePurpose string `json:"vaccinename,omitempty" validate:"required"`
// 	DueDate time.Time `json:"duedate,omitempty" validate:"required"`
// }

type Dog struct {
	Name string `bson:"name,omitempty" validate:"required"`
	Breed string `bson:"breed,omitempty" validate:"required"`
	DateOfBirth  time.Time `bson:"dateofbirth,omitempty" validate:"required"`
	Sex string `bson:"sex,omitempty" validate:"required"`
	TrainerId string `bson:"trainerId,omitempty"`
	// Vacination []DogVaccines `json:"vaccinations"`
}

type DogRequest struct {
	Name string `json:"name,omitempty" validate:"required"`
	Breed string `json:"breed,omitempty" validate:"required"`
	DateOfBirth  string `json:"dateofbirth,omitempty" validate:"required"`
	Sex string `json:"sex,omitempty" validate:"required"`
	TrainerId string `bson:"trainerId,omitempty"`
	// Vacination []DogVaccines `json:"vaccinations"`
}

type UpdateDogRequest struct {
	Name *string `json:"name,omitempty"`
	Breed *string `json:"breed,omitempty"`
	DateOfBirth  *time.Time `json:"dateofbirth,omitempty"`
	Sex *string `json:"sex,omitempty"`
	TrainerId string `bson:"trainerId,omitempty"`
	// Vacination []DogVaccines `json:"vaccinations"`
}

