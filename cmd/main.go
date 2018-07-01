package main

import (
	"flag"
	"log"

	"github.com/alvaroaleman/envoy-configurator/pkg/api"
	"github.com/alvaroaleman/envoy-configurator/pkg/controller"
)

func main() {
	caFile := flag.String("cafile", "", "Path to the cacert file clients will be verified against")
	crtFile := flag.String("crtfile", "", "Path to the serving cert")
	keyFile := flag.String("keyfile", "", "Path to the serving certs key")
	apiListenAddress := flag.String("api-listen-address", "0.0.0.0:8090", "Address the api listens on")
	flag.Parse()
	if *caFile == "" || *crtFile == "" || *keyFile == "" || *apiListenAddress == "" {
		log.Fatalln("All of cafile, crtfile, keyfile and api-listen-address must be a non-empty string!")
	}

	stopChannel := make(chan struct{})
	controller := controller.New("192.168.0.39", []int{80, 443}, "127.0.0.1:8080")

	api, err := api.New(*caFile, *crtFile, *keyFile, *apiListenAddress, controller)
	if err != nil {
		log.Fatalf("Failed to create api: %v", err)
	}
	api.MustRun()
	controller.MustRun()

	<-stopChannel
}
