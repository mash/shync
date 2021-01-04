package shync

import (
	"fmt"
	"io/ioutil"

	"github.com/mash/shync/log"
	"github.com/mash/shync/shopify"
)

func Push(c Config) error {
	log.Infof("opening Chrome")
	client := shopify.NewClient(c.Head)
	// you might want to skip cancel() to keep the browser window around
	if false {
		defer client.Close()
	}

	// run task list
	log.Infof("logging into Shopify")
	if err := client.Login(c.Store, c.Username, c.Password); err != nil {
		return fmt.Errorf("Push: %w", err)
	}

	for _, id := range c.Ids() {
		log.Infof("working on %s", id)
		if err := pushSingle(c, client, id); err != nil {
			return fmt.Errorf("Push: failed processing: %s, error: %w", id, err)
		}
	}
	return nil
}

func pushSingle(c Config, client *shopify.Client, id string) error {
	var subject, body string
	if err := read(c, c.SubjectFile(id), &subject); err != nil {
		if c.IgnoreMissing {
			log.Infof("%s not found, ignoring", c.SubjectFile(id))
		} else {
			return err
		}
	}
	if err := read(c, c.BodyFile(id), &body); err != nil {
		if c.IgnoreMissing {
			log.Infof("%s not found, ignoring", c.BodyFile(id))
		} else {
			return err
		}
	}
	if subject == "" && body == "" {
		return nil
	}
	if err := client.UpdateEmailTemplate(shopify.Templates[0], subject, body); err != nil {
		return err
	}
	return nil
}

func read(c Config, file string, content *string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	*content = string(b)
	return nil
}
