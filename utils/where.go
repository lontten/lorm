package utils

func GenwhereToken(columns []string) []string {
	tokens := make([]string, 0)
	for _, column := range columns {
		tokens = append(tokens, column+" = ? ")
	}
	return tokens
}
