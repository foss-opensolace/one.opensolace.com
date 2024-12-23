package middleware

import (
	"fmt"
	"strconv"
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

		code := c.Response().StatusCode()
		originalBody := c.Response().Body()
		requestID := c.GetRespHeader(fiber.HeaderXRequestID)

		response := Response{
			RequestID: requestID,
			Time:      time.Now().Format(time.RFC3339),
			Status:    code,
		}

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

		if code > 399 {
			response.Exception = response.Data
			response.Data = nil
		}

		if err != nil {
			response.Exception = err.Error()

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				c.Status(e.Code)
			} else {
				c.Status(fiber.StatusInternalServerError)
			}

			c.Locals("err", err.Error())
		}

		response.Status = code

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
