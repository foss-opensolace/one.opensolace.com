package exception

import (
	"github.com/gofiber/fiber/v2"
)

var exceptionIdTag = "exception_id"

type Exception uint16

/*
Errors explained:

Layout: CXZZ

C - Unknown (1) | Internal (2) | Validation (3) | Auth (4) | User (5) | API Key (6)
X - Unknown (0) | Non-Critical (1) | Critical (2);
Z - Error code
*/
const (
	Unknown Exception = 1000

	OneOrManyValidation Exception = 2101 /* Example: Expected: n >= 1 | Error: 0 */
	InvalidFieldType    Exception = 2102 /* Example: Expected: 1 | Error: "1" */
	InvalidFieldLayout  Exception = 2103 /* Example: Expected: "2006-10-20" | Error: "2006-20-10" */

	AuthInvalidCredentials Exception = 4101
	AuthDuplicated         Exception = 4201

	UserNotFound Exception = 5101

	MissingAPIKeyHeader     Exception = 6101
	InvalidAPIKey           Exception = 6102
	CannotUseAPIKey         Exception = 6103
	APIKeyMissingPermission Exception = 6104
)

func GetID(c *fiber.Ctx) *Exception {
	exception := c.Locals(exceptionIdTag)

	if exception != nil {
		id := exception.(Exception)
		return &id
	}

	return nil
}

func SetID(c *fiber.Ctx, i Exception) {
	c.Locals(exceptionIdTag, i)
}
