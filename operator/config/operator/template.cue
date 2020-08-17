constants: {
	name:    "gke-test-cluster-operator"
	project: "cilium-ci"
}

_workload: {
	metadata: {
		name: "\(constants.name)"
		labels: name: "\(constants.name)"
		namespace: "\(parameters.namespace)"
		annotations: "fluxcd.io/automated": "true"
	}
	spec: _workloadSpec
}

_extra_rbac_ClusterRoleAndBinding: [...{}]

_workloadSpec: {
	template: {
		metadata: labels: name: "\(constants.name)"
		spec: {
			serviceAccount: "\(constants.name)"
			volumes: [{
				name: "tmp"
				emptyDir: {}
			}]
			containers: [{
				command: _command
				image:   "\(parameters.image)"
				name:    "operator"
				volumeMounts: [{
					name:      "tmp"
					mountPath: "/tmp"
				}]
				resources: {
					limits: {
						cpu:    "100m"
						memory: "100Mi"
					}
					requests: {
						cpu:    "100m"
						memory: "100Mi"
					}
				}
			}]
			terminationGracePeriodSeconds: 10
		}
	}
}

_command: [...string]
_optionalLogviewDomainFlag: [...string]

if !parameters.test {
	_workload: {
		apiVersion: "apps/v1"
		kind:       "Deployment"
	}
	_workloadSpec: {
		replicas: 1
		selector: matchLabels: name: "\(constants.name)"
		template: {
			metadata: labels: name: "\(constants.name)"
			spec: {
				containers: [{
					env: [{
						name: "GITHUB_TOKEN"
						valueFrom:
							secretKeyRef: {
								optional: true
								name:     "\(constants.name)-github-token"
								key:      "token"
							}
					}]
				}]
			}
		}
	}
	_command: [
			"/usr/bin/\(constants.name)",
			"--enable-leader-election",
	] + _optionalLogviewDomainFlag
}

if parameters.logviewDomain != null && len(parameters.logviewDomain) > 0{
	_optionalLogviewDomainFlag: [
		"--logview-domain=\(parameters.logviewDomain)",
	]
}

if parameters.test {
	_workload: {
		apiVersion: "batch/v1"
		kind:       "Job"
	}

	_workloadSpec: {
		backoffLimit: 0
		template: spec: {
			restartPolicy: "Never"
			containers: [{
				name: "operator"
				env: [{
					name: "NAMESPACE"
					valueFrom: fieldRef: fieldPath: "metadata.namespace"
				}]
			}]
		}
	}
	_command: [
		"test.gke-test-cluster-operator-controllers",
		"-test.v",
		"-test.timeout=5m",
		"-resource-prefix=\(parameters.namespace)",
	]

	// in principle, there should be no difference in RBAC configuration
	// between test and regular deployments, but it's currently inevitable
	// to allow creation/deletion of namespaces while testing the operator
	_extra_rbac_ClusterRoleAndBinding: [
		{
			apiVersion: "rbac.authorization.k8s.io/v1"
			kind:       "ClusterRole"
			metadata:
				name: "\(parameters.namespace)-\(constants.name)-extra"
			rules: [{
				apiGroups: [""]
				resources: ["namespaces"]
				verbs: ["create", "delete"]
			}]
		}, {
			apiVersion: "rbac.authorization.k8s.io/v1beta1"
			kind:       "ClusterRoleBinding"
			metadata: {
				name: "\(parameters.namespace)-\(constants.name)-extra"
				labels: name: "\(constants.name)"
			}
			roleRef: {
				kind:     "ClusterRole"
				name:     "\(parameters.namespace)-\(constants.name)-extra"
				apiGroup: "rbac.authorization.k8s.io"
			}
			subjects: [{
				kind:      "ServiceAccount"
				name:      "\(constants.name)"
				namespace: "\(parameters.namespace)"
			}]
		},
	]
}

_adminServiceAccountEmail: "\(constants.name)@\(constants.project).iam.gserviceaccount.com"
_adminServiceAccountRef:   "serviceAccount:\(constants.project).svc.id.goog[\(parameters.namespace)/\(constants.name)]"

_serviceAccount: {
	apiVersion: "v1"
	kind:       "ServiceAccount"
	metadata: {
		name: "\(constants.name)"
		labels: name: "\(constants.name)"
		namespace: "\(parameters.namespace)"
		annotations: {
			"iam.gke.io/gcp-service-account": "\(_adminServiceAccountEmail)"
		}
	}
}

_rbac_ClusterRoleBinding: {
	apiVersion: "rbac.authorization.k8s.io/v1beta1"
	kind:       "ClusterRoleBinding"
	metadata: {
		name: "\(parameters.namespace)-\(constants.name)"
		labels: name: "\(constants.name)"
	}
	roleRef: {
		kind:     "ClusterRole"
		name:     "\(constants.name)"
		apiGroup: "rbac.authorization.k8s.io"
	}
	subjects: [{
		kind:      "ServiceAccount"
		name:      "\(constants.name)"
		namespace: "\(parameters.namespace)"
	}]
}

