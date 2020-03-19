package utility

func BuildFormGetConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/" + config["name"].(string)
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

func BuildTablePostConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/" + config["name"].(string)
	krakendConfig["method"] = "POST"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"
	extraConfig["schema"] = config["schema"]
	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = "/data"
	backendConfig["host"] = []string{extra["http_target"]}

	krakendConfig["backend"] = []map[string]interface{}{backendConfig}

	endpoint = append(endpoint, krakendConfig)

	return endpoint
}

func BuildTableGetConfig(config map[string]interface{}, extra map[string]string) []interface{} {
	krakendConfig := make(map[string]interface{})
	extraConfig := make(map[string]interface{})
	backendConfig := make(map[string]interface{})
	endpoint := make([]interface{}, 0)

	krakendConfig["endpoint"] = "/api/" + config["name"].(string)
	krakendConfig["method"] = "GET"
	krakendConfig["output_encoding"] = "no-op"

	extraConfig["decoding"] = "json"
	krakendConfig["extra_config"] = extraConfig

	backendConfig["url_pattern"] = "/data"
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

	krakendConfig["endpoint"] = "/api/" + config["name"].(string) + "/{id}"
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
