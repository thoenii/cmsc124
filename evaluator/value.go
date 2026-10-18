package evaluator

import "strconv"

// Value is any runtime value: float64, string, bool, or nil (none).
type Value = any

// Stringify is the single place that decides how a value prints.
func Stringify(v Value) string {
	switch x := v.(type) {
	case nil:
		return "dead"
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64) // 5.0 -> "5", 2.5 -> "2.5"
	case string:
		return x
	case bool:
		return strconv.FormatBool(x)
	}
	return "<unknown>"
}
