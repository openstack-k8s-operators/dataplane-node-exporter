// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Miguel Lavalle

package ovn

import (
	"bufio"
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/openstack-k8s-operators/dataplane-node-exporter/appctl"
	"github.com/openstack-k8s-operators/dataplane-node-exporter/collectors/lib"
	"github.com/openstack-k8s-operators/dataplane-node-exporter/config"
	"github.com/openstack-k8s-operators/dataplane-node-exporter/log"
	"github.com/openstack-k8s-operators/dataplane-node-exporter/ovsdb"
	"github.com/openstack-k8s-operators/dataplane-node-exporter/ovsdb/ovs"
	"github.com/prometheus/client_golang/prometheus"
)

func collectopenvSwitch(externaIds map[string]string, ch chan<- prometheus.Metric) {
	for name, metric := range openvSwitch {
		value, ok := externaIds[name]
		if !ok {
			continue
		}
		if !config.MetricSets().Has(metric.Set) {
			continue
		}

		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Errf("%s: %s: %s", name, value, err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(metric.Desc(), metric.ValueType, val)
	}
}

func collectopenvSwitchBoolean(externaIds map[string]string, ch chan<- prometheus.Metric) {
	for name, metric := range openvSwitchBoolean {
		value, ok := externaIds[name]
		if !ok {
			continue
		}
		if !config.MetricSets().Has(metric.Set) {
			continue
		}

		val := 1.0
		if value != "True" {
			val = 0
		}

		ch <- prometheus.MustNewConstMetric(metric.Desc(), metric.ValueType, val)
	}
}

func collectopenvSwitchLabels(externaIds map[string]string, ch chan<- prometheus.Metric) {
	for name, metric := range openvSwitchLabels {
		label, ok := externaIds[name]
		if !ok {
			continue
		}
		if !config.MetricSets().Has(metric.Set) {
			continue
		}

		value := 1.0
		ch <- prometheus.MustNewConstMetric(metric.Desc(), metric.ValueType, value, []string{label}...)
	}
}

func makeMetric(name, value string) prometheus.Metric {
	m, ok := ovnController[name]
	if !ok {
		return nil
	}
	if !config.MetricSets().Has(m.Set) {
		return nil
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Errf("%s: %s: %s", name, value, err)
		return nil
	}

	return prometheus.MustNewConstMetric(m.Desc(), m.ValueType, val)
}

// "vconn_sent                 0.0/sec     0.083/sec        0.0767/sec   total: 131870"
var coverageRe = regexp.MustCompile(`^(\w+)\s+.*\s+total: (\d+)$`)

func collectCoverageMetrics(ch chan<- prometheus.Metric) {
	buf := appctl.OvnController("coverage/show")
	if buf == "" {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(buf))
	for scanner.Scan() {
		line := scanner.Text()

		match := coverageRe.FindStringSubmatch(line)
		if match != nil {
			metric := makeMetric(match[1], match[2])
			if metric != nil {
				ch <- metric
			}
		}
	}
}

type Collector struct{}

func (Collector) Name() string {
	return "ovn"
}

func (Collector) Metrics() []lib.Metric {
	var res []lib.Metric
	for _, g := range metrics {
		for _, m := range *g {
			res = append(res, m)
		}
	}
	return res
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	lib.DescribeEnabledMetrics(c, ch)
}

func (Collector) Collect(ch chan<- prometheus.Metric) {
	// collect items from the ExternalIDs field in the OpenvSwitch table
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var vswitch ovs.OpenvSwitch
	err := ovsdb.Get(ctx, &vswitch)
	if err != nil {
		log.Errf("OvsdbGet(vswitch): %s", err)
		return
	}
	collectopenvSwitch(vswitch.ExternalIDs, ch)
	collectopenvSwitchBoolean(vswitch.ExternalIDs, ch)
	collectopenvSwitchLabels(vswitch.ExternalIDs, ch)

	// collect the ovn-controller coverage metrics
	collectCoverageMetrics(ch)
}
