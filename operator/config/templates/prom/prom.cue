// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package prom

import "encoding/yaml"

_name: "prom"

_promviewSDConfigs: [{
        role: "service"
        selectors: [{
                role: "service"
                label: "component=promview"
        }]
}]

_promviewRelabelConfigs: [{
        source_labels: [
                "__meta_kubernetes_namespace",
        ]
        target_label: "test_cluster_namespace"
}, {
        source_labels: [
                "__meta_kubernetes_service_label_cluster",
        ]
        target_label: "test_cluster_name"
}]

// having more then one item in `match[]` will be encoded using repeated keys in the URL,
// but promview (or rather `k8s.io/client-go`'s `ProxyGet`) doesn't work with repeated keys;
// it also seems plausible to separate each of the federated jobs
_configMapData: {
	global: scrape_interval: "15s"
	scrape_configs: [{
		job_name: "gke-test-cluster-operator-promview-metrics"
		kubernetes_sd_configs: _promviewSDConfigs
		relabel_configs: _promviewRelabelConfigs
        }, {
                job_name: "gke-test-cluster-operator-promview-federate-kubernetes-apiservers"
                kubernetes_sd_configs: _promviewSDConfigs
                honor_labels: true
                metrics_path: '/federate'
                params: {
                        "match[]": [ '{job="kubernetes-apiservers"}' ]
                }
                relabel_configs: _promviewRelabelConfigs
        }, {
                job_name: "gke-test-cluster-operator-promview-federate-envoy"
                kubernetes_sd_configs: _promviewSDConfigs
                honor_labels: true
                metrics_path: '/federate'
                params: {
                        "match[]": [ '{job="envoy"}' ]
                }
                relabel_configs: _promviewRelabelConfigs
        }, {
                job_name: "gke-test-cluster-operator-promview-federate-kubernetes-pods"
                kubernetes_sd_configs: _promviewSDConfigs
                honor_labels: true
                metrics_path: '/federate'
                params: {
                        "match[]": [ '{job="kubernetes-pods"}' ]
                }
                relabel_configs: _promviewRelabelConfigs
        }, {
                job_name: "gke-test-cluster-operator-promview-federate-kubernetes-nodes"
                kubernetes_sd_configs: _promviewSDConfigs
                honor_labels: true
                metrics_path: '/federate'
                params: {
                        "match[]": [ '{job="kubernetes-nodes"}'
                        ]
                }
                relabel_configs: _promviewRelabelConfigs
	}, {
		job_name: "gke-test-cluster-operator-promview-federate-cadvisor"
                kubernetes_sd_configs: _promviewSDConfigs
		honor_labels: true
		metrics_path: '/federate'
		params: {
			"match[]": [ '{job="cadvisor"}' ]
		}
                relabel_configs: _promviewRelabelConfigs
	}, {
		job_name: "kubernetes-apiservers"
		kubernetes_sd_configs: [{
			role: "endpoints"
		}]
		scheme: "https"
		tls_config: {
			ca_file:              "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
			insecure_skip_verify: true
		}
		bearer_token_file: "/var/run/secrets/kubernetes.io/serviceaccount/token"
		relabel_configs: [{
			source_labels: [
				"__meta_kubernetes_namespace",
				"__meta_kubernetes_service_name",
				"__meta_kubernetes_endpoint_port_name",
			]
			action: "keep"
			regex:  "default;kubernetes;https"
		}, {
			source_labels: [
				"__meta_kubernetes_namespace",
			]
			target_label: "kubernetes_namespace"
		}, {
			source_labels: [
				"__meta_kubernetes_endpoints_name",
			]
			target_label: "_kube_service"
		}]
	}, {
		job_name:     "envoy"
		metrics_path: "/stats/prometheus"
		kubernetes_sd_configs: [{
			role: "pod"
		}]
		relabel_configs: [{
			source_labels: [
				"__meta_kubernetes_pod_container_name",
			]
			action: "keep"
			regex:  "^envoy$"
		}, {
			source_labels: [
				"__address__",
				"__meta_kubernetes_pod_annotation_prometheus_io_port",
			]
			action:       "replace"
			target_label: "__address__"
			regex:        "([^:]+)(?::\\\\d+)?;(\\\\d+)"
			replacement:  "${1}:9901"
		}, {
			source_labels: [
				"__meta_kubernetes_namespace",
			]
			action:       "replace"
			target_label: "kubernetes_namespace"
		}, {
			source_labels: [
				"__meta_kubernetes_pod_name",
			]
			action:       "replace"
			target_label: "kubernetes_pod_name"
		}]
		metric_relabel_configs: [{
			source_labels: [
				"cluster_name",
			]
			regex:  "(outbound|inbound|prometheus_stats).*"
			action: "drop"
		}, {
			source_labels: [
				"tcp_prefix",
			]
			regex:  "(outbound|inbound|prometheus_stats).*"
			action: "drop"
		}, {
			source_labels: [
				"listener_address",
			]
			regex:  "(.+)"
			action: "drop"
		}, {
			source_labels: [
				"http_conn_manager_listener_prefix",
			]
			regex:  "(.+)"
			action: "drop"
		}, {
			source_labels: [
				"http_conn_manager_prefix",
			]
			regex:  "(.+)"
			action: "drop"
		}, {
			source_labels: [
				"__name__",
			]
			regex:  "envoy_tls.*"
			action: "drop"
		}, {
			source_labels: [
				"__name__",
			]
			regex:  "envoy_tcp_downstream.*"
			action: "drop"
		}, {
			source_labels: [
				"__name__",
			]
			regex:  "envoy_http_(stats|admin).*"
			action: "drop"
		}, {
			source_labels: [
				"__name__",
			]
			regex:  "envoy_cluster_(lb|retry|bind|internal|max|original).*"
			action: "drop"
		}]
	}, {
		job_name: "kubernetes-pods"
		kubernetes_sd_configs: [{
			role: "pod"
		}]
		relabel_configs: [{
			source_labels: [
				"__meta_kubernetes_pod_annotation_prometheus_io_scrape",
			]
			action: "drop"
			regex:  "false"
		}, {
			source_labels: [
				"__meta_kubernetes_pod_annotation_prometheus_io_scheme",
			]
			action:       "replace"
			target_label: "__scheme__"
			regex:        "^(https?)$"
			replacement:  "$1"
		}, {
			source_labels: [
				"__meta_kubernetes_pod_annotation_prometheus_io_path",
			]
			action:       "replace"
			target_label: "__metrics_path__"
			regex:        "^(.+)$"
			replacement:  "$1"
		}, {
			source_labels: [
				"__address__",
				"__meta_kubernetes_pod_annotation_prometheus_io_port",
			]
			action:       "replace"
			target_label: "__address__"
			regex:        "([^:]+)(?::\\\\d+)?;(\\\\d+)"
			replacement:  "$1:$2"
		}, {
			source_labels: [
				"__meta_kubernetes_pod_container_port_number",
				"__meta_kubernetes_pod_annotation_prometheus_io_port",
			]
			separator: ""
			action:    "drop"
			regex:     "^$"
		}, {
			source_labels: [
				"__meta_kubernetes_namespace",
			]
			target_label: "kubernetes_namespace"
		}, {
			source_labels: [
				"__meta_kubernetes_pod_name",
			]
			target_label: "kubernetes_pod_name"
		}, {
			source_labels: [
				"__meta_kubernetes_pod_name",
				"__meta_kubernetes_pod_node_name",
			]
			target_label: "node"
			regex:        "^prom-node-exporter-.+;(.+)$"
			replacement:  "$1"
		}, {
			source_labels: [
				"_kube_service",
				"__meta_kubernetes_pod_name",
			]
			target_label: "_kube_service"
			regex:        "^;(kube-.*)-(?:ip|gke)-.*$"
			replacement:  "$1"
		}, {
			source_labels: [
				"_kube_service",
				"__meta_kubernetes_pod_name",
			]
			target_label: "_kube_service"
			regex:        "^;(.*?)(?:(?:-[0-9bcdf]+)?-[0-9a-z]{5}|-[0-9]+)?$"
			replacement:  "$1"
		}, {
			source_labels: [
				"_kube_service",
				"__meta_kubernetes_pod_name",
			]
			regex:        "^;(.+)$"
			target_label: "_kube_service"
			replacement:  "$1"
		}]
	}, {
		job_name: "kubernetes-nodes"
		kubernetes_sd_configs: [{
			role: "node"
		}]
		tls_config: insecure_skip_verify: true
		bearer_token_file: "/var/run/secrets/kubernetes.io/serviceaccount/token"
		relabel_configs: [{
			target_label: "__scheme__"
			replacement:  "https"
		}, {
			target_label: "kubernetes_namespace"
			replacement:  "default"
		}, {
			target_label: "_kube_service"
			replacement:  "kubelet"
		}, {
			target_label: "__address__"
			replacement:  "kubernetes.default.svc:443"
		}, {
			source_labels: [
				"__meta_kubernetes_node_name",
			]
			regex:        "(.+)"
			target_label: "__metrics_path__"
			replacement:  "/api/v1/nodes/${1}/proxy/metrics"
		}]
	}, {
		job_name: "cadvisor"
		kubernetes_sd_configs: [{
			role: "node"
		}]
		tls_config: insecure_skip_verify: true
		bearer_token_file: "/var/run/secrets/kubernetes.io/serviceaccount/token"
		scheme:            "https"
		relabel_configs: [{
			target_label: "kubernetes_namespace"
			replacement:  "default"
		}, {
			target_label: "_kube_service"
			replacement:  "cadvisor"
		}, {
			target_label: "__address__"
			replacement:  "kubernetes.default.svc:443"
		}, {
			source_labels: [
				"__meta_kubernetes_node_name",
			]
			regex:        "(.+)"
			target_label: "__metrics_path__"
			replacement:  "/api/v1/nodes/${1}/proxy/metrics/cadvisor"
		}]
		metric_relabel_configs: [{
			source_labels: [
				"__name__",
			]
			regex:  "container_(network_tcp_usage_total|network_udp_usage_total|tasks_state|cpu_load_average_10s)"
			action: "drop"
		}, {
			source_labels: [
				"__name__",
				"id",
			]
			regex:  "^container_.*;/system.slice/run-.*scope$"
			action: "drop"
		}, {
			source_labels: [
				"_kube_pod_name",
				"pod_name",
			]
			target_label: "_kube_pod_name"
			regex:        "^;(kube-.*)-(?:ip|gke)-.*$"
			replacement:  "$1"
		}, {
			source_labels: [
				"_kube_pod_name",
				"pod_name",
			]
			target_label: "_kube_pod_name"
			regex:        "^;(.*?)(?:(?:-[0-9bcdf]+)?-[0-9a-z]{5}|-[0-9]+)?$"
			replacement:  "$1"
		}, {
			source_labels: [
				"_kube_pod_name",
				"pod_name",
			]
			regex:        "^;(.+)$"
			target_label: "_kube_pod_name"
			replacement:  "$1"
		}]
	}]
}

