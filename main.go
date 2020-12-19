package main

import (
	"carcompare/api"
	"carcompare/api/components/providers"
	"carcompare/api/components/providers/autotrader"
	"carcompare/api/components/providers/ebay"
	"carcompare/api/components/providers/facebook"
	"carcompare/api/database"
	"carcompare/cache"
	"log"
	"os"

	"github.com/apsdehal/go-logger"
)

func main() {
	database.Connect()
	cache.Connect()

	l, err := logger.New("", 1, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	l.SetLogLevel(logger.DebugLevel)

	manager, err := providers.NewManager(l)
	if err != nil {
		l.Fatalf("%s", err)
	}

	autotraderProvider := autotrader.NewProvider(l)
	if err := manager.RegisterProvider("autotrader", autotraderProvider); err != nil {
		l.Fatalf("%s", err)
	}

	ebayProvider := ebay.NewProvider(l)
	if err := manager.RegisterProvider("ebay", ebayProvider); err != nil {
		l.Fatalf("%s", err)
	}

	facebookProvider := facebook.NewProvider(l)
	if err := manager.RegisterProvider("facebook", facebookProvider); err != nil {
		l.Fatalf("%s", err)
	}

	api.Run(manager, l)

	defer cache.RedisDB.Close()
}
