// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

const (
	ClusterAccessResourcesTemplateName = "iam"
	TestInfraWorkloadsTemplateName     = "infra"
	PromResourcesTemplateName          = "prom"
)

type Config struct {
	BaseDirectory string

	templates map[string]*template.Generator
}

func (c *Config) Load() error {
	entries, err := ioutil.ReadDir(c.BaseDirectory)
	if err != nil {
		return fmt.Errorf("unable to list avaliable config templates in %q: %w", c.BaseDirectory, err)
	}

	c.templates = map[string]*template.Generator{}

	for _, entry := range entries {
		if entry.IsDir() {
			// both path.Join and filpath.Join break this by striping leading `./`,
			// just like Go, relative package path in must be prefixed with `./`
			// (or `../`)
			fullPath := c.BaseDirectory + "/" + entry.Name()
			template := &template.Generator{
				InputDirectory: fullPath,
			}
			if err := template.CompileAndValidate(); err != nil {
				return fmt.Errorf("unable to load config template from %q: %w", fullPath, err)
			}
			c.templates[entry.Name()] = template
		}
	}

	if len(c.templates) == 0 {
		return fmt.Errorf("no config templates found in %q", c.BaseDirectory)
	}
	return nil
}

func (c *Config) HaveExistingTemplate(name string) bool {
	_, ok := c.templates[name]
	return ok
}

func (c *Config) ExistingTemplates() []string {
	templates := []string{}
	for template := range c.templates {
		templates = append(templates, template)
	}
	return templates
}

func (c *Config) ApplyDefaults(templateName string, defaults *v1alpha1.TestClusterGKE) error {
	if !c.HaveExistingTemplate(templateName) {
		return fmt.Errorf("no such template: %q", templateName)
	}
	template := c.templates[templateName]
	template, err := template.WithDefaults(defaults.WithoutTypeMeta())
	if err != nil {
		return err
	}
	c.templates[templateName] = template
	return nil
}
func (c *Config) ApplyDefaultsForClusterAccessResources(defaults *v1alpha1.TestClusterGKE) error {
	return c.ApplyDefaults(ClusterAccessResourcesTemplateName, defaults)
}

func (c *Config) ApplyDefaultsForTestInfraWorkloads(defaults *v1alpha1.TestClusterGKE) error {
	return c.ApplyDefaults(TestInfraWorkloadsTemplateName, defaults)
}

func (c *Config) renderTemplateAsJSON(cluster *v1alpha1.TestClusterGKE, templateName string) ([]byte, error) {
	if cluster == nil {
		return nil, fmt.Errorf("unexpected nil object")
	}
	if cluster.Name == "" {
		return nil, fmt.Errorf("unexpected unnamed object")
	}
	switch templateName {
	case ClusterAccessResourcesTemplateName:
	case TestInfraWorkloadsTemplateName:
		if cluster.Spec.JobSpec == nil {
			return nil, fmt.Errorf("unexpected nil jobSpec")
		}
		if cluster.Spec.JobSpec.Runner == nil {
			return nil, fmt.Errorf("unexpected nil jobSpec.runner")
		}
		if cluster.Spec.JobSpec.Runner.Image == nil || *cluster.Spec.JobSpec.Runner.Image == "" {
			return nil, fmt.Errorf("unexpected nil/empty jobSpec.runner.image")
		}
		if cluster.Status.ClusterName == nil {
			return nil, fmt.Errorf("unexpected nil status.clusterName")
		}
		cluster = cluster.DeepCopy()
		cluster.Name = *cluster.Status.ClusterName
	case "":
		if cluster.Spec.ConfigTemplate == nil || *cluster.Spec.ConfigTemplate == "" {
			return nil, fmt.Errorf("unexpected nil/empty configTemplate")
		}
		templateName = *cluster.Spec.ConfigTemplate
		if templateName == TestInfraWorkloadsTemplateName || templateName == ClusterAccessResourcesTemplateName {
			return nil, fmt.Errorf("cannot create cluster directly with configTemplate=%q", templateName)
		}
	}

	if !c.HaveExistingTemplate(templateName) {
		return nil, fmt.Errorf("no such template: %q", templateName)
	}
	template := c.templates[templateName]
	template, err := template.WithResource(cluster.WithoutTypeMeta())
	if err != nil {
		return nil, err
	}
	return template.RenderJSON()
}

func (c *Config) RenderClusterCoreResourcesAsJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	return c.renderTemplateAsJSON(cluster, "")
}

func (c *Config) RenderClusterAccessResourcesAsJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	return c.renderTemplateAsJSON(cluster, ClusterAccessResourcesTemplateName)
}

func (c *Config) RenderTestInfraWorkloadsAsJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	return c.renderTemplateAsJSON(cluster, TestInfraWorkloadsTemplateName)
}

func (c *Config) RenderPromResourcesAsJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	return c.renderTemplateAsJSON(cluster, PromResourcesTemplateName)
}

func (c *Config) RenderAllClusterResources(cluster *v1alpha1.TestClusterGKE, generateName bool) (*unstructured.UnstructuredList, error) {
	allResources := &unstructured.UnstructuredList{}
	coreResources := &unstructured.UnstructuredList{}
	accessResources := &unstructured.UnstructuredList{}

	if generateName {
		generatedName := cluster.Name + "-" + utilrand.String(5)
		if cluster.Status.ClusterName == nil {
			// store generated name in status
			cluster.Status.ClusterName = &generatedName
		} else {
			// if name is already stored in status, use that instead
			generatedName = *cluster.Status.ClusterName
		}
		// make a copy and use new name as input to generator
		cluster = cluster.DeepCopy()
		cluster.Name = generatedName
	}

	coreResourcesData, err := c.RenderClusterCoreResourcesAsJSON(cluster)
	if err != nil {
		return nil, err
	}

	if err := coreResources.UnmarshalJSON(coreResourcesData); err != nil {
		return nil, err
	}

	accessResourcesData, err := c.RenderClusterAccessResourcesAsJSON(cluster)
	if err != nil {
		return nil, err
	}

	if err := accessResources.UnmarshalJSON(accessResourcesData); err != nil {
		return nil, err
	}

	promResourcesData, err := c.RenderPromResourcesAsJSON(cluster)
	if err != nil {
		return nil, err
	}

	systemConfigMap, err := toUnstructured(&corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-system",
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
			},
		},
		BinaryData: map[string][]byte{
			"init-manifest": promResourcesData,
		},
	})
	if err != nil {
		return nil, err
	}

	allResources.Items = append(allResources.Items, coreResources.Items...)
	allResources.Items = append(allResources.Items, accessResources.Items...)
	allResources.Items = append(allResources.Items, *systemConfigMap)

	return allResources, nil
}

// toUnstructured convers runtime.Object so that it can be appended to UnstructuredList
// with all other resouces
func toUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err := u.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return u, nil
}

func (c *Config) RenderTestInfraWorkloads(cluster *v1alpha1.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	jobRunnerResources := &unstructured.UnstructuredList{}

	jobRunnerResourcesData, err := c.RenderTestInfraWorkloadsAsJSON(cluster)
	if err != nil {
		return nil, err
	}

	if err := jobRunnerResources.UnmarshalJSON(jobRunnerResourcesData); err != nil {
		return nil, err
	}

	return jobRunnerResources, nil
}