#PromManifestTemplate: {
	apiVersion: "v1"
	kind:       "List"
	items: [{
		apiVersion: "v1"
		kind:       "Namespace"
		metadata: name: _name
	}, {
		apiVersion: "v1"
		kind:       "Service"
		metadata: {
			name: _name
			labels: name: _name
			namespace: _name
		}
		spec: {
			ports: [{
				name:       _name
				port:       80
				protocol:   "TCP"
				targetPort: 8080
			}]
			selector: name: _name
		}
	}, {
		apiVersion: "v1"
		kind:       "ServiceAccount"
		metadata: {
			name: _name
			labels: name: _name
			namespace: _name
		}
	}, {
		apiVersion: "rbac.authorization.k8s.io/v1"
		kind:       "ClusterRole"
		metadata: {
			name: _name
			labels: name: _name
		}
		rules: [{
			apiGroups: [
				"",
			]
			resources: [
				"configmaps",
				"endpoints",
				"limitranges",
				"namespaces",
				"nodes",
				"nodes/proxy",
				"persistentvolumeclaims",
				"persistentvolumes",
				"pods",
				"replicationcontrollers",
				"resourcequotas",
				"secrets",
				"services",
			]
			verbs: [
				"get",
				"list",
				"watch",
			]
		}, {
			apiGroups: [
				"extensions",
			]
			resources: [
				"daemonsets",
				"deployments",
				"ingresses",
				"replicasets",
			]
			verbs: [
				"get",
				"list",
				"watch",
			]
		}, {
			apiGroups: [
				"apps",
			]
			resources: [
				"daemonsets",
				"deployments",
				"replicasets",
				"statefulsets",
			]
			verbs: [
				"get",
				"list",
				"watch",
			]
		}, {
			apiGroups: [
				"batch",
			]
			resources: [
				"cronjobs",
				"jobs",
			]
			verbs: [
				"get",
				"list",
				"watch",
			]
		}, {
			apiGroups: [
				"autoscaling",
			]
			resources: [
				"horizontalpodautoscalers",
			]
			verbs: [
				"get",
				"list",
				"watch",
			]
		}, {
			apiGroups: [
				"policy",
			]
			resources: [
				"poddisruptionbudgets",
			]
			verbs: [
				"get",
				"list",
				"watch",
			]
		}, {
			nonResourceURLs: [
				"/metrics",
			]
			verbs: [
				"get",
			]
		}]
	}, {
		apiVersion: "rbac.authorization.k8s.io/v1"
		kind:       "ClusterRoleBinding"
		metadata: {
			name: _name
			labels: name: _name
		}
		roleRef: {
			kind:     "ClusterRole"
			name:     _name
			apiGroup: "rbac.authorization.k8s.io"
		}
		subjects: [{
			kind:      "ServiceAccount"
			name:      _name
			namespace: _name
		}]
	}, {
		apiVersion: "apps/v1"
		kind:       "Deployment"
		metadata: {
			name: _name
			labels: name: _name
			namespace: _name
		}
		spec: {
			replicas:             1
			revisionHistoryLimit: 2
			selector: matchLabels: name: _name
			template: {
				metadata: {
					annotations: "prometheus.io.scrape": "true"
					labels: name:                        _name
				}
				spec: {
					containers: [{
						name: _name
						args: [
							"--config.file=/etc/prometheus/prometheus.yml",
							"--web.listen-address=:8080",
							"--storage.tsdb.retention=2h",
							"--web.enable-lifecycle",
						]
						env: []
						image:           "docker.io/prom/prometheus:v2.20.1"
						imagePullPolicy: "IfNotPresent"
						ports: [{
							containerPort: 8080
							protocol:      "TCP"
						}]
						resources: requests: {
							cpu:    "200m"
							memory: "1000Mi"
						}
						volumeMounts: [{
							name:      "etc-prometheus"
							mountPath: "/etc/prometheus"
						}]
					}]
					serviceAccountName: _name
					volumes: [{
						name: "etc-prometheus"
						configMap: name: _name
					}]
				}
			}
		}
	}, {
		apiVersion: "apps/v1"
		kind:       "DaemonSet"
		metadata: {
			name: "prom-node-exporter"
			labels: name: "prom-node-exporter"
			namespace: _name
		}
		spec: {
			minReadySeconds: 5
			selector: matchLabels: name: "prom-node-exporter"
			template: {
				metadata: {
					annotations: "prometheus.io.scrape": "true"
					labels: name:                        "prom-node-exporter"
				}
				spec: {
					containers: [{
						name: "prom-node-exporter"
						env: []
						image:           "docker.io/prom/node-exporter:v1.0.1"
						imagePullPolicy: "IfNotPresent"
						ports: [{
							containerPort: 9100
							protocol:      "TCP"
						}]
						resources: requests: {
							cpu:    "10m"
							memory: "20Mi"
						}
						securityContext: privileged: true
					}]
					hostNetwork:        true
					hostPID:            true
					serviceAccountName: _name
					tolerations: [{
						effect:   "NoSchedule"
						operator: "Exists"
					}, {
						effect:   "NoExecute"
						operator: "Exists"
					}]
				}
			}
			updateStrategy: type: "RollingUpdate"
		}
	}, {
		apiVersion: "apps/v1"
		kind:       "Deployment"
		metadata: {
			name: "kube-state-metrics"
			labels: name: "kube-state-metrics"
			namespace: _name
		}
		spec: {
			replicas:             1
			revisionHistoryLimit: 2
			selector: matchLabels: name: "kube-state-metrics"
			template: {
				metadata: {
					annotations: "prometheus.io.scrape": "true"
					labels: name:                        "kube-state-metrics"
				}
				spec: {
					containers: [{
						name: "kube-state-metrics"
						args: [
							"--collectors=cronjobs,daemonsets,deployments,endpoints,horizontalpodautoscalers,ingresses,jobs,limitranges,namespaces,nodes,persistentvolumeclaims,persistentvolumes,poddisruptionbudgets,pods,resourcequotas,services,statefulsets",
						]
						env: []
						image: "quay.io/coreos/kube-state-metrics:v1.9.7"
						ports: [{
							name:          "metrics"
							containerPort: 8080
						}]
						resources: requests: {
							cpu:    "10m"
							memory: "20Mi"
						}
					}]
					serviceAccountName: _name
				}
			}
		}
	}, {
		apiVersion: "v1"
		kind:       "ConfigMap"
		metadata: {
			name: _name
			labels: name: _name
			namespace: _name
		}
		data: "prometheus.yml": yaml.Marshal(_configMapData)
	}]
}

template: #PromManifestTemplate
