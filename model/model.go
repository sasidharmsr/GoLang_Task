package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// My User Model
type User struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name"`
	DOB         time.Time          `json:"dob"`
	Address     Address            `json:"address"`
	Description string             `json:"description"`
	CreatedAt   primitive.DateTime `json:"created_at"`
	Following   []User             `json:"following"`
	Followers   []User             `json:"followers"`
}

// Model For Returning Response
type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type Address struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Addresss  string  `json:"address"`
}
