package config_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/isovalent/gke-test-cluster-management/operator/pkg/config"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		c := &Config{
			BaseDir: "./nonexistent",
		}

		err := c.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unable to list avaliable config templates in "./nonexistent": open ./nonexistent: no such file or directory`))
	}

	{
		c := &Config{
			BaseDir: "./",
		}

		err := c.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`no config templates found in "./"`))
	}

	{
		c := &Config{
			BaseDir: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).To(Not(HaveOccurred()))

		g.Expect(c.HaveExistingTemplate("basic")).To(BeTrue())
	}

	{
		c := &Config{
			BaseDir: "../../config/templates",
		}

		err := c.Load()
		g.Expect(err).To(Not(HaveOccurred()))

		_, err = c.RenderJSON(&v1alpha1.TestClusterGKE{})
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`invalid test cluster object`))
	}

	{
		c := &Config{
			BaseDir: "../../config/templates",
		}
		templateName := "basic"
		cluster := &v1alpha1.TestClusterGKE{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: v1alpha1.TestClusterGKESpec{
				ConfigTemplate: &templateName,
			},
		}
		err := c.Load()
		g.Expect(err).To(Not(HaveOccurred()))

		g.Expect(c.HaveExistingTemplate("basic")).To(BeTrue())
		_, err = c.RenderJSON(cluster)
		g.Expect(err).To(Not(HaveOccurred()))

		//	g.Expect(err.Error()).To(Equal(`invalid test cluster object`))

	}
}
