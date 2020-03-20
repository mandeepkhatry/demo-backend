package utility

func BuildFormGetConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/" + "form" + "/" + "{" + config["table"].(string) + "}"
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"
	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = "/form"
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

	krakendConfig["endpoint"] = "/api/" + "table" + "/" + "{name}"
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

func BuildDataGetConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)
	schema := make(map[string]interface{})

	schema["required"] = []string{"query"}
	schema["title"] = "query"
	schema["type"] = "object"

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

func BuildUpdateConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/" + config["table"].(string) + "/{id}"
	krakendConfig["method"] = "POST"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = "/data/{id}"
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildSearchConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	
	krakendConfig["endpoint"] = "/api/table/search"
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"

	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = krakendConfig["endpoint"]
	backendConfig["host" = []string{extra["http_target"]}

	krakenbackendConfig["backend"] = []map[string]interface{}{backbackendConfig}

	endpoint = append(endpoint, krakendConkrakendConfig)

	return endpoint
}