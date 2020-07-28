// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/isovalent/gke-test-cluster-management/operator/pkg/config"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterResources(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		c := &Config{
			BaseDirectory: "./nonexistent",
		}

		err := c.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unable to list avaliable config templates in "./nonexistent": open ./nonexistent: no such file or directory`))
	}

	{
		c := &Config{
			BaseDirectory: "./",
		}

		err := c.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`no config templates found in "./"`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		g.Expect(c.HaveExistingTemplate("basic")).To(BeTrue())
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderClusterCoreResourcesAsJSON(&v1alpha1.TestClusterGKE{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil/empty configTemplate`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		templateName := "basic"
		region := "us-west1"
		machineType := "n1-standard-4"

		defCluster := &v1alpha1.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
			},
			Spec: v1alpha1.TestClusterGKESpec{
				MachineType: &machineType,
				Location:    &region,
				Region:      &region,
			},
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		g.Expect(c.HaveExistingTemplate(templateName)).To(BeTrue())

		err = c.ApplyDefaults(templateName, defCluster)
		g.Expect(err).ToNot(HaveOccurred())

		err = c.ApplyDefaultsForClusterAccessResources(defCluster)
		g.Expect(err).ToNot(HaveOccurred())

		{
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "other",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					ConfigTemplate: &templateName,
					Region:         new(string),
					Location:       new(string),
				},
			}
			*cluster.Spec.Region = "europe-west2"
			*cluster.Spec.Location = "europe-west2-b"

			coreResourcesData, err := c.RenderClusterCoreResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			const coreResourcesExpected = `
			{
				"kind": "List",
				"apiVersion": "v1",
				"items": [
				  {
					"kind": "ContainerCluster",
					"apiVersion": "container.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/remove-default-node-pool": "true",
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "location": "europe-west2-b",
					  "initialNodeCount": 1,
					  "loggingService": "logging.googleapis.com/kubernetes",
					  "masterAuth": {
						"clientCertificateConfig": {
						  "issueClientCertificate": false
						}
					  },
					  "monitoringService": "monitoring.googleapis.com/kubernetes",
					  "networkRef": {
						"name": "baz"
					  },
					  "subnetworkRef": {
						"name": "baz"
					  }
					}
				  },
				  {
					"kind": "ContainerNodePool",
					"apiVersion": "container.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "location": "europe-west2-b",
					  "initialNodeCount": 0,
					  "clusterRef": {
						"name": "baz"
					  },
					  "management": {
						"autoRepair": false,
						"autoUpgrade": false
					  },
					  "nodeConfig": {
						"metadata": {
						  "disable-legacy-endpoints": "true"
						},
						"machineType": "n1-standard-4",
						"diskSizeGb": 100,
						"diskType": "pd-standard",
						"oauthScopes": [
						  "https://www.googleapis.com/auth/logging.write",
						  "https://www.googleapis.com/auth/monitoring"
						]
					  }
					}
				  },
				  {
					"kind": "ComputeNetwork",
					"apiVersion": "compute.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "autoCreateSubnetworks": false,
					  "deleteDefaultRoutesOnCreate": false,
					  "routingMode": "REGIONAL"
					}
				  },
				  {
					"kind": "ComputeSubnetwork",
					"apiVersion": "compute.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "region": "europe-west2",
					  "networkRef": {
						"name": "baz"
					  },
					  "ipCidrRange": "10.128.0.0/20"
					}
				  }
				]
			}`
			g.Expect(coreResourcesData).To(MatchJSON(coreResourcesExpected))

			accessResourcesData, err := c.RenderClusterAccessResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			const accessResourcesExpected = `
			{
				"kind": "List",
				"apiVersion": "v1",
				"items": [
				  {
					"kind": "IAMServiceAccount",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "displayName": "baz-admin"
					}
				  },
				  {
					"kind": "IAMPolicy",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "resourceRef": {
						"name": "baz-admin",
						"kind": "IAMServiceAccount",
						"apiVersion": "iam.cnrm.cloud.google.com/v1beta1"
					  },
					  "bindings": [
						{
						  "role": "roles/iam.workloadIdentityUser",
						  "members": [
							"serviceAccount:cilium-ci.svc.id.goog[other/baz-admin]"
						  ]
						}
					  ]
					}
				  },
				  {
					"kind": "ServiceAccount",
					"apiVersion": "v1",
					"metadata": {
					  "name": "baz-admin",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci",
						"iam.gke.io/gcp-service-account": "baz-admin@cilium-ci.iam.gserviceaccount.com"
					  }
					}
				  }
				]
			}`
			g.Expect(accessResourcesData).To(MatchJSON(accessResourcesExpected))

			objs, err := c.RenderAllClusterResources(cluster, false)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(7))
		}

		{
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
			}

			coreResourcesData, err := c.RenderClusterCoreResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			const coreResourcesExpected = `
			{
				"kind": "List",
				"apiVersion": "v1",
				"items": [
				  {
					"kind": "ContainerCluster",
					"apiVersion": "container.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/remove-default-node-pool": "true",
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "location": "us-west1",
					  "initialNodeCount": 1,
					  "loggingService": "logging.googleapis.com/kubernetes",
					  "masterAuth": {
						"clientCertificateConfig": {
						  "issueClientCertificate": false
						}
					  },
					  "monitoringService": "monitoring.googleapis.com/kubernetes",
					  "networkRef": {
						"name": "bar"
					  },
					  "subnetworkRef": {
						"name": "bar"
					  }
					}
				  },
				  {
					"kind": "ContainerNodePool",
					"apiVersion": "container.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "location": "us-west1",
					  "initialNodeCount": 0,
					  "clusterRef": {
						"name": "bar"
					  },
					  "management": {
						"autoRepair": false,
						"autoUpgrade": false
					  },
					  "nodeConfig": {
						"metadata": {
						  "disable-legacy-endpoints": "true"
						},
						"machineType": "n1-standard-4",
						"diskSizeGb": 100,
						"diskType": "pd-standard",
						"oauthScopes": [
						  "https://www.googleapis.com/auth/logging.write",
						  "https://www.googleapis.com/auth/monitoring"
						]
					  }
					}
				  },
				  {
					"kind": "ComputeNetwork",
					"apiVersion": "compute.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "autoCreateSubnetworks": false,
					  "deleteDefaultRoutesOnCreate": false,
					  "routingMode": "REGIONAL"
					}
				  },
				  {
					"kind": "ComputeSubnetwork",
					"apiVersion": "compute.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "region": "us-west1",
					  "networkRef": {
						"name": "bar"
					  },
					  "ipCidrRange": "10.128.0.0/20"
					}
				  }
				]
			}`
			g.Expect(coreResourcesData).To(MatchJSON(coreResourcesExpected))

			accessResourcesData, err := c.RenderClusterAccessResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			const accessResourcesExpected = `
			{
				"kind": "List",
				"apiVersion": "v1",
				"items": [
				  {
					"kind": "IAMServiceAccount",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "displayName": "bar-admin"
					}
				  },
				  {
					"kind": "IAMPolicy",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					},
					"spec": {
					  "resourceRef": {
						"name": "bar-admin",
						"kind": "IAMServiceAccount",
						"apiVersion": "iam.cnrm.cloud.google.com/v1beta1"
					  },
					  "bindings": [
						{
						  "role": "roles/iam.workloadIdentityUser",
						  "members": [
							"serviceAccount:cilium-ci.svc.id.goog[default/bar-admin]"
						  ]
						}
					  ]
					}
				  },
				  {
					"kind": "ServiceAccount",
					"apiVersion": "v1",
					"metadata": {
					  "name": "bar-admin",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci",
						"iam.gke.io/gcp-service-account": "bar-admin@cilium-ci.iam.gserviceaccount.com"
					  }
					}
				  }
				]
			}`
			g.Expect(accessResourcesData).To(MatchJSON(accessResourcesExpected))

			objs, err := c.RenderAllClusterResources(cluster, false)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(7))
		}

		{
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "baz",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
			}

			objs, err := c.RenderAllClusterResources(cluster, true)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(7))

			name := objs.Items[1].GetName()
			g.Expect(name).To(HavePrefix("baz-"))
			g.Expect(name).To(HaveLen(9))

			for _, obj := range objs.Items {
				labels := obj.GetLabels()
				g.Expect(labels).To(HaveKeyWithValue("cluster", name))
				g.Expect(obj.GetName()).To(HavePrefix(name))
			}
		}

		{
			clusterName := "baz-abc2sax"
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "baz",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
				Status: v1alpha1.TestClusterGKEStatus{
					ClusterName: &clusterName,
				},
			}

			objs, err := c.RenderAllClusterResources(cluster, true)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(7))

			name := objs.Items[1].GetName()
			g.Expect(name).To(Equal(clusterName))

			for _, obj := range objs.Items {
				labels := obj.GetLabels()
				g.Expect(labels).To(HaveKeyWithValue("cluster", name))
				g.Expect(obj.GetName()).To(HavePrefix(name))
			}
		}

		{
			cluster := &v1alpha1.TestClusterGKE{
				Spec: v1alpha1.TestClusterGKESpec{
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
			}

			_, err := c.RenderClusterCoreResourcesAsJSON(cluster)
			g.Expect(err).To(HaveOccurred())
			// this is another weird error from CUE, but that's what you get when optional field is unspecified on export...
			g.Expect(err.Error()).ToNot(Equal(`cue: marshal error: template.items.0.metadata.name: field "name" is optional`))
			g.Expect(err.Error()).To(Equal(`unexpected unnamed object`))
		}
	}
}

