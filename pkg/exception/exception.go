package exception

import (
	"github.com/gofiber/fiber/v2"
)

var exceptionIdTag = "exception_id"

type Exception uint16

/*
Errors explained:

Layout: CXZZ

C - Unknown (1) | Internal (2) | Generic (3) | Validation (4) | Auth (5) | API Key (6)
X - Unknown (0) | Non-Critical (1) | Critical (2);
Z - Error code
*/
const (
	Unknown Exception = 1000

	GenericNotFound Exception = 3101

	OneOrManyValidation Exception = 4101 /* Example: Expected: n >= 1 | Error: 0 */
	InvalidFieldType    Exception = 4102 /* Example: Expected: 1 | Error: "1" */
	InvalidFieldLayout  Exception = 4103 /* Example: Expected: "2006-10-20" | Error: "2006-20-10" */

	AuthInvalidCredentials Exception = 5101
	AuthDuplicated         Exception = 5201

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
