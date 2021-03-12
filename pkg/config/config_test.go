// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/isovalent/gke-test-cluster-operator/pkg/config"

	"github.com/isovalent/gke-test-cluster-operator/api/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

		_, err = c.RenderClusterCoreResourcesAsJSON(&v1alpha2.TestClusterGKE{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil/empty configTemplate`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		templateName := "basic"
		region := "us-west1"
		nodes := 3
		machineType := "n1-standard-4"

		defCluster := &v1alpha2.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
			},
			Spec: v1alpha2.TestClusterGKESpec{
				Nodes:       &nodes,
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
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "other",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					ConfigTemplate: &templateName,
					Project:        new(string),
					Region:         new(string),
					Location:       new(string),
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: new(string),
				},
			}
			*cluster.Spec.Project = "cilium-ci"
			*cluster.Spec.Region = "europe-west2"
			*cluster.Spec.Location = "europe-west2-b"
			*cluster.Status.ClusterName = "baz"

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
					  "initialNodeCount": 3,
					  "clusterRef": {
						"name": "baz"
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
					  "name": "baz-admin",
					  "namespace": "other",
					  "labels": {
						"cluster": "baz"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					}
				  },
				  {
					"kind": "IAMPolicyMember",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz-workload-identity",
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
					  "role": "roles/iam.workloadIdentityUser",
					  "member":	"serviceAccount:cilium-ci.svc.id.goog[other/baz-admin]"
					}
				  },
				  {
					"kind": "IAMPolicyMember",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "baz-cluster-admin",
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
						"external": "projects/cilium-ci",
						"kind": "Project",
						"apiVersion": "resourcemanager.cnrm.cloud.google.com/v1beta1"
					  },
					  "role": "roles/container.clusterAdmin",
					  "member":	"serviceAccount:baz-admin@cilium-ci.iam.gserviceaccount.com"
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

			objs, err := c.RenderAllClusterResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(9))
		}

		{
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					Project:        new(string),
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: new(string),
				},
			}
			*cluster.Spec.Project = "cilium-ci"
			*cluster.Status.ClusterName = "bar"

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
					  "initialNodeCount": 3,
					  "clusterRef": {
						"name": "bar"
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
					  "name": "bar-admin",
					  "namespace": "default",
					  "labels": {
						"cluster": "bar"
					  },
					  "annotations": {
						"cnrm.cloud.google.com/project-id": "cilium-ci"
					  }
					}
				  },
				  {
					"kind": "IAMPolicyMember",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar-workload-identity",
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
					  "role": "roles/iam.workloadIdentityUser",
					  "member": "serviceAccount:cilium-ci.svc.id.goog[default/bar-admin]"
					}
				  },
				  {
					"kind": "IAMPolicyMember",
					"apiVersion": "iam.cnrm.cloud.google.com/v1beta1",
					"metadata": {
					  "name": "bar-cluster-admin",
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
						"external": "projects/cilium-ci",
						"kind": "Project",
						"apiVersion": "resourcemanager.cnrm.cloud.google.com/v1beta1"
					  },
					  "role": "roles/container.clusterAdmin",
					  "member":	"serviceAccount:bar-admin@cilium-ci.iam.gserviceaccount.com"
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

			objs, err := c.RenderAllClusterResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(9))
		}

		{
			generatedName := "baz-abc1xad"
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "baz",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					Project:        new(string),
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: &generatedName,
				},
			}
			*cluster.Spec.Project = "cilium-ci"

			objs, err := c.RenderAllClusterResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(9))

			name := objs.Items[1].GetName()
			g.Expect(name).To(Equal(generatedName))

			for _, obj := range objs.Items {
				labels := obj.GetLabels()
				g.Expect(labels).To(HaveKeyWithValue("cluster", cluster.Name))
			}
		}

		{
			clusterName := "baz-abc2sax"
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "baz",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					Project:        new(string),
					ConfigTemplate: &templateName,
					MachineType:    &machineType,
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: &clusterName,
				},
			}
			*cluster.Spec.Project = "cilium-ci"

			objs, err := c.RenderAllClusterResources(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(9))

			name := objs.Items[1].GetName()
			g.Expect(name).To(Equal(clusterName))

			for _, obj := range objs.Items {
				labels := obj.GetLabels()
				g.Expect(labels).To(HaveKeyWithValue("cluster", cluster.Name))
			}
		}

		{
			cluster := &v1alpha2.TestClusterGKE{
				Spec: v1alpha2.TestClusterGKESpec{
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

func TestTestInfraWorkloadsResources(t *testing.T) {
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

		_, err = c.RenderTestInfraWorkloadsAsJSON(&v1alpha2.TestClusterGKE{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderTestInfraWorkloadsAsJSON(&v1alpha2.TestClusterGKE{ObjectMeta: metav1.ObjectMeta{Name: "foo"}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderTestInfraWorkloadsAsJSON(&v1alpha2.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			Spec: v1alpha2.TestClusterGKESpec{
				JobSpec: &v1alpha2.TestClusterGKEJobSpec{},
			}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec.runner`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderTestInfraWorkloadsAsJSON(&v1alpha2.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{Name: "foo"},
			Spec: v1alpha2.TestClusterGKESpec{
				JobSpec: &v1alpha2.TestClusterGKEJobSpec{
					Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{},
				},
			}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil/empty jobSpec.runner.image`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		runnerImage := "cilium-ci/cilium-e2e:latest"

		_, err = c.RenderTestInfraWorkloadsAsJSON(&v1alpha2.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"},
			Spec: v1alpha2.TestClusterGKESpec{
				JobSpec: &v1alpha2.TestClusterGKEJobSpec{
					Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
						Image: &runnerImage,
					},
				},
			}})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil status.clusterName`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		defProject := "cilium-ci"
		defRunnerImage := "cilium-ci/cilium-e2e:latest"
		defRunnerInitImage := "quay.io/isovalent/gke-test-cluster-initutil:660e365e201df32d61efd57a112c19d242743ae6"
		defCluster := &v1alpha2.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
			},
			Spec: v1alpha2.TestClusterGKESpec{
				Project:  &defProject,
				Region:   new(string),
				Location: new(string),
				JobSpec: &v1alpha2.TestClusterGKEJobSpec{
					Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
						Image:     &defRunnerImage,
						InitImage: &defRunnerInitImage,
					},
				},
			},
		}
		*defCluster.Spec.Region = "europe-west2"
		*defCluster.Spec.Location = "europe-west2-b"

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		err = c.ApplyDefaultsForTestInfraWorkloads(defCluster)
		g.Expect(err).ToNot(HaveOccurred())

		{
			generatedName := "baz-a0b1c2"
			runnerImage := "cilium-ci/cilium-e2e:80d4133f2b9317a0f08fcff9b2f8d625ea9f7b7a"
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "other",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					JobSpec: &v1alpha2.TestClusterGKEJobSpec{
						Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
							Image:   &runnerImage,
							Command: []string{"app.test", "-test.v"},
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "bar",
							}},
							ConfigMap: &generatedName,
						},
					},
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: &generatedName,
				},
			}

			data, err := c.RenderTestInfraWorkloadsAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			var expected = getInfraWithDashboard(generatedName, `
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
						"component": "test-runner",
						"cluster": "baz-a0b1c2"
					  }
					},
					"spec": {
					  "backoffLimit": 0,
					  "template": {
						"metadata": {
						  "labels": {
							"cluster": "baz-a0b1c2",
							"component": "test-runner"
						  }
						},
						"spec": {
						  "enableServiceLinks": false,
						  "automountServiceAccountToken": false,
						  "volumes": [
							{
							  "name": "credentials",
							  "emptyDir": {}
							},
							{
							  "name": "config-system",
							  "configMap": {
								"optional": true,
							    "name": "baz-a0b1c2-system"
							  }
							},
							{
							  "name": "config-user",
							  "configMap": {
							    "name": "baz-a0b1c2"
							  }
							}
						  ],
						  "initContainers": [
							{
							  "name": "initutil",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "baz-a0b1c2-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "baz-a0b1c2"
								},
								{
								  "name": "FOO",
								  "value": "bar"
								}
							  ],
							  "image": "quay.io/isovalent/gke-test-cluster-initutil:660e365e201df32d61efd57a112c19d242743ae6",
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								},
								{
								  "name": "config-user",
								  "mountPath": "/config/user"
								}
							  ]
							}
						  ],
						  "containers": [
							{
							  "name": "test-runner",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "baz-a0b1c2-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "baz-a0b1c2"
								},
								{
									"name": "FOO",
									"value": "bar"
								}
							  ],
							  "command": [
								"app.test",
								"-test.v"
							  ],
							  "image": "cilium-ci/cilium-e2e:80d4133f2b9317a0f08fcff9b2f8d625ea9f7b7a",
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								},
								{
								  "name": "config-user",
								  "mountPath": "/config/user"
								}
							  ]
							}
						  ],
						  "dnsPolicy": "ClusterFirst",
						  "restartPolicy": "Never",
						  "serviceAccountName": "baz-a0b1c2-admin"
						}
					  }
					}
				  },
				  {
					"kind": "Deployment",
					"apiVersion": "apps/v1",
					"metadata": {
					  "name": "baz-a0b1c2-promview",
					  "namespace": "other",
					  "labels": {
						"component": "promview",
						"cluster": "baz-a0b1c2"
					  }
					},
					"spec": {
					  "selector": {
						"matchLabels": {
						  "component": "promview",
						  "cluster": "baz-a0b1c2"
						}
					  },
					  "template": {
						"metadata": {
						  "labels": {
							"component": "promview",
							"cluster": "baz-a0b1c2"
						  },
						  "annotations": {
							"prometheus.io.scrape": "false"
						  }
						},
						"spec": {
						  "volumes": [
							{
							  "name": "credentials",
							  "emptyDir": {}
							},
							{
							  "name": "config-system",
							  "configMap": {
								"optional": true,
								"name": "baz-a0b1c2-system"
							  }
							},
							{
							  "name": "config-user",
							  "configMap": {
								"name": "baz-a0b1c2"
							  }
							}
						  ],
						  "initContainers": [
							{
							  "name": "initutil",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "baz-a0b1c2-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "baz-a0b1c2"
								},
								{
								  "name": "FOO",
								  "value": "bar"
								}
							  ],
							  "image": "quay.io/isovalent/gke-test-cluster-initutil:660e365e201df32d61efd57a112c19d242743ae6",
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								},
								{
								  "name": "config-user",
								  "mountPath": "/config/user"
								}
							  ]
							}
						  ],
						  "containers": [
							{
							  "name": "promview",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "baz-a0b1c2-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "baz-a0b1c2"
								},
								{
								  "name": "FOO",
								  "value": "bar"
								}
							  ],
							  "resources": {
								"limits": {
								  "cpu": "100m",
								  "memory": "400Mi"
								},
								"requests": {
								  "cpu": "100m",
								  "memory": "400Mi"
								}
							  },
							  "image": "quay.io/isovalent/gke-test-cluster-promview:7695938dcf3a6e4f0e7fb9537091103259aed46e",
							  "command": [
								"/usr/bin/gke-test-cluster-promview"
							  ],
							  "ports": [
								{
								  "name": "http",
								  "containerPort": 8080
								}
							  ],
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								},
								{
								  "name": "config-user",
								  "mountPath": "/config/user"
								}
							  ]
							}
						  ],
						  "terminationGracePeriodSeconds": 10,
						  "serviceAccountName": "baz-a0b1c2-admin",
						  "automountServiceAccountToken": false,
						  "enableServiceLinks": false
						}
					  },
					  "replicas": 2
					}
				  },
				  {
					"kind": "Service",
					"apiVersion": "v1",
					"metadata": {
					  "name": "baz-a0b1c2-promview",
					  "namespace": "other",
					  "labels": {
						"component": "promview",
						"cluster": "baz-a0b1c2"
					  }
					},
					"spec": {
					  "selector": {
						"component": "promview",
						"cluster": "baz-a0b1c2"
					  },
					  "ports": [
						{
						  "name": "promview",
						  "port": 80,
						  "targetPort": 8080
						}
					  ]
					}
				  },
				  {
                  "kind": "ConfigMap",
                  "apiVersion": "v1",
                  "metadata": {
                    "name": "dashboard-baz-a0b1c2-cilium",
                    "namespace": "grafana",
                    "labels": {
                      "grafana_dashboard": "1",
					  "cluster": "baz-a0b1c2",
					  "component": "dashboard"
                    }
                  },
                  "data": {
                    "dashboard-baz-a0b1c2-cilium.json": CILIUM_DASHBOARD_DATA
                  }
				  },
				  {
                  "kind": "ConfigMap",
                  "apiVersion": "v1",
                  "metadata": {
                    "name": "dashboard-baz-a0b1c2-cilium-operator",
                    "namespace": "grafana",
                    "labels": {
                      "grafana_dashboard": "1",
					  "cluster": "baz-a0b1c2",
					  "component": "dashboard"
                    }
                  },
                  "data": {
                    "dashboard-baz-a0b1c2-cilium-operator.json": CILIUM_OPERATOR_DASHBOARD_DATA
                  }
                },
				{
                  "kind": "ConfigMap",
                  "apiVersion": "v1",
                  "metadata": {
                    "name": "dashboard-baz-a0b1c2-hubble",
                    "namespace": "grafana",
                    "labels": {
                      "grafana_dashboard": "1",
					  "cluster": "baz-a0b1c2",
					  "component": "dashboard"
                    }
                  },
                  "data": {
                    "dashboard-baz-a0b1c2-hubble.json": HUBBLE_DASHBOARD_DATA
                  }
				  }
				]
			}`)
			g.Expect(data).To(MatchJSON(expected))

			objs, err := c.RenderTestInfraWorkloads(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(6))
		}

		{
			generatedName := "bar-0a1b2c"
			runnerImage := "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a"
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "bar",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					JobSpec: &v1alpha2.TestClusterGKEJobSpec{
						Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
							Image: &runnerImage,
						},
					},
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: &generatedName,
				},
			}

			data, err := c.RenderTestInfraWorkloadsAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())

			var expected = getInfraWithDashboard(generatedName, `
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
						"component": "test-runner",
						"cluster": "bar-0a1b2c"
					  }
					},
					"spec": {
					  "backoffLimit": 0,
					  "template": {
						"metadata": {
						  "labels": {
							"cluster": "bar-0a1b2c",
							"component": "test-runner"
						  }
						},
						"spec": {
						  "enableServiceLinks": false,
						  "automountServiceAccountToken": false,
						  "volumes": [
							{
							  "name": "credentials",
							  "emptyDir": {}
							},
							{
							  "name": "config-system",
							  "configMap": {
								"optional": true,
							    "name": "bar-0a1b2c-system"
							  }
							}
						  ],
						  "initContainers": [
							{
							  "name": "initutil",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "bar-0a1b2c-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "bar-0a1b2c"
								}
							  ],
							  "image": "quay.io/isovalent/gke-test-cluster-initutil:660e365e201df32d61efd57a112c19d242743ae6",
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								}
							  ]
							}
						  ],
						  "containers": [
							{
							  "name": "test-runner",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "bar-0a1b2c-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "bar-0a1b2c"
								}
							  ],
							  "image": "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a",
							  "command": [],
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								}
							  ]
							}
						  ],
						  "dnsPolicy": "ClusterFirst",
						  "restartPolicy": "Never",
						  "serviceAccountName": "bar-0a1b2c-admin"
						}
					  }
					}
				  },
				  {
					"kind": "Deployment",
					"apiVersion": "apps/v1",
					"metadata": {
					  "name": "bar-0a1b2c-promview",
					  "namespace": "default",
					  "labels": {
						"component": "promview",
						"cluster": "bar-0a1b2c"
					  }
					},
					"spec": {
					  "selector": {
						"matchLabels": {
						  "component": "promview",
						  "cluster": "bar-0a1b2c"
						}
					  },
					  "template": {
						"metadata": {
						  "labels": {
							"component": "promview",
							"cluster": "bar-0a1b2c"
						  },
						  "annotations": {
							"prometheus.io.scrape": "false"
						  }
						},
						"spec": {
						  "volumes": [
							{
							  "name": "credentials",
							  "emptyDir": {}
							},
							{
							  "name": "config-system",
							  "configMap": {
								"optional": true,
								"name": "bar-0a1b2c-system"
							  }
							}
						  ],
						  "initContainers": [
							{
							  "name": "initutil",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "bar-0a1b2c-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "bar-0a1b2c"
								}
							  ],
							  "image": "quay.io/isovalent/gke-test-cluster-initutil:660e365e201df32d61efd57a112c19d242743ae6",
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								}
							  ]
							}
						  ],
						  "containers": [
							{
							  "name": "promview",
							  "env": [
								{
								  "name": "KUBECONFIG",
								  "value": "/credentials/kubeconfig"
								},
								{
								  "name": "SERVICE_ACCOUNT",
								  "value": "bar-0a1b2c-admin@cilium-ci.iam.gserviceaccount.com"
								},
								{
								  "name": "CLUSTER_LOCATION",
								  "value": "europe-west2-b"
								},
								{
								  "name": "CLUSTER_NAME",
								  "value": "bar-0a1b2c"
								}
							  ],
							  "resources": {
								"limits": {
								  "cpu": "100m",
								  "memory": "400Mi"
								},
								"requests": {
								  "cpu": "100m",
								  "memory": "400Mi"
								}
							  },
							  "image": "quay.io/isovalent/gke-test-cluster-promview:7695938dcf3a6e4f0e7fb9537091103259aed46e",
							  "command": [
								"/usr/bin/gke-test-cluster-promview"
							  ],
							  "ports": [
								{
								  "name": "http",
								  "containerPort": 8080
								}
							  ],
							  "volumeMounts": [
								{
								  "name": "credentials",
								  "mountPath": "/credentials"
								},
								{
								  "name": "config-system",
								  "mountPath": "/config/system"
								}
							  ]
							}
						  ],
						  "terminationGracePeriodSeconds": 10,
						  "serviceAccountName": "bar-0a1b2c-admin",
						  "automountServiceAccountToken": false,
						  "enableServiceLinks": false
						}
					  },
					  "replicas": 2
					}
				  },
				  {
					"kind": "Service",
					"apiVersion": "v1",
					"metadata": {
					  "name": "bar-0a1b2c-promview",
					  "namespace": "default",
					  "labels": {
						"component": "promview",
						"cluster": "bar-0a1b2c"
					  }
					},
					"spec": {
					  "selector": {
						"component": "promview",
						"cluster": "bar-0a1b2c"
					  },
					  "ports": [
						{
						  "name": "promview",
						  "port": 80,
						  "targetPort": 8080
						}
					  ]
					}
				  },
				  {
                  "kind": "ConfigMap",
                  "apiVersion": "v1",
                  "metadata": {
                    "name": "dashboard-bar-0a1b2c-cilium",
                    "namespace": "grafana",
                    "labels": {
                      "grafana_dashboard": "1",
					  "cluster": "bar-0a1b2c",
					  "component": "dashboard"
                    }
                  },
                  "data": {
                    "dashboard-bar-0a1b2c-cilium.json": CILIUM_DASHBOARD_DATA
                  }
				  },
				  {
                  "kind": "ConfigMap",
                  "apiVersion": "v1",
                  "metadata": {
                    "name": "dashboard-bar-0a1b2c-cilium-operator",
                    "namespace": "grafana",
                    "labels": {
                      "grafana_dashboard": "1",
					  "cluster": "bar-0a1b2c",
					  "component": "dashboard"
                    }
                  },
                  "data": {
                    "dashboard-bar-0a1b2c-cilium-operator.json": CILIUM_OPERATOR_DASHBOARD_DATA
                  }
                },
				{
                  "kind": "ConfigMap",
                  "apiVersion": "v1",
                  "metadata": {
                    "name": "dashboard-bar-0a1b2c-hubble",
                    "namespace": "grafana",
                    "labels": {
                      "grafana_dashboard": "1",
					  "cluster": "bar-0a1b2c",
					  "component": "dashboard"
                    }
                  },
                  "data": {
                    "dashboard-bar-0a1b2c-hubble.json": HUBBLE_DASHBOARD_DATA
                  }
				  }
				]
			}`)
			g.Expect(data).To(MatchJSON(expected))

			objs, err := c.RenderTestInfraWorkloads(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(6))
		}

		{
			generatedName := "baz-x2a8332"
			runnerImage := "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a"
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name: "baz",
				},
				Spec: v1alpha2.TestClusterGKESpec{
					JobSpec: &v1alpha2.TestClusterGKEJobSpec{
						Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
							Image: &runnerImage,
						}},
				},
				Status: v1alpha2.TestClusterGKEStatus{
					ClusterName: &generatedName,
				},
			}

			objs, err := c.RenderTestInfraWorkloads(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(6))

			g.Expect(objs.Items[0].GetName()).To(Equal("test-runner-baz-x2a8332"))

			for _, obj := range objs.Items {
				labels := obj.GetLabels()
				g.Expect(labels).To(HaveKeyWithValue("cluster", "baz-x2a8332"))
				g.Expect(labels).To(HaveKey("component"))
				g.Expect(obj.GetName()).To(ContainSubstring("baz-x2a8332"))
			}
		}

		{
			runnerImage := "cilium-ci/cilium-e2e:0d725ea9f7ba0f08fcff48133f2b9319b2f8d67a"
			cluster := &v1alpha2.TestClusterGKE{
				Spec: v1alpha2.TestClusterGKESpec{
					JobSpec: &v1alpha2.TestClusterGKEJobSpec{
						Runner: &v1alpha2.TestClusterGKEJobRunnerSpec{
							Image: &runnerImage,
						}},
				},
			}

			_, err := c.RenderTestInfraWorkloadsAsJSON(cluster)
			g.Expect(err).To(HaveOccurred())
			// this is another weird error from CUE, but that's what you get when optional field is unspecified on export...
			g.Expect(err.Error()).ToNot(Equal(`cue: marshal error: template.items.0.metadata.name: field "name" is optional`))
			g.Expect(err.Error()).To(Equal(`unexpected unnamed object`))
		}
	}
}

