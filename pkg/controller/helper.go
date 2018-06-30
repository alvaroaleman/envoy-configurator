package controller

import (
	"fmt"
)

func getLDSAnswer(name, backend string, port int) string {
	return fmt.Sprintf(`
    {
      "@type": "type.googleapis.com/envoy.api.v2.Listener",
      "address": {
        "socket_address": {
          "address": "0.0.0.0",
          "port_value": %v
        }
      },
      "filter_chains": [
        {
          "filters": [
            {
              "config": {
                "cluster": "%s",
                "idle_timeout": "10s",
                "stat_prefix": "%s"
              },
              "name": "envoy.tcp_proxy"
            }
          ]
        }
      ],
      "name": "%s"
    }
		`, port, backend, name, name)
}

func getCDSAnswer(clusterName string) string {
	return fmt.Sprintf(`
	    {
      "@type": "type.googleapis.com/envoy.api.v2.Cluster",
      "connect_timeout": "0.25s",
      "eds_cluster_config": {
        "eds_config": {
          "api_config_source": {
            "api_type": "REST",
            "cluster_names": "xds_cluster",
            "refresh_delay": "1s"
          }
        }
      },
      "lb_policy": "ROUND_ROBIN",
      "name": "%s",
      "type": "EDS"
    }
	`, clusterName)
}

func getEDSAnswer(address string, port int) string {
	return fmt.Sprintf(`
        {
          "lb_endpoints": [
            {
              "endpoint": {
                "address": {
                  "socket_address": {
                    "address": "%s",
                    "port_value": %v
                  }
                }
              }
            }
          ]
        }
`, address, port)
}
