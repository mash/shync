package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mash/shync"
	"github.com/mash/shync/cmd/mw"
	"github.com/mash/shync/log"
	"github.com/mash/shync/shopify"
	"gopkg.in/alecthomas/kingpin.v2"
)

// injected via go build -ldflags
var (
	Version   string
	BuildDate string
)

var (
	app         = kingpin.New("shync", "Shopify email template syncer")
	logLevel    = app.Flag("loglevel", "Log level (debug, info, error)").Default("info").Envar("SHYNC_LOGLEVEL").Enum("debug", "info", "error")
	debugChrome = app.Flag("debug-chrome", "Use non headless Chrome").Bool()
	store       = app.Flag("store", "The Shopify store URL. eg: `https://{shopname}.myshopify.com`").Envar("SHYNC_STORE").String()
	username    = app.Flag("username", "The Shopify admin username").Envar("SHYNC_USERNAME").String()
	password    = app.Flag("password", "The Shopify admin password").Envar("SHYNC_PASSWORD").String()

	version = app.Command("version", "Show version")

	ids = app.Command("ids", "Show email template ids")

	checkout          = app.Command("checkout", "Checkout email templates into a directory")
	checkoutDir       = checkout.Arg("dir", "Output directory, or - for stdout").Default(".").Envar("SHYNC_DIR").String()
	checkoutAll       = checkout.Flag("all", "Checkout all email templates").Short('a').Bool()
	checkoutTemplates = checkout.Flag("id", "Email template identifier to checkout").Short('i').Strings()

	push          = app.Command("push", "Push email templates to Shopify")
	pushDir       = push.Arg("dir", "Input directory where the email templates exist").Default(".").Envar("SHYNC_DIR").String()
	pushAll       = push.Flag("all", "Push all email templates").Short('a').Bool()
	pushTemplates = push.Flag("id", "Email template identifier to push").Short('i').Strings()
)

func main() {
	log.Info("version", Version, "buildDate", BuildDate)
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case version.FullCommand():
		break
	case ids.FullCommand():
		for _, v := range shopify.Templates {
			fmt.Println(v)
		}
	case checkout.FullCommand():
		c := shync.Config{
			Store:        *store,
			Username:     *username,
			Password:     *password,
			Dir:          *checkoutDir,
			AllTemplates: *checkoutAll,
			Templates:    *checkoutTemplates,
			Head:         *debugChrome,
		}
		if err := c.Check(); err != nil {
			log.Errorf("checkout: %s", err)
			return
		}
		fn := mw.Recover(mw.StatusLog(shync.Checkout))
		fn(c)
	case push.FullCommand():
		c := shync.Config{
			Store:        *store,
			Username:     *username,
			Password:     *password,
			Dir:          *pushDir,
			AllTemplates: *pushAll,
			Templates:    *pushTemplates,
			Head:         *debugChrome,
		}
		if err := c.Check(); err != nil {
			log.Errorf("push: %s", err)
			return
		}
		// fail fast
		if err := c.CheckReadable(); err != nil {
			log.Errorf("push: %s", err)
			return
		}
		fn := mw.Recover(mw.StatusLog(shync.Push))
		fn(c)
	}
}
