package args

import "strings"

func ConvertArgsSliceToMap(args []string) map[string]string {
	result := make(map[string]string, len(args))
	for _, v := range args {
		a := strings.Split(v, "=")
		key := a[0]
		value := a[1]
		result[key] = value
	}
	return result
}
