// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"io/ioutil"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
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

func (c *Config) renderTemplateAsJSON(cluster *v1alpha1.TestClusterGKE, templateName string) ([]byte, error) {
	if cluster == nil {
		return nil, fmt.Errorf("unexpected nil object")
	}
	if templateName == "" {
		if cluster.Spec.ConfigTemplate == nil || *cluster.Spec.ConfigTemplate == "" {
			return nil, fmt.Errorf("unexpected nil/empty configTemplate")
		}
		templateName = *cluster.Spec.ConfigTemplate
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

func (c *Config) RenderCoreResourcesAsJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	return c.renderTemplateAsJSON(cluster, "")
}

func (c *Config) RenderAccessResourcesAsJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	return c.renderTemplateAsJSON(cluster, "iam")
}

func (c *Config) RenderAllClusterResources(cluster *v1alpha1.TestClusterGKE, generateName bool) (*unstructured.UnstructuredList, error) {
	allResources := &unstructured.UnstructuredList{}
	coreResources := &unstructured.UnstructuredList{}
	accessResources := &unstructured.UnstructuredList{}

	if generateName {
		clusterName := cluster.Name + "-" + utilrand.String(5)
		if cluster.Status.ClusterName == nil {
			// store generated name in status
			cluster.Status.ClusterName = &clusterName
		} else {
			// if name is already stored in status, use that instead
			clusterName = *cluster.Status.ClusterName
		}
		// make a copy and use new name as input to generator
		cluster = cluster.DeepCopy()
		cluster.Name = clusterName
	}

	coreResourcesData, err := c.RenderCoreResourcesAsJSON(cluster)
	if err != nil {
		return nil, err
	}

	if err := coreResources.UnmarshalJSON(coreResourcesData); err != nil {
		return nil, err
	}

	accessResourcesData, err := c.RenderAccessResourcesAsJSON(cluster)
	if err != nil {
		return nil, err
	}

	if err := accessResources.UnmarshalJSON(accessResourcesData); err != nil {
		return nil, err
	}

	allResources.Items = append(allResources.Items, coreResources.Items...)
	allResources.Items = append(allResources.Items, accessResources.Items...)

	return allResources, nil
}
