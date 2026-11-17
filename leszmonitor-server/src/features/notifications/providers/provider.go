package providers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Config interface {
	Validate() error
	Redact()
	KeepSecretsFrom(previous Config)
	PlaceholderSecrets(envPrefix string)
}

var nonEnvChars = regexp.MustCompile(`[^A-Z0-9]+`)

func envPlaceholder(parts ...string) string {
	name := nonEnvChars.ReplaceAllString(strings.ToUpper(strings.Join(parts, "_")), "_")
	return "${" + strings.Trim(name, "_") + "}"
}

func validateHTTPURL(field, raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("%s must be a valid http or https URL", field)
	}
	return nil
}
