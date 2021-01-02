package shync

import (
	"fmt"

	"github.com/mash/shync/shopify"
)

func Checkout(c Config) error {
	client := shopify.NewClient(true)
	// you might want to skip cancel() to keep the browser window around
	if false {
		defer client.Close()
	}

	// run task list
	if err := client.Login(c.Store, c.Username, c.Password); err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}

	var subject, body string
	if err := client.FetchEmailTemplate(shopify.Templates[0], &subject, &body); err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}
	return nil
}