func TestTestRunnerJobResources(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderTestRunnerJobResourcesAsJSON(&v1alpha1.TestClusterGKE{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderTestRunnerJobResourcesAsJSON(&v1alpha1.TestClusterGKE{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderTestRunnerJobResourcesAsJSON(&v1alpha1.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			Spec: v1alpha1.TestClusterGKESpec{
				JobSpec: &v1alpha1.TestClusterGKEJobSpec{},
			}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil/empty runnerImage`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		runnerImage := "cilium-ci/cilium-e2e:latest"

		_, err = c.RenderTestRunnerJobResourcesAsJSON(&v1alpha1.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"},
			Spec: v1alpha1.TestClusterGKESpec{
				JobSpec: &v1alpha1.TestClusterGKEJobSpec{
					RunnerImage: &runnerImage,
				},
			}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil status.clusterName`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		runnerImage := "cilium-ci/cilium-e2e:latest"
		defCluster := &v1alpha1.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
			},
			Spec: v1alpha1.TestClusterGKESpec{
				JobSpec: &v1alpha1.TestClusterGKEJobSpec{
					RunnerImage: &runnerImage,
				},
			},
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		err = c.ApplyDefaultsForTestRunnerJobResources(defCluster)
		g.Expect(err).ToNot(HaveOccurred())

		{
			actualName := "baz-a0b1c2"
			runnerImage := "cilium-ci/cilium-e2e:80d4133f2b9317a0f08fcff9b2f8d625ea9f7b7a"
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "other",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					JobSpec: &v1alpha1.TestClusterGKEJobSpec{
						RunnerImage: &runnerImage,
					},
				},
				Status: v1alpha1.TestClusterGKEStatus{
					ClusterName: &actualName,
				},
			}

			data, err := c.RenderTestRunnerJobResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			const expected = `
			{
				"kind": "List",
				"apiVersion": "v1",
				"items": [
				  {
					"kind": "Job",
					"apiVersion": "batch/v1",
					"metadata": {
					  "name": "test-runner-baz-a0b1c2",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz-a0b1c2"
					  }
					},
					"spec": {
					  "backoffLimit": 0,
					  "template": {
						"metadata": {
						  "labels": {
							"cluster": "baz-a0b1c2"
						  }
						},
						"spec": {
						  "volumes": [],
						  "initContainers": [
							{
							  "name": "get-credentials",
							  "command": [
								"bash",
								"-c",
								"until gcloud auth list \"--format=value(account)\" | grep baz-a0b1c2-admin@cilium-ci.iam.gserviceaccount.com ; do sleep 1 ; done"
							  ],
							  "image": "google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a",
							  "volumeMounts": []
							}
						  ],
						  "containers": [
							{
							  "name": "test-runner",
							  "command": [
								"bash",
								"-l"
							  ],
							  "image": "cilium-ci/cilium-e2e:80d4133f2b9317a0f08fcff9b2f8d625ea9f7b7a",
							  "volumeMounts": [],
							  "tty": true
							}
						  ],
						  "dnsPolicy": "ClusterFirst",
						  "restartPolicy": "Never",
						  "serviceAccountName": "baz-a0b1c2-admin"
						}
					  }
					}
				  }
				]
			}`
			g.Expect(data).To(MatchJSON(expected))

			objs, err := c.RenderTestRunnerJobResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(1))
		}

		{
			actualName := "bar-0a1b2c"
			runnerImage := "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a"
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "bar",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					JobSpec: &v1alpha1.TestClusterGKEJobSpec{
						RunnerImage: &runnerImage,
					},
				},
				Status: v1alpha1.TestClusterGKEStatus{
					ClusterName: &actualName,
				},
			}

			data, err := c.RenderTestRunnerJobResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			const expected = `
			{
				"kind": "List",
				"apiVersion": "v1",
				"items": [
				  {
					"kind": "Job",
					"apiVersion": "batch/v1",
					"metadata": {
					  "name": "test-runner-bar-0a1b2c",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar-0a1b2c"
					  }
					},
					"spec": {
					  "backoffLimit": 0,
					  "template": {
						"metadata": {
						  "labels": {
							"cluster": "bar-0a1b2c"
						  }
						},
						"spec": {
						  "volumes": [],
						  "initContainers": [
							{
							  "name": "get-credentials",
							  "command": [
								"bash",
								"-c",
								"until gcloud auth list \"--format=value(account)\" | grep bar-0a1b2c-admin@cilium-ci.iam.gserviceaccount.com ; do sleep 1 ; done"
							  ],
							  "image": "google/cloud-sdk:slim@sha256:a2bade78228faad59a16c36d440f10cfef58a6055cd997d19e258c59c78a409a",
							  "volumeMounts": []
							}
						  ],
						  "containers": [
							{
							  "name": "test-runner",
							  "command": [
								"bash",
								"-l"
							  ],
							  "image": "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a",
							  "volumeMounts": [],
							  "tty": true
							}
						  ],
						  "dnsPolicy": "ClusterFirst",
						  "restartPolicy": "Never",
						  "serviceAccountName": "bar-0a1b2c-admin"
						}
					  }
					}
				  }
				]
			}`
			g.Expect(data).To(MatchJSON(expected))

			objs, err := c.RenderTestRunnerJobResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(1))
		}

		{
			actualName := "baz-x2a8332"
			runnerImage := "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a"
			cluster := &v1alpha1.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "baz",
				},
				Spec: v1alpha1.TestClusterGKESpec{
					JobSpec: &v1alpha1.TestClusterGKEJobSpec{
						RunnerImage: &runnerImage,
					},
				},
				Status: v1alpha1.TestClusterGKEStatus{
					ClusterName: &actualName,
				},
			}

			objs, err := c.RenderTestRunnerJobResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(1))

			name := objs.Items[0].GetName()
			g.Expect(name).To(Equal("test-runner-baz-x2a8332"))

			for _, obj := range objs.Items {
				labels := obj.GetLabels()
				g.Expect(labels).To(HaveKeyWithValue("cluster", "baz-x2a8332"))
				g.Expect(obj.GetName()).To(Equal(name))
			}
		}

		{
			runnerImage := "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a"
			cluster := &v1alpha1.TestClusterGKE{
				Spec: v1alpha1.TestClusterGKESpec{
					JobSpec: &v1alpha1.TestClusterGKEJobSpec{
						RunnerImage: &runnerImage,
					},
				},
			}

			_, err := c.RenderTestRunnerJobResourcesAsJSON(cluster)
			g.Expect(err).To(HaveOccurred())
			// this is another weird error from CUE, but that's what you get when optional field is unspecified on export...
			g.Expect(err.Error()).ToNot(Equal(`cue: marshal error: template.items.0.metadata.name: field "name" is optional`))
			g.Expect(err.Error()).To(Equal(`unexpected unnamed object`))
		}
	}
}
