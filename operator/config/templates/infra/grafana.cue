// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package infra

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"

_generatedName: resource.metadata.name | *resource.status.clusterName

_dashboardConfigMapData: {
	"annotations": {
		"list": [
		{
			"builtIn": 1,
			"datasource": "-- Grafana --",
			"enable": true,
			"hide": true,
			"iconColor": "rgba(0, 211, 255, 1)",
			"name": "Annotations & Alerts",
			"type": "dashboard"
		}
		]
	},
	"description": "Monitors Kubernetes cluster using Prometheus. Shows overall cluster CPU / Memory / Filesystem Uses cAdvisor metrics only.",
	"editable": true,
	"gnetId": 315,
	"graphTooltip": 0,
	"links": [],
	"panels": [
		{
		  "collapsed": false,
		  "datasource": null,
		  "gridPos": {
			"h": 1,
			"w": 24,
			"x": 0,
			"y": 0
		  },
		  "id": 33,
		  "panels": [],
		  "title": "Network I/O pressure",
		  "type": "row"
		},
		{
		  "aliasColors": {},
		  "bars": false,
		  "dashLength": 10,
		  "dashes": false,
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "fill": 1,
		  "fillGradient": 0,
		  "grid": {},
		  "gridPos": {
			"h": 6,
			"w": 24,
			"x": 0,
			"y": 1
		  },
		  "height": "200px",
		  "hiddenSeries": false,
		  "id": 32,
		  "isNew": true,
		  "legend": {
			"alignAsTable": false,
			"avg": true,
			"current": true,
			"max": false,
			"min": false,
			"rightSide": false,
			"show": false,
			"sideWidth": 200,
			"sort": "current",
			"sortDesc": true,
			"total": false,
			"values": true
		  },
		  "lines": true,
		  "linewidth": 2,
		  "links": [],
		  "nullPointMode": "connected",
		  "percentage": false,
		  "pluginVersion": "7.1.1",
		  "pointradius": 5,
		  "points": false,
		  "renderer": "flot",
		  "seriesOverrides": [],
		  "spaceLength": 10,
		  "stack": false,
		  "steppedLine": false,
		  "targets": [
			{
			  "expr": "sum (rate (container_network_receive_bytes_total{kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=~\"\(_generatedName)\"}[1m]))",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "legendFormat": "Received",
			  "metric": "network",
			  "refId": "A",
			  "step": 10
			},
			{
			  "expr": "- sum (rate (container_network_transmit_bytes_total{kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}[1m]))",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "legendFormat": "Sent",
			  "metric": "network",
			  "refId": "B",
			  "step": 10
			}
		  ],
		  "thresholds": [],
		  "timeFrom": null,
		  "timeRegions": [],
		  "timeShift": null,
		  "title": "Network I/O pressure",
		  "tooltip": {
			"msResolution": false,
			"shared": true,
			"sort": 0,
			"value_type": "cumulative"
		  },
		  "type": "graph",
		  "xaxis": {
			"buckets": null,
			"mode": "time",
			"name": null,
			"show": true,
			"values": []
		  },
		  "yaxes": [
			{
			  "format": "Bps",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": true
			},
			{
			  "format": "Bps",
			  "label": null,
			  "logBase": 1,
			  "max": null,
			  "min": null,
			  "show": false
			}
		  ],
		  "yaxis": {
			"align": false,
			"alignLevel": null
		  }
		},
		{
		  "collapsed": false,
		  "datasource": null,
		  "gridPos": {
			"h": 1,
			"w": 24,
			"x": 0,
			"y": 7
		  },
		  "id": 34,
		  "panels": [],
		  "title": "Total usage",
		  "type": "row"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": true,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "percent",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": true,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 5,
			"w": 8,
			"x": 0,
			"y": 8
		  },
		  "height": "180px",
		  "id": 4,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "50%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (container_memory_working_set_bytes{id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}) / sum (machine_memory_bytes{kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}) * 100",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "65, 90",
		  "title": "Cluster memory usage",
		  "type": "singlestat",
		  "valueFontSize": "80%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": true,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "percent",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": true,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 5,
			"w": 8,
			"x": 8,
			"y": 8
		  },
		  "height": "180px",
		  "id": 6,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "50%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (rate (container_cpu_usage_seconds_total{id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}[1m])) / sum (machine_cpu_cores{kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}) * 100",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "65, 90",
		  "title": "Cluster CPU usage (1m avg)",
		  "type": "singlestat",
		  "valueFontSize": "80%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": true,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "percent",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": true,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 5,
			"w": 8,
			"x": 16,
			"y": 8
		  },
		  "height": "180px",
		  "id": 7,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "50%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (container_fs_usage_bytes{device=~\"^/dev/[sv]d[a-z][1-9]$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}) / sum (container_fs_limit_bytes{device=~\"^/dev/[sv]d[a-z][1-9]$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}) * 100",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "legendFormat": "",
			  "metric": "",
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "65, 90",
		  "title": "Cluster filesystem usage",
		  "type": "singlestat",
		  "valueFontSize": "80%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": false,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "bytes",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": false,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 3,
			"w": 4,
			"x": 0,
			"y": 13
		  },
		  "height": "1px",
		  "id": 9,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "20%",
		  "prefix": "",
		  "prefixFontSize": "20%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (container_memory_working_set_bytes{id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"})",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "",
		  "title": "Used",
		  "type": "singlestat",
		  "valueFontSize": "50%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": false,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "bytes",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": false,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 3,
			"w": 4,
			"x": 4,
			"y": 13
		  },
		  "height": "1px",
		  "id": 10,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "50%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (machine_memory_bytes{kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"})",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "",
		  "title": "Total",
		  "type": "singlestat",
		  "valueFontSize": "50%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": false,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "none",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": false,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 3,
			"w": 4,
			"x": 8,
			"y": 13
		  },
		  "height": "1px",
		  "id": 11,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": " cores",
		  "postfixFontSize": "30%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (rate (container_cpu_usage_seconds_total{id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"}[1m]))",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "",
		  "title": "Used",
		  "type": "singlestat",
		  "valueFontSize": "50%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": false,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "none",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": false,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 3,
			"w": 4,
			"x": 12,
			"y": 13
		  },
		  "height": "1px",
		  "id": 12,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": " cores",
		  "postfixFontSize": "30%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (machine_cpu_cores{kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"})",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "",
		  "title": "Total",
		  "type": "singlestat",
		  "valueFontSize": "50%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": false,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "bytes",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": false,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 3,
			"w": 4,
			"x": 16,
			"y": 13
		  },
		  "height": "1px",
		  "id": 13,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "50%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (container_fs_usage_bytes{device=~\"^/dev/[sv]d[a-z][1-9]$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"})",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "",
		  "title": "Used",
		  "type": "singlestat",
		  "valueFontSize": "50%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		},
		{
		  "cacheTimeout": null,
		  "colorBackground": false,
		  "colorValue": false,
		  "colors": [
			"rgba(50, 172, 45, 0.97)",
			"rgba(237, 129, 40, 0.89)",
			"rgba(245, 54, 54, 0.9)"
		  ],
		  "datasource": "Prometheus",
		  "decimals": 2,
		  "editable": true,
		  "error": false,
		  "fieldConfig": {
			"defaults": {
			  "custom": {}
			},
			"overrides": []
		  },
		  "format": "bytes",
		  "gauge": {
			"maxValue": 100,
			"minValue": 0,
			"show": false,
			"thresholdLabels": false,
			"thresholdMarkers": true
		  },
		  "gridPos": {
			"h": 3,
			"w": 4,
			"x": 20,
			"y": 13
		  },
		  "height": "1px",
		  "id": 14,
		  "interval": null,
		  "isNew": true,
		  "links": [],
		  "mappingType": 1,
		  "mappingTypes": [
			{
			  "name": "value to text",
			  "value": 1
			},
			{
			  "name": "range to text",
			  "value": 2
			}
		  ],
		  "maxDataPoints": 100,
		  "nullPointMode": "connected",
		  "nullText": null,
		  "postfix": "",
		  "postfixFontSize": "50%",
		  "prefix": "",
		  "prefixFontSize": "50%",
		  "rangeMaps": [
			{
			  "from": "null",
			  "text": "N/A",
			  "to": "null"
			}
		  ],
		  "sparkline": {
			"fillColor": "rgba(31, 118, 189, 0.18)",
			"full": false,
			"lineColor": "rgb(31, 120, 193)",
			"show": false
		  },
		  "tableColumn": "",
		  "targets": [
			{
			  "expr": "sum (container_fs_limit_bytes{device=~\"^/dev/[sv]d[a-z][1-9]$\",id=\"/\",kubernetes_io_hostname=~\"^$Node$\",test_cluster_name=\"\(_generatedName)\"})",
			  "interval": "10s",
			  "intervalFactor": 1,
			  "refId": "A",
			  "step": 10
			}
		  ],
		  "thresholds": "",
		  "title": "Total",
		  "type": "singlestat",
		  "valueFontSize": "50%",
		  "valueMaps": [
			{
			  "op": "=",
			  "text": "N/A",
			  "value": "null"
			}
		  ],
		  "valueName": "current"
		}
	  ],
	  "refresh": "10s",
	  "schemaVersion": 26,
	  "style": "dark",
	  "tags": [
		"kubernetes"
	  ],
	  "templating": {
		"list": [
		  {
			"allValue": ".*",
			"current": {
			  "selected": false,
			  "text": "All",
			  "value": "$__all"
			},
			"datasource": "Prometheus",
			"definition": "",
			"hide": 0,
			"includeAll": true,
			"label": null,
			"multi": false,
			"name": "Node",
			"options": [],
			"query": "label_values(kubernetes_io_hostname)",
			"refresh": 1,
			"regex": "",
			"skipUrlSync": false,
			"sort": 0,
			"tagValuesQuery": "",
			"tags": [],
			"tagsQuery": "",
			"type": "query",
			"useTags": false
		  }
		]
	  },
	  "time": {
		"from": "now-5m",
		"to": "now"
	  },
	  "timepicker": {
		"refresh_intervals": [
		  "10s",
		  "30s",
		  "1m",
		  "5m",
		  "15m",
		  "30m",
		  "1h",
		  "2h",
		  "1d"
		],
		"time_options": [
		  "5m",
		  "15m",
		  "1h",
		  "6h",
		  "12h",
		  "24h",
		  "2d",
		  "7d",
		  "30d"
		]
	  },
	  "timezone": "browser",
	  "title": "Test cluster monitoring for \(_generatedName)",
	  "version": 1
 }

resource: v1alpha2.#TestClusterGKE
