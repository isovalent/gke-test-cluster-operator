package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/template"
)

type Config struct {
	BaseDir string

	templates map[string]*template.Generator
}

func (c *Config) Load() error {
	entries, err := ioutil.ReadDir(c.BaseDir)
	if err != nil {
		return fmt.Errorf("unable to list avaliable config templates in %q: %w", c.BaseDir, err)
	}

	c.templates = map[string]*template.Generator{}

	for _, entry := range entries {
		if entry.IsDir() {
			fullPath := filepath.Join(c.BaseDir, entry.Name())
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
		return fmt.Errorf("no config templates found in %q", c.BaseDir)
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

func (c *Config) RenderJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	if cluster == nil || cluster.Spec.ConfigTemplate == nil {
		return nil, fmt.Errorf("invalid test cluster object")
	}
	templateName := *cluster.Spec.ConfigTemplate
	if !c.HaveExistingTemplate(templateName) {
		return nil, fmt.Errorf("no such template: %q", templateName)
	}

	return c.templates[templateName].ToJSON(cluster)
}
