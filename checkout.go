package shync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mash/shync/log"
	"github.com/mash/shync/shopify"
)

func Checkout(c Config) error {
	log.Infof("opening Chrome")
	client := shopify.NewClient(c.Head)
	// you might want to skip cancel() to keep the browser window around
	if false {
		defer client.Close()
	}

	// run task list
	log.Infof("logging into Shopify")
	if err := client.Login(c.Store, c.Username, c.Password); err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}

	for _, id := range c.Ids() {
		log.Infof("checking out %s", id)
		if err := checkoutSingle(c, client, id); err != nil {
			return fmt.Errorf("Checkout: failed processing: %s, error: %w", id, err)
		}
	}
	return nil
}

func checkoutSingle(c Config, client *shopify.Client, id string) error {
	var subject, body string
	if err := client.FetchEmailTemplate(id, &subject, &body); err != nil {
		return err
	}
	return write(c, id, subject, body)
}

func write(c Config, id, subject, body string) error {
	if c.Out == "-" {
		fmt.Println(body)
		return nil
	}
	subjectFile := filepath.Join(c.Out, fmt.Sprintf("%s.txt", id))
	bodyFile := filepath.Join(c.Out, fmt.Sprintf("%s.html", id))

	sf, err := os.Create(subjectFile)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(sf, subject); err != nil {
		return err
	}

	bf, err := os.Create(bodyFile)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(bf, body); err != nil {
		return err
	}
	return nil
}
