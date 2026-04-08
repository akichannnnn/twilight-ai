package utils

// ToFloat64 coerces an any value to float64, handling common numeric types.
// Returns (0, false) when the value is nil or an unsupported type.
func ToFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

// ToInt coerces an any value to int, handling common numeric types.
// Floating-point values are truncated toward zero.
// Returns (0, false) when the value is nil or an unsupported type.
func ToInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	case float32:
		return int(n), true
	default:
		return 0, false
	}
}
