// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package infra

import "github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha2"

_generatedName: resource.metadata.name | *resource.status.clusterName

_dashboardConfigMapDataCiliumHubble:
    {
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
      "editable": true,
      "gnetId": null,
      "graphTooltip": 0,
      "id": 3,
      "links": [],
      "panels": [
        {
          "collapsed": false,
          "gridPos": {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 0
          },
          "id": 14,
          "panels": [],
          "title": "General Processing",
          "type": "row"
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 1
          },
          "id": 12,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "max",
              "fillBelowTo": "avg",
              "lines": false
            },
            {
              "alias": "avg",
              "fill": 0,
              "fillBelowTo": "min"
            },
            {
              "alias": "min",
              "lines": false
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "avg(sum(rate(hubble_flows_processed_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "avg",
              "refId": "A"
            },
            {
              "expr": "min(sum(rate(hubble_flows_processed_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "min",
              "refId": "B"
            },
            {
              "expr": "max(sum(rate(hubble_flows_processed_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "max",
              "refId": "C"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Flows processed Per Node",
          "tooltip": {
            "shared": true,
            "sort": 1,
            "value_type": "individual"
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
              "format": "ops",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 1
          },
          "id": 32,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_flows_processed_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, type)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{type}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Flows Types",
          "tooltip": {
            "shared": true,
            "sort": 2,
            "value_type": "individual"
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
              "format": "ops",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 6
          },
          "id": 59,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_flows_processed_total{type=\"L7\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, subtype)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{subtype}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "L7 Flow Distribution",
          "tooltip": {
            "shared": true,
            "sort": 2,
            "value_type": "individual"
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
              "format": "ops",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 6
          },
          "id": 60,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_flows_processed_total{type=\"Trace\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, subtype)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{subtype}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Trace Flow Distribution",
          "tooltip": {
            "shared": true,
            "sort": 2,
            "value_type": "individual"
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
              "format": "ops",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "collapsed": false,
          "gridPos": {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 11
          },
          "id": 16,
          "panels": [],
          "title": "Network",
          "type": "row"
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 12
          },
          "id": 33,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_flows_processed_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, verdict)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{verdict}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Forwarded vs Dropped",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "ops",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 12
          },
          "id": 8,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_drop_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, reason)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{reason}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Drop Reason",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 17
          },
          "id": 34,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": true,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum (rate(hubble_port_distribution_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, protocol)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{protocol}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Protocol Usage",
          "tooltip": {
            "shared": true,
            "sort": 2,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 17
          },
          "id": 18,
          "legend": {
            "alignAsTable": true,
            "avg": false,
            "current": true,
            "hideEmpty": false,
            "hideZero": false,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": true,
            "sort": "current",
            "sortDesc": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "topk(10, sum (rate(hubble_port_distribution_total{test_cluster_name=~\"\(_generatedName).*\",port!=\"0\"}[1m])) by (pod, port, protocol))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{port}}/{{protocol}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Top 10 Port Distribution",
          "tooltip": {
            "shared": true,
            "sort": 2,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 22
          },
          "id": 10,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "fin",
              "yaxis": 1
            },
            {
              "alias": "FIN",
              "yaxis": 2
            },
            {
              "alias": "RST",
              "yaxis": 2
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_tcp_flags_total{test_cluster_name=~\"\(_generatedName).*\",family=\"IPv4\"}[1m])) by (pod, flag)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{flag}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "TCPv4",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "alert": {
            "conditions": [
              {
                "evaluator": {
                  "params": [
                    0.2
                  ],
                  "type": "gt"
                },
                "operator": {
                  "type": "and"
                },
                "query": {
                  "params": [
                    "B",
                    "5m",
                    "now"
                  ]
                },
                "reducer": {
                  "params": [],
                  "type": "avg"
                },
                "type": "query"
              }
            ],
            "executionErrorState": "alerting",
            "for": "5m",
            "frequency": "1m",
            "handler": 1,
            "name": "Missing TCP SYN-ACK",
            "noDataState": "no_data",
            "notifications": []
          },
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 22
          },
          "id": 62,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "fin",
              "yaxis": 1
            },
            {
              "alias": "FIN",
              "yaxis": 2
            },
            {
              "alias": "RST",
              "yaxis": 2
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_tcp_flags_total{family=\"IPv4\", flag=\"SYN\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod) - sum(rate(hubble_tcp_flags_total{family=\"IPv4\", flag=\"SYN-ACK\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod)",
              "format": "time_series",
              "hide": false,
              "intervalFactor": 1,
              "legendFormat": "Missing SYN-ACK",
              "refId": "B"
            }
          ],
          "thresholds": [
            {
              "colorMode": "critical",
              "fill": true,
              "line": true,
              "op": "gt",
              "value": 0.2
            }
          ],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Missing TCPv4 SYN-ACKs",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 27
          },
          "id": 35,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "fin",
              "yaxis": 1
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_tcp_flags_total{family=\"IPv6\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, flag)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{flag}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "TCPv6",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "alert": {
            "conditions": [
              {
                "evaluator": {
                  "params": [
                    0.2
                  ],
                  "type": "gt"
                },
                "operator": {
                  "type": "and"
                },
                "query": {
                  "params": [
                    "B",
                    "5m",
                    "now"
                  ]
                },
                "reducer": {
                  "params": [],
                  "type": "avg"
                },
                "type": "query"
              }
            ],
            "executionErrorState": "alerting",
            "for": "5m",
            "frequency": "1m",
            "handler": 1,
            "name": "Missing TCPv6 SYN-ACKs alert",
            "noDataState": "no_data",
            "notifications": []
          },
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 27
          },
          "id": 63,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "fin",
              "yaxis": 1
            },
            {
              "alias": "FIN",
              "yaxis": 2
            },
            {
              "alias": "RST",
              "yaxis": 2
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_tcp_flags_total{family=\"IPv6\", flag=\"SYN\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod) - sum(rate(hubble_tcp_flags_total{family=\"IPv6\", flag=\"SYN-ACK\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod)",
              "format": "time_series",
              "hide": false,
              "intervalFactor": 1,
              "legendFormat": "Missing SYN-ACK",
              "refId": "B"
            }
          ],
          "thresholds": [
            {
              "colorMode": "critical",
              "fill": true,
              "line": true,
              "op": "gt",
              "value": 0.2
            }
          ],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Missing TCPv6 SYN-ACKs",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 32
          },
          "id": 31,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_icmp_total{family=\"IPv4\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, type)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{type}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "ICMPv4",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "alert": {
            "conditions": [
              {
                "evaluator": {
                  "params": [
                    0.1
                  ],
                  "type": "gt"
                },
                "operator": {
                  "type": "and"
                },
                "query": {
                  "params": [
                    "B",
                    "5m",
                    "now"
                  ]
                },
                "reducer": {
                  "params": [],
                  "type": "avg"
                },
                "type": "query"
              }
            ],
            "executionErrorState": "alerting",
            "for": "5m",
            "frequency": "1m",
            "handler": 1,
            "name": "Missing ICMPv4 Echo-Reply alert",
            "noDataState": "no_data",
            "notifications": []
          },
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 32
          },
          "id": 64,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_icmp_total{family=\"IPv4\", type=\"EchoRequest\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod) - sum(rate(hubble_icmp_total{family=\"IPv4\", type=\"EchoReply\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod)",

              "format": "time_series",
              "hide": false,
              "intervalFactor": 1,
              "legendFormat": "Missing ICMP Echo-Reply",
              "refId": "B"
            }
          ],
          "thresholds": [
            {
              "colorMode": "critical",
              "fill": true,
              "line": true,
              "op": "gt",
              "value": 0.1
            }
          ],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Missing ICMPv4 Echo-Reply",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 37
          },
          "id": 36,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_icmp_total{family=\"IPv6\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, type)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{type}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "ICMPv6",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 37
          },
          "id": 65,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_icmp_total{family=\"IPv6\", type=\"EchoRequest\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod) - sum(rate(hubble_icmp_total{family=\"IPv6\", type=\"EchoReply\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod)",

              "format": "time_series",
              "hide": false,
              "intervalFactor": 1,
              "legendFormat": "Missing ICMP Echo-Reply",
              "refId": "B"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Missing ICMPv6 Echo-Reply",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "collapsed": false,
          "gridPos": {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 42
          },
          "id": 42,
          "panels": [],
          "title": "Network Policy",
          "type": "row"
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 4,
            "w": 12,
            "x": 0,
            "y": 43
          },
          "id": 43,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_drop_total{reason=\"Policy denied (L3)\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, reason)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{reason}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Denies by Reason",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 4,
            "w": 12,
            "x": 12,
            "y": 43
          },
          "id": 61,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_drop_total{reason=\"Policy denied (L3)\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, protocol)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{protocol}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Denied Packets by Protocol",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 47
          },
          "id": 55,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "topk(10, sum(rate(hubble_drop_total{reason=\"Policy denied (L3)\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, source))",

              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{source}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Top 10 Source Pods with Denied Packets",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 47
          },
          "id": 54,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "topk(10, sum(rate(hubble_drop_total{reason=\"Policy denied (L3)\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, destination))",

              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{destination}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Top 10 Destination Pods with Denied Packets",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "pps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "collapsed": false,
          "gridPos": {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 52
          },
          "id": 47,
          "panels": [],
          "title": "HTTP",
          "type": "row"
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 6,
            "w": 12,
            "x": 0,
            "y": 53
          },
          "id": 45,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_http_requests_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, method)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{method}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "HTTP Requests",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 1,
              "format": "reqps",
              "label": "",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 6,
            "w": 12,
            "x": 12,
            "y": 53
          },
          "id": 49,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_http_responses_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, status)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{status}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "HTTP responses",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 1,
              "format": "reqps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 59
          },
          "id": 51,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "histogram_quantile(0.5, rate(hubble_http_request_duration_seconds_bucket{test_cluster_name=~\"\(_generatedName).*\"}[1m]))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{method}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "HTTP Request/Response Latency (p50)",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 59
          },
          "id": 58,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "histogram_quantile(0.99, rate(hubble_http_request_duration_seconds_bucket{test_cluster_name=~\"\(_generatedName).*\"}[1m]))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{method}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "HTTP Request/Response Latency (p99)",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 64
          },
          "id": 53,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": true,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_http_requests_total{test_cluster_name=~\"\(_generatedName).*\"}[5m])) by (pod, protocol)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{protocol}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "HTTP Protocol Usage",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 0,
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "collapsed": false,
          "gridPos": {
            "h": 1,
            "w": 24,
            "x": 0,
            "y": 69
          },
          "id": 6,
          "panels": [],
          "title": "DNS",
          "type": "row"
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 8,
            "x": 0,
            "y": 70
          },
          "id": 2,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_dns_queries_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, qtypes)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{qtypes}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "DNS Requests",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "reqps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 8,
            "x": 8,
            "y": 70
          },
          "id": 4,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_dns_responses_total{rcode=\"No Error\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, qtypes)",

              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{qtypes}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "DNS responses",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 1,
              "format": "reqps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "alert": {
            "conditions": [
              {
                "evaluator": {
                  "params": [
                    0.5
                  ],
                  "type": "gt"
                },
                "operator": {
                  "type": "and"
                },
                "query": {
                  "params": [
                    "A",
                    "5m",
                    "now"
                  ]
                },
                "reducer": {
                  "params": [],
                  "type": "avg"
                },
                "type": "query"
              }
            ],
            "executionErrorState": "alerting",
            "for": "5m",
            "frequency": "1m",
            "handler": 1,
            "name": "DNS Request/Response Symmetry alert",
            "noDataState": "no_data",
            "notifications": []
          },
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 8,
            "x": 16,
            "y": 70
          },
          "id": 66,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_dns_queries_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, qtypes) - sum(rate(hubble_dns_responses_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, qtypes)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{qtypes}}",
              "refId": "A"
            }
          ],
          "thresholds": [
            {
              "colorMode": "critical",
              "fill": true,
              "line": true,
              "op": "gt",
              "value": 0.5
            }
          ],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Missing DNS Responses",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "format": "reqps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 75
          },
          "id": 40,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_dns_response_types_total{test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, type)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{type}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "DNS Response Record Type",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 1,
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 75
          },
          "id": 57,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_dns_responses_total{rcode=\"No Error\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod,ips_returned)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{ips_returned}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "DNS Response IPs Returned",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 1,
              "format": "reqps",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 0,
            "y": 80
          },
          "id": 28,
          "legend": {
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "show": true,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "sum(rate(hubble_dns_responses_total{rcode!=\"No Error\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, qtypes, rcode)",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{rcode}} ({{qtypes}})",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "DNS Errors",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 4,
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 5,
            "w": 12,
            "x": 12,
            "y": 80
          },
          "id": 56,
          "legend": {
            "alignAsTable": false,
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "rightSide": false,
            "show": true,
            "sideWidth": null,
            "total": false,
            "values": false
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "topk(10,sum(rate(hubble_dns_responses_total{rcode!=\"No Error\",test_cluster_name=~\"\(_generatedName).*\"}[1m])) by (pod, destination))",

              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{destination}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Pods with DNS errors",
          "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
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
              "decimals": 4,
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "prometheus",
          "fill": 1,
          "gridPos": {
            "h": 6,
            "w": 24,
            "x": 0,
            "y": 85
          },
          "id": 20,
          "legend": {
            "alignAsTable": true,
            "avg": false,
            "current": true,
            "max": false,
            "min": false,
            "rightSide": true,
            "show": true,
            "sort": "current",
            "sortDesc": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "nullPointMode": "null",
          "options": {},
          "percentage": false,
          "pointradius": 2,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "topk(10, sum(rate(hubble_dns_queries_total{test_cluster_name=~\"\(_generatedName).*\"}[10m])*60) by (query, qtypes))",
              "format": "time_series",
              "intervalFactor": 1,
              "legendFormat": "{{query}} ({{qtypes}})",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Top 10 DNS Queries per minute",
          "tooltip": {
            "shared": true,
            "sort": 2,
            "value_type": "individual"
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
              "format": "none",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "refresh": "30s",
      "schemaVersion": 18,
      "style": "dark",
      "tags": [],
      "templating": {
        "list": []
      },
      "time": {
        "from": "now-6h",
        "to": "now"
      },
      "timepicker": {
        "refresh_intervals": [
          "5s",
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
      "timezone": "",
      "title": "Hubble",
      "uid": "5HftnJAWz",
      "version": 24
    }

resource: v1alpha2.#TestClusterGKE
