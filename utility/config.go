package utility

func BuildFormPostConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/" + "form"
	krakendConfig["method"] = "POST"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"
	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = "/form"
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildFormsGetConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/" + "form"
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildFormGetConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/" + "form" + "/" + "{table}"
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = "/form/{table}"
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildDataGetConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/" + "query"
	krakendConfig["method"] = "POST"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"
	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildSearchConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)
	krakendConfig["endpoint"] = "/api/table/{table}"
	krakendConfig["querystring_params"] = []string{"*"}
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildDataPostConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/" + "table" + "/" + config["table"].(string)
	krakendConfig["method"] = "POST"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"
	extraConfig["schema"] = config["schema"]
	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildUpdateConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/table/{table}/{_id}"
	krakendConfig["method"] = "PATCH"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildDeleteConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/table/{table}/{_id}"
	krakendConfig["method"] = "DELETE"
	krakendConfig["output_encoding"] = "no-op"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildDataFetchConfig(extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/table/{table}/{_id}"
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint

}
