package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/v2/discovery:endpoints", edsHandler)
	http.HandleFunc("/v2/discovery:listeners", ldsHandler)
	http.HandleFunc("/", debug)
	http.ListenAndServe(":8000", nil)
}

const (
	hardCodedEdsAnswer = `
{
  "resources": [
    {
      "@type": "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
      "cluster_name": "backend_0",
      "endpoints": [
        {
          "lb_endpoints": [
            {
              "endpoint": {
                "address": {
                  "socket_address": {
                    "address": "192.168.0.39",
                    "port_value": 80
                  }
                }
              }
            }
          ]
        }
      ]
    }
  ],
  "version_info": "0"
}
`
	hardCodedLdsAnswer = `
{
  "resources": [
    {
      "@type": "type.googleapis.com/envoy.api.v2.Listener",
      "address": {
        "socket_address": {
          "address": "0.0.0.0",
          "port_value": 10000
        }
      },
      "filter_chains": [
        {
          "filters": [
            {
              "config": {
                "cluster": "backend_0",
                "idle_timeout": "10s",
                "stat_prefix": "frontend_0"
              },
              "name": "envoy.tcp_proxy"
            }
          ]
        }
      ],
      "name": "listener_0"
    }
  ],
  "version_info": "0"
}

				`
)

func edsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(hardCodedEdsAnswer))
	if err != nil {
		fmt.Printf("Error writing response: %v", err)
	}
	return
}

func ldsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(hardCodedLdsAnswer))
	if err != nil {
		fmt.Printf("Error writing response: %v", err)
	}
	return
}

func debug(w http.ResponseWriter, r *http.Request) {
	var body []byte
	_, err := r.Body.Read(body)
	r.Body.Close()
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("Error reading body: %v", err)
		return
	}
	fmt.Printf("url: %s\n", r.URL.String())
	fmt.Printf("url: %s\n", r.URL.RawPath)
	fmt.Printf("url: %s\n", r.URL.RawQuery)
	fmt.Printf("url: %s\n", r.URL.Opaque)
	fmt.Printf("url: %s\n", r.URL.RequestURI())
	fmt.Printf("Body: %s\n", string(body))
	return
}
