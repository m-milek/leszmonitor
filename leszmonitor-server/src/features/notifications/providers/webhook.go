package providers

import (
	"fmt"
	"strings"
)

type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

func (c *WebhookConfig) Validate() error {
	if err := validateHTTPURL("url", c.URL); err != nil {
		return err
	}
	for name, value := range c.Headers {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("header name cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("header %s cannot have an empty value", name)
		}
	}
	return nil
}

func (c *WebhookConfig) Redact() {
	for name := range c.Headers {
		c.Headers[name] = ""
	}
}

func (c *WebhookConfig) KeepSecretsFrom(previous Config) {
	prev, ok := previous.(*WebhookConfig)
	if !ok {
		return
	}
	for name, value := range c.Headers {
		if value == "" {
			c.Headers[name] = prev.Headers[name]
		}
	}
}
