package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

func Helmet() fiber.Handler {
	return helmet.New()
}

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format: `
[${green}Client${reset}]:     ${ip}:${port}
[${green}Timestamp${reset}]:  ${time}
[${green}Rec/Sent${reset}]:   ${bytesReceived}/${bytesSent}
[${green}Process ID${reset}]: ${pid}
[${green}Request${reset}]:    ${status} ${method} ${yellow}${url}${reset} via ${magenta}${ua}${reset}
[${green}Request ID${reset}]: ${respHeader:X-Request-ID}

[${red}Latency${reset}]: ${latency}
[${red}Error${reset}]:   ${locals:err}

[${blue}API KEY${reset}]: ${black}${reqHeader:X-API-KEY}${reset}
[${blue}Request Body${reset}]:
${black}${body}${reset}
`,
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	})
}

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowHeaders:     "X-API-KEY, Content-Type, Authorization",
		AllowMethods:     "GET, POST, PATCH, PUT, DELETE",
	})
}

func RequestId() fiber.Handler {
	return requestid.New()
}
