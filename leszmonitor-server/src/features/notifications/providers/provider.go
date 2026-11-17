package providers

import (
	"fmt"
	"net/url"
)

type Config interface {
	Validate() error
	Redact()
	KeepSecretsFrom(previous Config)
}

func validateHTTPURL(field, raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("%s must be a valid http or https URL", field)
	}
	return nil
}
