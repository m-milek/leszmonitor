package providers

import (
	"fmt"
	"strings"
)

type NtfyConfig struct {
	ServerURL string `json:"serverUrl"`
	Topic     string `json:"topic"`
	Token     string `json:"token,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Priority  int    `json:"priority,omitempty"`
}

func (c *NtfyConfig) Validate() error {
	if err := validateHTTPURL("serverUrl", c.ServerURL); err != nil {
		return err
	}
	if strings.TrimSpace(c.Topic) == "" {
		return fmt.Errorf("topic cannot be empty")
	}
	if c.Token != "" && c.Username != "" {
		return fmt.Errorf("token and username/password are mutually exclusive")
	}
	if (c.Username == "") != (c.Password == "") {
		return fmt.Errorf("username and password must be set together")
	}
	if c.Priority < 0 || c.Priority > 5 {
		return fmt.Errorf("priority must be between 1 and 5")
	}
	return nil
}

func (c *NtfyConfig) Redact() {
	c.Token = ""
	c.Password = ""
}

func (c *NtfyConfig) KeepSecretsFrom(previous Config) {
	prev, ok := previous.(*NtfyConfig)
	if !ok {
		return
	}
	if c.Token == "" && c.Username == "" {
		c.Token = prev.Token
	}
	if c.Password == "" && c.Username != "" && c.Username == prev.Username {
		c.Password = prev.Password
	}
}