_iam_clusterAdminAccess: [
	{
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMServiceAccount"
		metadata: {
			name:      "\(constants.name)"
			namespace: "\(parameters.namespace)"
			labels: name: "\(constants.name)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(constants.project)"
			}
		}
	},
	{
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMPolicyMember"
		metadata: {
			name: "\(constants.name)-workload-identity"
			labels: name: "\(constants.name)"
			namespace: "\(parameters.namespace)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(constants.project)"
			}
		}
		spec: {
			member: "\(_adminServiceAccountRef)"
			role:   "roles/iam.workloadIdentityUser"
			resourceRef: {
				apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
				kind:       "IAMServiceAccount"
				name:       "\(constants.name)"
			}
		}
	},
	{
		apiVersion: "iam.cnrm.cloud.google.com/v1beta1"
		kind:       "IAMPolicyMember"
		metadata: {
			name: "\(constants.name)-cluster-admin"
			labels: name: "\(constants.name)"
			namespace: "\(parameters.namespace)"
			annotations: {
				"cnrm.cloud.google.com/project-id": "\(constants.project)"
			}
		}
		spec: {
			member: "serviceAccount:\(_adminServiceAccountEmail)"
			role:   "roles/container.admin" // "roles/owner"
			resourceRef: {
				// At the moment ContainerCluster cannot be referenced here, so it's at project level for now
				// (see https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/248)
				apiVersion: "resourcemanager.cnrm.cloud.google.com/v1beta1"
				kind:       "Project"
				external:   "projects/cilium-ci"
			}
		}

	},
]

_extra_certManager_IssuerAndCertificater: [...{}]

if parameters.certManager {
	_extra_certManager_IssuerAndCertificater: [{
		apiVersion: "cert-manager.io/v1beta1"
		kind:       "Issuer"
		metadata: {
			name: constants.name
			labels: name: constants.name
			namespace: parameters.namespace
		}
		spec: selfSigned: {}
	}, {
		apiVersion: "cert-manager.io/v1beta1"
		kind:       "Certificate"
		metadata: {
			name: constants.name
			labels: name: constants.name
			namespace: parameters.namespace
		}
		spec: {
			dnsNames: [
				"\(constants.name)",
				"\(constants.name).\(parameters.namespace).svc",
				"\(constants.name).\(parameters.namespace).svc.cluster.local",
			]
			issuerRef: {
				kind: "Issuer"
				name: "\(constants.name)"
			}
			secretName: "\(constants.name)-webhook-server-cert"
		}
	}]
}

// TODO: webhook - these were provided by kubebuilder as kustomize pathches,
// and need to be ported over to CUE and tested

// // The following patch adds a directive for certmanager to inject CA into the CRD
// // CRD conversion requires k8s 1.13 or later.
// apiVersion: "apiextensions.k8s.io/v1beta1"
// kind:       "CustomResourceDefinition"
// metadata: {
// 	annotations: "cert-manager.io/inject-ca-from": "$(CERTIFICATE_NAMESPACE)/$(CERTIFICATE_NAME)"
// 	name: "testclustergkes.clusters.ci.cilium.io"
// }
// // The following patch adds a directive for certmanager to inject CA into the CRD
// // CRD conversion requires k8s 1.13 or later.
// apiVersion: "apiextensions.k8s.io/v1beta1"
// kind:       "CustomResourceDefinition"
// metadata: {
// 	annotations: "cert-manager.io/inject-ca-from": "$(CERTIFICATE_NAMESPACE)/$(CERTIFICATE_NAME)"
// 	name: "testclusterpoolgkes.clusters.ci.cilium.io"
// }
// // The following patch enables conversion webhook for CRD
// // CRD conversion requires k8s 1.13 or later.
// apiVersion: "apiextensions.k8s.io/v1beta1"
// kind:       "CustomResourceDefinition"
// metadata: name: "testclustergkes.clusters.ci.cilium.io"
// spec: conversion: {
// 	strategy: "Webhook"
// 	webhookClientConfig: {
// 		// this is "\n" used as a placeholder, otherwise it will be rejected by the apiserver for being blank,
// 		// but we're going to set it later using the cert-manager (or potentially a patch if not using cert-manager)
// 		caBundle: "Cg=="
// 		service: {
// 			namespace: "system"
// 			name:      "webhook-service"
// 			path:      "/convert"
// 		}
// 	}
// }
// // The following patch enables conversion webhook for CRD
// // CRD conversion requires k8s 1.13 or later.
// apiVersion: "apiextensions.k8s.io/v1beta1"
// kind:       "CustomResourceDefinition"
// metadata: name: "testclusterpoolgkes.clusters.ci.cilium.io"
// spec: conversion: {
// 	strategy: "Webhook"
// 	webhookClientConfig: {
// 		// this is "\n" used as a placeholder, otherwise it will be rejected by the apiserver for being blank,
// 		// but we're going to set it later using the cert-manager (or potentially a patch if not using cert-manager)
// 		caBundle: "Cg=="
// 		service: {
// 			namespace: "system"
// 			name:      "webhook-service"
// 			path:      "/convert"
// 		}
// 	}
// }

