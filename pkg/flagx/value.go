package flagx

type Value interface {
	Parse(...string) (int, error)
}

func FormatSlice(v Value) []string {
	if v, ok := v.(interface{ Format() []string }); ok {
		return v.Format()
	}
	if v, ok := v.(interface{ Format() string }); ok {
		if f := v.Format(); f != "" {
			return []string{f}
		}
	}
	return nil
}

func Format(v Value) string {
	if v, ok := v.(interface{ Format() []string }); ok {
		if ff := v.Format(); len(ff) == 1 {
			return ff[0]
		}
	}
	if v, ok := v.(interface{ Format() string }); ok {
		return v.Format()
	}
	return ""
}

func IsInlined(v Value) bool {
	if v, ok := v.(interface{ IsInlined() bool }); ok {
		return v.IsInlined()
	}
	return false
}

func SelectValue(predicate bool, trueValue Value, falseValue Value) Value {
	if predicate {
		return trueValue
	}
	return falseValue
}