func TestPromResources(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		{
			cluster := &v1alpha2.TestClusterGKE{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "baz",
					Namespace: "other",
				},
			}

			cluster.Default()

			data, err := c.RenderPromResourcesAsJSON(cluster)
			g.Expect(err).ToNot(HaveOccurred())
			objs := &unstructured.UnstructuredList{}
			g.Expect(objs.UnmarshalJSON(data)).To(Succeed())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(9))
		}
	}
}

func TestToUnstructured(t *testing.T) {
	g := NewGomegaWithT(t)

	runtimeCM := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "name",
			Namespace: "namespace",
			Labels: map[string]string{
				"labelname": "labelvalue",
			},
		},
		BinaryData: map[string][]byte{
			"datakey": []byte("value"),
		},
	}

	unstructuredCM, err := ToUnstructured(&runtimeCM)

	g.Expect(err).ToNot(HaveOccurred())

	runtimeJSON, err := json.Marshal(runtimeCM)
	g.Expect(err).ToNot(HaveOccurred())

	unstructuredJSON, err := json.Marshal(unstructuredCM)
	g.Expect(err).ToNot(HaveOccurred())

	g.Expect(runtimeJSON).To(MatchJSON(unstructuredJSON))
}

func getInfraWithDashboard(clusterName, infra string) string {
	for dashboardDataPlaceholder, dashboardTemplate := range dashboardMap {
		dashboard := strings.Replace(dashboardTemplate, clusterNamePlaceholder, clusterName, -1)
		infra = strings.Replace(infra, dashboardDataPlaceholder, dashboard, -1)
	}

	return infra
}
