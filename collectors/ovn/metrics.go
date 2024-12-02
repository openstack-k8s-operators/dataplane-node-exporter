package ovn

import (
	"github.com/openstack-k8s-operators/dataplane-node-exporter/collectors/lib"
	"github.com/openstack-k8s-operators/dataplane-node-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

var metrics = map[string]lib.Metric{
	"encap_type": {
		Name:        "ovnc_encap_type",
		Description: "encapsulation type that a chassis should use to connect to this node.",
		ValueType:   prometheus.GaugeValue,
		Set:         config.METRICS_BASE,
	},
	"bridge_mappings": {
		Name:        "ovnc_bridge_mappings",
		Description: "list  of  key-value  pairs that map a physical network name to a local ovs bridge that provides connectivity  to that  network.",
		ValueType:   prometheus.GaugeValue,
		Set:         config.METRICS_BASE,
	},
}
