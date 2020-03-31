package utility

import (
	"strings"

	valid "github.com/asaskevich/govalidator"
)

func ConvertParamsToQuery(params map[string]string) string {
	query := "@" + params["table"] + " "
	delete(params, "table")
	queryParams := make([]string, 0)

	for k, v := range params {
		eachParam := ""
		if valid.HasWhitespace(v) {
			count := 0
			words := strings.Split(v, " ")
			for _, eachWord := range words {
				if valid.IsAlpha(eachWord) {
					count++
				}
			}
			if count == len(words) {
				eachParam = k + "=" + "\"" + v + "\""
			} else {
				eachParam = k + "=" + v
			}
		} else {
			if valid.IsAlpha(v) {
				eachParam = k + "=" + "\"" + v + "\""
			} else {
				eachParam = k + "=" + v
			}

		}

		queryParams = append(queryParams, eachParam)
	}
	return (query + strings.Join(queryParams, " AND "))
}
