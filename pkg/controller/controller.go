package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var Debug bool

const (
	clusterNamePrefix = "cluster_%v"
)

type Controller struct {
	BackendAddress string
	Ports          []int
	listenAddress  string
}

type envoyRequest struct {
	Node           map[string]string `json:"node"`
	RessourceNames []string          `json:"resource_names"`
}

func New(backendAddress string, ports []int, listenAddress string) *Controller {
	return &Controller{BackendAddress: backendAddress, Ports: ports, listenAddress: listenAddress}
}

func (c *Controller) MustRun() {
	http.HandleFunc("/v2/discovery:endpoints", handlerMiddleware(c.getEDSAnswer))
	http.HandleFunc("/v2/discovery:listeners", handlerMiddleware(c.getLDSAnswer))
	http.HandleFunc("/v2/discovery:clusters", handlerMiddleware(c.getCDSAnswer))
	http.HandleFunc("/", defaultHandler)
	go func() {
		if err := http.ListenAndServe(c.listenAddress, nil); err != nil {
			log.Fatalf("Controller failed to listen on %s: %v", c.listenAddress, err)
		}
	}()
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got an unexpected request on %s\n", r.URL.String())
	w.WriteHeader(http.StatusBadRequest)
	return
}

func handlerMiddleware(handleFunc func(*http.Request) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		answer, err := handleFunc(r)
		if err != nil {
			log.Printf("Error executing handle func: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(answer)
		if Debug {
			log.Printf("Response:\n---\n%s\n---\n", string(answer))
			log.Printf("Query: %s\n", r.URL.RawQuery)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return
			}
			log.Printf("Body: %s", string(body))
		}
		return
	}
}

func (c *Controller) getLDSAnswer(_ *http.Request) ([]byte, error) {
	// Bytes.WriteString apparently never returns an error, but panics
	// if the strings is too large
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			log.Printf("Recovered panic in getLDSAnswer: %v\n", r)
		}
	}()

	answer := bytes.Buffer{}
	answer.WriteString(`{"version_info": "0", "resources": [`)

	maxIndex := len(c.Ports) - 1
	for idx, port := range c.Ports {
		answer.WriteString(getLDSAnswer(fmt.Sprintf("frontend_%v", port), fmt.Sprintf(clusterNamePrefix, port), port))
		if idx < maxIndex {
			answer.WriteString(",")
		}
	}
	answer.WriteString("]}")
	return answer.Bytes(), nil
}

func (c *Controller) getCDSAnswer(_ *http.Request) ([]byte, error) {
	// Bytes.WriteString apparently never returns an error, but panics
	// if the strings is too large
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			log.Printf("Recovered panic in getCDSAnswer: %v\n", r)
		}
	}()

	answer := bytes.Buffer{}
	answer.WriteString(`{"version_info": "0", "resources": [`)
	maxIndex := len(c.Ports) - 1
	for idx, port := range c.Ports {
		answer.WriteString(getCDSAnswer(fmt.Sprintf(clusterNamePrefix, port)))
		if idx < maxIndex {
			answer.WriteString(",")
		}
	}
	answer.WriteString("]}")
	return answer.Bytes(), nil

}

func (c *Controller) getEDSAnswer(request *http.Request) ([]byte, error) {
	// Bytes.WriteString apparently never returns an error, but panics
	// if the strings is too large
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			log.Printf("Recovered panic in getDDSAnswer: %v\n", r)
		}
	}()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}
	var decodedRequest envoyRequest
	if err := json.Unmarshal(body, &decodedRequest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %v", err)
	}
	if len(decodedRequest.RessourceNames) != 1 {
		return nil, fmt.Errorf("expected to get exactly one ressourceNames, got %v", len(decodedRequest.RessourceNames))
	}

	var clusterPort int
	for _, port := range c.Ports {
		if fmt.Sprintf(clusterNamePrefix, port) == decodedRequest.RessourceNames[0] {
			clusterPort = port
		}
	}
	if clusterPort == 0 {
		return nil, fmt.Errorf("received request for unknown cluster %s", decodedRequest.RessourceNames[0])
	}

	answer := bytes.Buffer{}
	answer.WriteString(fmt.Sprintf(`{"version_info": "0", "resources": [
    {
      "@type": "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
      "cluster_name": "%s",
      "endpoints": [`, decodedRequest.RessourceNames[0]))
	answer.WriteString(getEDSAnswer(c.BackendAddress, clusterPort))
	answer.WriteString("]}]}")
	return answer.Bytes(), nil
}
