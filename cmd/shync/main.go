package main

import (
	"os"

	"github.com/mash/shync"
	"github.com/mash/shync/cmd/mw"
	"github.com/mash/shync/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// injected via go build -ldflags
var (
	Version   string
	BuildDate string
	GoVersion string
)

var (
	app      = kingpin.New("shync", "Shopify email template syncer")
	logLevel = app.Flag("loglevel", "Log level (debug, info, error)").Default("info").Envar("SHYNC_LOGLEVEL").Enum("debug", "info", "error")
	store    = app.Flag("store", "The Shopify store URL. eg: `https://{shopname}.myshopify.com`").Required().Envar("SHYNC_STORE").String()
	username = app.Flag("username", "The Shopify admin username").Required().Envar("SHYNC_USERNAME").String()
	password = app.Flag("password", "The Shopify admin password").Required().Envar("SHYNC_PASSWORD").String()

	version = app.Command("version", "Show version")

	download    = app.Command("download", "Download email templates into a directory")
	downloadOut = download.Flag("out", "Output directory").Default(".").Envar("SHYNC_OUTDIR").Short('o').String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case version.FullCommand():
		log.Info("version", Version, "buildDate", BuildDate, "goVersion", GoVersion)
	case download.FullCommand():
		c := shync.Config{
			Store:    *store,
			Username: *username,
			Password: *password,
			Out:      *downloadOut,
		}
		fn := mw.Recover(mw.StatusLog(shync.Download))
		fn(c)
	}
}
