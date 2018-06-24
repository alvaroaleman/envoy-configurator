package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", dump)
	http.ListenAndServe(":8000", nil)
}

const hardCodedEdsAnswer = `
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

func dump(w http.ResponseWriter, r *http.Request) {
	var body []byte
	_, err := r.Body.Read(body)
	r.Body.Close()
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("Error reading body: %v", err)
		return
	}
	_, err = w.Write([]byte(hardCodedEdsAnswer))
	if err != nil {
		fmt.Errorf("Error writting response: %v", err)
	}
	return
}
