// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package job_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/isovalent/gke-test-cluster-management/operator/pkg/job"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestJob(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		c := &Config{
			BaseDirectory: "./nonexistent",
		}

		err := c.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unable to load job template from "./nonexistent": cannot find package "./nonexistent/job"`))
	}

	{
		c := &Config{
			BaseDirectory: "./",
		}

		err := c.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unable to load job template from "./": cannot find package "./job"`))
	}

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

		_, err = c.RenderJSON(&v1alpha1.TestClusterGKE{}, "")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderJSON(&v1alpha1.TestClusterGKE{}, "")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil jobSpec`))
	}

	{
		c := &Config{
			BaseDirectory: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).ToNot(HaveOccurred())

		_, err = c.RenderJSON(&v1alpha1.TestClusterGKE{
			Spec: v1alpha1.TestClusterGKESpec{
				JobSpec: &v1alpha1.TestClusterGKEJobSpec{},
			}}, "")
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unexpected nil/empty runnerImage`))
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

		err = c.ApplyDefaults(defCluster)
		g.Expect(err).ToNot(HaveOccurred())

		{
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
			}

			data, err := c.RenderJSON(cluster, "baz-a0b1c2")
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

			objs, err := c.RenderObjects(cluster, "baz-a0b1c2")
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(1))
		}

		{
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
			}

			data, err := c.RenderJSON(cluster, "bar-0a1b2c")
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

			objs, err := c.RenderObjects(cluster, "bar-0a1b2c")
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(objs).ToNot(BeNil())
			g.Expect(objs.Items).To(HaveLen(1))
		}

		{
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
			}

			objs, err := c.RenderObjects(cluster, "baz-x2a8332")
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

			_, err := c.RenderJSON(cluster, "")
			g.Expect(err).To(HaveOccurred())
			// this is another weird error from CUE, but that's what you get when optional field is unspecified on export...
			g.Expect(err.Error()).To(Equal(`cue: marshal error: template.items.0.metadata.name: field "name" is optional`))
		}
	}
}
