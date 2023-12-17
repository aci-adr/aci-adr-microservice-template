package common

import (
	"FxService/model/dto/response"
	"FxService/model/validation"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
)

func ParamValidationMiddleware[T any](validationRules []validation.ValidationRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errors []response.Error
		for _, rule := range validationRules {
			paramValue := c.DefaultQuery(rule.ParamName, "")

			if rule.Required && paramValue == "" {
				errors = append(errors, response.Error{
					Code:    "INVALID_INPUT",
					Message: fmt.Sprintf("%s is required", rule.ParamName),
					Details: fmt.Sprintf("%s is required", rule.ParamName),
				})
				continue
			}

			if rule.ParamType == "int" {
				_, err := strconv.Atoi(paramValue)
				if err != nil {
					errors = append(errors, response.Error{
						Code:    "INVALID_INPUT",
						Message: "Data type invalid",
						Details: fmt.Sprintf("%s must be an integer", rule.ParamName),
					})
				}
			} else if rule.ParamType == "float" {
				_, err := strconv.ParseFloat(paramValue, 64)
				if err != nil {
					errors = append(errors, response.Error{
						Code:    "INVALID_INPUT",
						Message: "Data type invalid",
						Details: fmt.Sprintf("%s must be an float", rule.ParamName),
					})
				}
			} else if rule.ParamType == "date" {
				_, err := time.Parse(time.RFC3339, paramValue)
				if err != nil {
					errors = append(errors, response.Error{
						Code:    "INVALID_INPUT",
						Message: "Data type invalid",
						Details: fmt.Sprintf("%s must be an date", rule.ParamName),
					})
				}
			} else if rule.ParamType == "string" {
				// You can add more string-specific validation here
			}
		}
		if len(errors) > 0 {
			c.IndentedJSON(http.StatusOK, GetSimpleResponse[T](nil, response.BadRequest, &errors))
			c.Abort()
			return
		}

		// Continue to the next middleware or route handler
		c.Next()
	}
}

func ValidateAndReturnBody[T any](c *gin.Context) (T, error) {
	var data T
	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusOK, GetSimpleResponse[T](nil, response.BadRequest, translateError(err)))
		c.Abort()
		return data, err
	}
	return data, nil
}

func translateError(err error) *[]response.Error {
	var errors []response.Error
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, response.Error{
					Code:    "INVALID_INPUT",
					Message: fmt.Sprintf("The field %s is required.", e.Field()),
					Details: fmt.Sprintf("The field %s is required. it %s is not supplied", e.Field(), e.Param()),
				})
			}
		}
	} else if marshallingErr, ok := err.(*json.UnmarshalTypeError); ok {
		errors = append(errors, response.Error{
			Code:    "INVALID_INPUT",
			Message: "Error converting data type",
			Details: fmt.Sprintf("The field %s must be a %s", marshallingErr.Field, marshallingErr.Type.String()),
		})
	}
	return &errors
}

func ParamValidationMiddlewareFiber[T any](validationRules []validation.ValidationRule) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errors []response.Error
		for _, rule := range validationRules {
			paramValue := c.Query(rule.ParamName, "")

			if rule.Required && paramValue == "" {
				errors = append(errors, response.Error{
					Code:    "INVALID_INPUT",
					Message: fmt.Sprintf("%s is required", rule.ParamName),
					Details: fmt.Sprintf("%s is required", rule.ParamName),
				})
				continue
			}

			if rule.ParamType == "int" {
				_, err := strconv.Atoi(paramValue)
				if err != nil {
					errors = append(errors, response.Error{
						Code:    "INVALID_INPUT",
						Message: "Data type invalid",
						Details: fmt.Sprintf("%s must be an integer", rule.ParamName),
					})
				}
			} else if rule.ParamType == "float" {
				_, err := strconv.ParseFloat(paramValue, 64)
				if err != nil {
					errors = append(errors, response.Error{
						Code:    "INVALID_INPUT",
						Message: "Data type invalid",
						Details: fmt.Sprintf("%s must be an float", rule.ParamName),
					})
				}
			} else if rule.ParamType == "date" {
				_, err := time.Parse(time.RFC3339, paramValue)
				if err != nil {
					errors = append(errors, response.Error{
						Code:    "INVALID_INPUT",
						Message: "Data type invalid",
						Details: fmt.Sprintf("%s must be an date", rule.ParamName),
					})
				}
			} else if rule.ParamType == "string" {
				// You can add more string-specific validation here
			}
		}
		if len(errors) > 0 {
			return c.Status(fiber.StatusOK).JSON(GetSimpleResponse[T](nil, response.BadRequest, &errors))
		}

		// Continue to the next middleware or route handler
		return c.Next()
	}
}