// [{
// 	// Adds namespace to all resources.
// 	namespace: "operator-system"

// 	// Value of this field is prepended to the
// 	// names of all resources, e.g. a deployment named
// 	// "wordpress" becomes "alices-wordpress".
// 	// Note that it should also match with the prefix (text before '-') of the namespace
// 	// field above.
// 	namePrefix: "operator-"

// 	// Labels to add to all resources and selectors.
// 	//commonLabels:
// 	//  someName: someValue

// 	bases: [
// 		"../crd",
// 		"../rbac",
// 		"../manager",
// 	]
// 	// [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in 
// 	// crd/kustomization.yaml
// 	//- ../webhook
// 	// [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
// 	//- ../certmanager
// 	// [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'. 
// 	//- ../prometheus

// 	patchesStrategicMerge:
// 	// Protect the /metrics endpoint by putting it behind auth.
// 	// If you want your controller-manager to expose the /metrics
// 	// endpoint w/o any authn/z, please comment the following line.
// 	[
// 		"manager_auth_proxy_patch.yaml",
// 	]

// 	// [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in 
// 	// crd/kustomization.yaml
// 	//- manager_webhook_patch.yaml
// 	// [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
// 	// Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
// 	// 'CERTMANAGER' needs to be enabled to use ca injection
// 	//- webhookcainjection_patch.yaml
// 	// the following config is for teaching kustomize how to do var substitution
// 	vars: null
// }]
// [{
// 	// This patch inject a sidecar container which is a HTTP proxy for the 
// 	// controller manager, it performs RBAC authorization against the Kubernetes API using SubjectAccessReviews.
// 	apiVersion: "apps/v1"
// 	kind:       "Deployment"
// 	metadata: {
// 		name:      "controller-manager"
// 		namespace: "system"
// 	}
// 	spec: template: spec: containers: [{
// 		name:  "kube-rbac-proxy"
// 		image: "gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0"
// 		args: [
// 			"--secure-listen-address=0.0.0.0:8443",
// 			"--upstream=http://127.0.0.1:8080/",
// 			"--logtostderr=true",
// 			"--v=10",
// 		]
// 		ports: [{
// 			containerPort: 8443
// 			name:          "https"
// 		}]
// 	}, {
// 		name: "manager"
// 		args: [
// 			"--metrics-addr=127.0.0.1:8080",
// 			"--enable-leader-election",
// 		]
// 	}]
// }]
// [{
// 	apiVersion: "apps/v1"
// 	kind:       "Deployment"
// 	metadata: {
// 		name:      "controller-manager"
// 		namespace: "system"
// 	}
// 	spec: template: spec: {
// 		containers: [{
// 			name: "manager"
// 			ports: [{
// 				containerPort: 9443
// 				name:          "webhook-server"
// 				protocol:      "TCP"
// 			}]
// 			volumeMounts: [{
// 				mountPath: "/tmp/k8s-webhook-server/serving-certs"
// 				name:      "cert"
// 				readOnly:  true
// 			}]
// 		}]
// 		volumes: [{
// 			name: "cert"
// 			secret: {
// 				defaultMode: 420
// 				secretName:  "webhook-server-cert"
// 			}
// 		}]
// 	}
// }]
// [{
// 	// This patch add annotation to admission webhook config and
// 	// the variables $(CERTIFICATE_NAMESPACE) and $(CERTIFICATE_NAME) will be substituted by kustomize.
// 	apiVersion: "admissionregistration.k8s.io/v1beta1"
// 	kind:       "MutatingWebhookConfiguration"
// 	metadata: {
// 		name: "mutating-webhook-configuration"
// 		annotations: "cert-manager.io/inject-ca-from": "$(CERTIFICATE_NAMESPACE)/$(CERTIFICATE_NAME)"
// 	}
// }, {
// 	apiVersion: "admissionregistration.k8s.io/v1beta1"
// 	kind:       "ValidatingWebhookConfiguration"
// 	metadata: {
// 		name: "validating-webhook-configuration"
// 		annotations: "cert-manager.io/inject-ca-from": "$(CERTIFICATE_NAMESPACE)/$(CERTIFICATE_NAME)"
// 	}
// }]

_core_items: [
	_serviceAccount,
	_workload,
	_rbac_ClusterRoleBinding,
]

#WorkloadTemplate: {
	kind:       "List"
	apiVersion: "v1"
	items:
		_core_items +
		_iam_clusterAdminAccess +
		_extra_rbac_ClusterRoleAndBinding +
		_extra_certManager_IssuerAndCertificater +
		[]
}

#WorkloadParameters: {
	namespace:      string
	image:          string
	test:           bool
	logviewDomain?: string
	certManager: bool
}

parameters: #WorkloadParameters
template:   #WorkloadTemplate
