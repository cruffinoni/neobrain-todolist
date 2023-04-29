package utils

import (
	"log"
)

type InternalServerErrorBuilder struct {
	Error string `json:"error"`
}

func NewInternalServerErrorBuilder(err error) *InternalServerErrorBuilder {
	log.Printf("internal server error: %v", err)
	return &InternalServerErrorBuilder{
		Error: "internal server error",
	}
}

type BadRequestBuilder struct {
	Error string `json:"error"`
}

func NewBadRequestBuilder(error string) *BadRequestBuilder {
	return &BadRequestBuilder{Error: error}
}

type StatusOKBuilder struct {
	Message string `json:"message"`
}

func NewStatusOKBuilder(message string) *StatusOKBuilder {
	return &StatusOKBuilder{Message: message}
}
