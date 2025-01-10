package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/foss-opensolace/api.opensolace.com/pkg/exception"
	"github.com/foss-opensolace/api.opensolace.com/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
)

type Response struct {
	Data        any                  `json:"data"`
	Exception   any                  `json:"exception"`
	ExceptionID *exception.Exception `json:"exception_id"`
	RequestID   string               `json:"request_id"`
	IssuedAt    string               `json:"issued_at"`
	Status      int                  `json:"status"`
}

func Interceptor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		originalBody := c.Response().Body()
		requestID := c.GetRespHeader(fiber.HeaderXRequestID)

		response := Response{
			RequestID:   requestID,
			ExceptionID: exception.GetID(c),
			IssuedAt:    time.Now().Format(time.RFC3339),
		}

		response.Status = c.Response().StatusCode()

		var parsedData any
		if err := json.Unmarshal(originalBody, &parsedData); err == nil {
			response.Data = &parsedData
		} else {
			if len(originalBody) > 0 {
				bodyString := string(originalBody)

				switch string(originalBody) {
				case "true", "false":
					response.Data = string(originalBody) == "true"

				default:
					if number, err := parseNumber(bodyString); err == nil {
						response.Data = number
					} else {
						response.Data = bodyString
					}
				}
			}
		}

		if err != nil {
			response.Data = err.Error()

			if e, ok := err.(*fiber.Error); ok {
				if errors.Is(err, fiber.ErrUnprocessableEntity) {
					response.Data = "No body found when expected"
				}

				response.Status = e.Code
			} else if e, ok := err.(exception.FieldTypeError); ok {
				response.ExceptionID = utils.ToPtr(exception.IdInvalidFieldType)
				response.Data = e.Error()

				response.Status = fiber.StatusBadRequest
			} else if e, ok := err.(exception.FieldLayoutError); ok {
				response.ExceptionID = utils.ToPtr(exception.IdInvalidFieldLayout)
				response.Data = e.Error()

				response.Status = fiber.StatusBadRequest
			} else if e, ok := err.(*multierror.Error); ok {
				response.ExceptionID = utils.ToPtr(exception.IdOneOrManyValidation)
				response.Data = e.Errors

				response.Status = fiber.StatusBadRequest
			} else if strings.Contains(err.Error(), "input json is empty") {
				response.Data = "No body found when expected"

				response.Status = fiber.StatusBadRequest
			} else {
				response.Status = fiber.StatusInternalServerError
			}
		}

		if response.Status > 399 {
			response.Exception = response.Data
			response.Data = nil

			if response.ExceptionID == nil {
				response.ExceptionID = utils.ToPtr(exception.IdUnknown)
			}
		}

		c.Locals("err", response.Exception)
		c.Status(response.Status)

		return c.JSON(response)
	}
}

func parseNumber(input string) (any, error) {
	if intValue, err := strconv.ParseInt(input, 10, 64); err == nil {
		return intValue, nil
	}

	if floatValue, err := strconv.ParseFloat(input, 64); err == nil {
		return floatValue, nil
	}

	return nil, fmt.Errorf("not a valid number")
}
