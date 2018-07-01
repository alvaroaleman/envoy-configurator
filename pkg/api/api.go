package api

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/alvaroaleman/envoy-configurator/pkg/controller"
)

type configurationRequest struct {
	Ports []int `json:"ports"`
}

type Handler struct {
	listenAddress   string
	tlsConfig       *tls.Config
	envoyController *controller.Controller
}

func New(cafile, certfile, keyfile, listenAddress string, envoyController *controller.Controller) (*Handler, error) {

	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return &Handler{listenAddress: listenAddress, tlsConfig: tlsConfig, envoyController: envoyController}, nil
}

func (c *Handler) MustRun() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/v1/update", c.updateRequestHandler)
	server := http.Server{
		Handler:   mux,
		Addr:      c.listenAddress,
		TLSConfig: c.tlsConfig,
	}
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("Failed to start apiserver: %v", err)
		}
	}()
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got an unexpected request on %s\n", r.URL.String())
	w.WriteHeader(http.StatusBadRequest)
	return
}

func (c *Handler) updateRequestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var request configurationRequest
	if err := json.Unmarshal(body, &request); err != nil {
		errMsg := fmt.Sprintf("Error unmarshalling body: %v", err)
		log.Printf(errMsg)
		if _, err := w.Write([]byte(errMsg)); err != nil {
			log.Printf("Error writing response body: %v", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	remoteAddrSlice := strings.Split(r.RemoteAddr, ":")
	remoteAddr := strings.Join(remoteAddrSlice[:len(remoteAddrSlice)-1], ":")

	if remoteAddr != c.envoyController.BackendAddress {
		log.Printf("Updating backend address from %s to %s", c.envoyController.BackendAddress, remoteAddr)
		c.envoyController.BackendAddress = remoteAddr
	}

	if request.Ports != nil &&
		c.envoyController.Ports != nil &&
		!reflect.DeepEqual(request.Ports, c.envoyController.Ports) {
		log.Printf("Updating ports from %v to %v", c.envoyController.Ports, request.Ports)
		c.envoyController.Ports = request.Ports
	}

	return
}
