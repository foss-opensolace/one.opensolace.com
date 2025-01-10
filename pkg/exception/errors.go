package exception

type FieldTypeError struct {
	Value string
	Field string
	Type  string
}

type FieldLayoutError struct {
	Value  string
	Layout string
}

func (fte FieldTypeError) Error() string {
	return "Cannot use " + fte.Value + " from '" + fte.Field + "' as a value of type " + fte.Type
}

func (fle FieldLayoutError) Error() string {
	return "Couldn't parse " + fle.Value + ". Expected layout: " + fle.Layout
}
