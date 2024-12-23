package middleware

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Data      any    `json:"data"`
	Exception any    `json:"exception"`
	RequestID string `json:"requestId"`
	Time      string `json:"time"`
	Status    int    `json:"status"`
}

func Interceptor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		originalBody := c.Response().Body()
		requestID := c.GetRespHeader(fiber.HeaderXRequestID)

		response := Response{
			RequestID: requestID,
			Time:      time.Now().Format(time.RFC3339),
		}

		response.Status = c.Response().StatusCode()

		var parsedData interface{}
		if err := sonic.Unmarshal(originalBody, &parsedData); err == nil {
			response.Data = parsedData
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

		if response.Status > 399 {
			response.Exception = response.Data
			response.Data = nil
		}

		if err != nil {
			response.Exception = err.Error()

			if e, ok := err.(*fiber.Error); ok {
				if errors.Is(err, fiber.ErrUnprocessableEntity) {
					response.Exception = "No body found when expected"
				}

				response.Status = e.Code
			} else {
				if strings.Contains(err.Error(), "input json is empty") {
					response.Exception = "No body found when expected"

					response.Status = fiber.StatusUnprocessableEntity
				} else {
					response.Status = fiber.StatusInternalServerError
				}
			}
		}

		c.Locals("err", response.Exception)
		c.Status(response.Status)

		return c.JSON(response)
	}
}

func parseNumber(input string) (interface{}, error) {
	if intValue, err := strconv.ParseInt(input, 10, 64); err == nil {
		return intValue, nil
	}

	if floatValue, err := strconv.ParseFloat(input, 64); err == nil {
		return floatValue, nil
	}

	return nil, fmt.Errorf("not a valid number")
}
