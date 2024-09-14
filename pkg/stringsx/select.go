package stringsx

func SelectString(predicate bool, trueValue string, falseValue string) string {
	if predicate {
		return trueValue
	}
	return falseValue
}
