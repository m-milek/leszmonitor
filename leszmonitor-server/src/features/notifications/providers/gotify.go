package providers

import "fmt"

type GotifyConfig struct {
	ServerURL string `json:"serverUrl" yaml:"serverUrl"`
	AppToken  string `json:"appToken,omitempty" yaml:"appToken,omitempty"`
	Priority  int    `json:"priority" yaml:"priority"`
}

func (c *GotifyConfig) Validate() error {
	if err := validateHTTPURL("serverUrl", c.ServerURL); err != nil {
		return err
	}
	if c.AppToken == "" {
		return fmt.Errorf("appToken cannot be empty")
	}
	if c.Priority < 0 || c.Priority > 10 {
		return fmt.Errorf("priority must be between 0 and 10")
	}
	return nil
}

func (c *GotifyConfig) Redact() {
	c.AppToken = ""
}

func (c *GotifyConfig) KeepSecretsFrom(previous Config) {
	prev, ok := previous.(*GotifyConfig)
	if !ok {
		return
	}
	if c.AppToken == "" {
		c.AppToken = prev.AppToken
	}
}

func (c *GotifyConfig) PlaceholderSecrets(envPrefix string) {
	c.AppToken = envPlaceholder(envPrefix, "app_token")
}
