package utility

import (
	"fmt"
	"strings"

	valid "github.com/asaskevich/govalidator"
)

func ConvertParamsToQuery(params map[string]string) string {
	query := "@" + params["table"] + " "
	delete(params, "table")
	queryParams := make([]string, 0)
	fmt.Println("params : ", params)
	for k, v := range params {
		eachParam := ""
		if valid.IsAlpha(v) {
			eachParam = k + "=" + "\"" + v + "\""
		} else {
			eachParam = k + "=" + v
		}
		queryParams = append(queryParams, eachParam)
	}
	return (query + strings.Join(queryParams, " AND "))
}
