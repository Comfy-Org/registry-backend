package mapper

// BoolPtrToBool function to convert a bool pointer to a bool, returns false if the pointer is nil
func BoolPtrToBool(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}
