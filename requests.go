package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()
var postValidate = validator.New()

func DecodeForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(dst, r.PostForm); err != nil {
		return err
	}

	return nil
}

func ValidatePostData(data interface{}) map[string]string {
	err := postValidate.Struct(data)
	postErrors := map[string]string{}
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			postErrors["system"] = "error"
		}

		for _, err := range err.(validator.ValidationErrors) {
			ferr := err.ActualTag()
			if ferr == "gte" {
				ferr = "Must be greater than '" + err.Param() + "' characters"
			} else if ferr == "lte" {
				ferr = "Must be less than '" + err.Param() + "' characters"
			} else if ferr == "required" {
				ferr = "Required"
			} else if ferr == "email" {
				ferr = "Must be an Email"
			} else if ferr == "numeric" {
				ferr = "Must contain only numbers"
			} else if ferr == "alphaunicode" {
				ferr = "No special characters"
			} else if ferr == "alphanum" {
				ferr = "No special characters"
			} else if ferr == "datetime" {
				ferr = "Must be in datetime format"
			} else if ferr == "boolean" {
				ferr = "Must be true or false"
			}
			postErrors[err.Field()] = ferr
		}
	}
	return postErrors
}
