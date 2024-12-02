// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Miguel Lavalle

package ovn

import (
	"github.com/openstack-k8s-operators/dataplane-node-exporter/collectors/lib"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct{}

func (Collector) Name() string {
	return "ovn"
}

func (Collector) Metrics() []lib.Metric {
	var res []lib.Metric
	for _, m := range metrics {
		res = append(res, m)
	}
	return res
}

func (Collector) Collect(ch chan<- prometheus.Metric) {
}
