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
	IdUnknown Exception = 1000

	IdGenericNotFound Exception = 3101

	IdOneOrManyValidation Exception = 4101 /* Example: Expected: n >= 1 | Error: 0 */
	IdInvalidFieldType    Exception = 4102 /* Example: Expected: 1 | Error: "1" */
	IdInvalidFieldLayout  Exception = 4103 /* Example: Expected: "2006-10-20" | Error: "2006-20-10" */

	IdAuthInvalidCredentials Exception = 5101
	IdAuthDuplicated         Exception = 5201

	IdMissingAPIKeyHeader     Exception = 6101
	IdInvalidAPIKey           Exception = 6102
	IdCannotUseAPIKey         Exception = 6103
	IdAPIKeyMissingPermission Exception = 6104
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
