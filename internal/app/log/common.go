package log

import "strings"

func (l *commonLogger) formatLogStack(stackValue string) string {
	stackSlice := strings.Split(stackValue, "\n")

	fPart := stackSlice[0]
	lPart := stackSlice[5:]

	result := []string{fPart}
	result = append(result, lPart...)

	return strings.Join(result, "\n")
}
