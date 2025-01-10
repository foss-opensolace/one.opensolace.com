package utils

var (
	False *bool = ToPtr(false)
	True  *bool = ToPtr(true)
)

func ToPtr[T any](obj T) *T {
	return &obj
}
