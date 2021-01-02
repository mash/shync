package shync

import (
	"fmt"

	"github.com/mash/shync/log"
	"github.com/mash/shync/shopify"
)

type Config struct {
	Store, Username, Password string
	In, Out                   string
	AllTemplates              bool
	Templates                 []string
}

func (c Config) Check() error {
	if c.Store == "" {
		return fmt.Errorf("config: store is required")
	}
	if c.Username == "" {
		return fmt.Errorf("config: username is required")
	}
	if c.Password == "" {
		return fmt.Errorf("config: password is required")
	}
	if c.AllTemplates && len(c.Templates) != 0 {
		log.Info("msg", "--all is set but --id is also set. --id wins")
	}
	if !c.AllTemplates && len(c.Templates) == 0 {
		return fmt.Errorf("config: at least one email template id is required")
	}
	if c.Out == "-" {
		if c.AllTemplates || len(c.Templates) > 1 {
			return fmt.Errorf("config: writing multiple email templates to stdout is not supported (why would you do that?)")
		}
		log.Debug("msg", "writing email template body to stdout")
	}
	for _, id := range c.Templates {
		if !shopify.IsValid(id) {
			return fmt.Errorf("config: invalid email template id: %s", id)
		}
	}
	return nil
}

func (c Config) Ids() []string {
	if len(c.Templates) > 0 {
		return c.Templates
	}
	return shopify.Templates
}
