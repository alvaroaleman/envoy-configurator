package main

import (
	"fmt"
	"net/http"

	"github.com/alvaroaleman/envoy-configurator/pkg/controller"
)

func main() {
	stopChannel := make(chan struct{})
	controller := controller.New("192.168.0.39", []int{80, 443}, "127.0.0.1:8080")
	controller.MustRun()

	<-stopChannel
}
