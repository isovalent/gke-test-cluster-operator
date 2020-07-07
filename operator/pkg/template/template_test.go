package template_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/isovalent/gke-test-cluster-management/operator/pkg/template"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/template/testtypes"
)

func TestGenerator(t *testing.T) {
	g := NewGomegaWithT(t)

	gen := &Generator{
		InputDirectory: "./testassets",
	}

	err := gen.CompileAndValidate()

	g.Expect(err).To(Not(HaveOccurred()))

	{
		cluster := testtypes.Cluster{}
		cluster.Metadata.Name = "foo1"
		cluster.Metadata.Namespace = "default"
		cluster.Spec.Location = "us-central1-a"

		js, err := gen.ToJSON(cluster)
		g.Expect(err).To(Not(HaveOccurred()))

		const expected = `
		{
			"kind": "List",
			"apiVersion": "v1",
			"items": [
			{
				"metadata": {
				"name": "foo1",
				"namespace": "default",
				"labels": {
					"cluster": "foo1"
				},
				"annotations": {
					"cnrm.cloud.google.com/remove-default-node-pool": "false"
				}
				},
				"spec": {
				"location": "us-central1-a",
				"networkRef": {
					"name": "foo1"
				},
				"subnetworkRef": {
					"name": "foo1"
				},
				"initialNodeCount": 1,
				"loggingService": "logging.googleapis.com/kubernetes",
				"monitoringService": "monitoring.googleapis.com/kubernetes",
				"masterAuth": {
					"clientCertificateConfig": {
					"issueClientCertificate": false
					}
				}
				},
				"kind": "ContainerCluster",
				"apiVersion": "container.cnrm.cloud.google.com/v1beta1"
			},
			{
				"metadata": {
				"name": "foo1",
				"namespace": "default",
				"labels": {
					"cluster": "foo1"
				}
				},
				"spec": {
				"routingMode": "REGIONAL",
				"autoCreateSubnetworks": false,
				"deleteDefaultRoutesOnCreate": false
				},
				"kind": "ComputeNetwork",
				"apiVersion": "compute.cnrm.cloud.google.com/v1beta1"
			},
			{
				"metadata": {
				"name": "foo1",
				"namespace": "default",
				"labels": {
					"cluster": "foo1"
				}
				},
				"spec": {
				"networkRef": {
					"name": "foo1"
				},
				"ipCidrRange": "10.128.0.0/20",
				"region": "us-central1"
				},
				"kind": "ComputeSubnetwork",
				"apiVersion": "compute.cnrm.cloud.google.com/v1beta1"
			}
			]
		}
		`
		g.Expect(js).To(MatchJSON(expected))
	}

	{
		cluster := map[string]string{
			"foo": "bar",
		}

		_, err := gen.ToJSON(cluster)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`cue: marshal error: template.items.0.metadata.name: field "foo" not allowed in closed struct`))
	}

	{
		cluster := map[string]string{}

		_, err := gen.ToJSON(cluster)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`cue: marshal error: template.items.0.metadata.name: incomplete value 'resource.metadata.name' in interpolation`))
	}

	{
		_, err := gen.ToJSON(0)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`cue: marshal error: template.items.0.metadata.name: conflicting values Cluster and 0 (mismatched types struct and int)`))
	}

	{
		gen := &Generator{
			InputDirectory: ".",
		}

		err := gen.CompileAndValidate()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`no CUE files in .`))
	}

	{
		gen := &Generator{
			InputDirectory: "./nonexistent",
		}

		err := gen.CompileAndValidate()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`cannot find package "./nonexistent"`))
	}
}
