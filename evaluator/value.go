package evaluator

import ("strconv"
		"fmt"
)

// Value is any runtime value: float64, string, bool, or nil (none).
type Value = any

// isTruthy is the single truthiness rule: only false and nil are falsey. 0 and "" are truthy
func isTruthy(v Value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

// isEqual compares supported runtime values.
func isEqual(a, b Value) bool {
    if a == nil || b == nil {
        return a == nil && b == nil
    }

    switch left := a.(type) {
    case float64:
        right, ok := b.(float64)
        return ok && left == right

    case string:
        right, ok := b.(string)
        return ok && left == right

    case bool:
        right, ok := b.(bool)
        return ok && left == right

    default:
        return false
    }
}

// Stringify is the single place that decides how a value prints.
func Stringify(v Value) string {
	switch x := v.(type) {
	case nil:
		return "dead"
	case float64:
		if x == 0 {
			return "0" // catches -0 too
		}
		return strconv.FormatFloat(x, 'f', -1, 64) // 5.0 -> "5", 2.5 -> "2.5"
	case string:
		return x
	case bool:
		return strconv.FormatBool(x)
	}
	return fmt.Sprintf("%v", v) // bad value shows up in test output
}
