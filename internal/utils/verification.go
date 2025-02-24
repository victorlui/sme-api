package utils

// Função auxiliar para verificar zero-values
func IsZeroValue(val interface{}) bool {
	switch v := val.(type) {
	case int:
		return v == 0
	case string:
		return v == ""
	case nil:
		return true
	default:
		return false
	}
}
