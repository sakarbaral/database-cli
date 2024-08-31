package models

import (
	"encoding/json"
)

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type Logger interface {
	Debug(string, ...interface{})
}

type Options struct {
	Logger
}
