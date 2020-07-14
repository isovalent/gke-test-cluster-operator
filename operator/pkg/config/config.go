package config

import (
	"fmt"
	"io/ioutil"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func (c *Config) RenderJSON(cluster *v1alpha1.TestClusterGKE) ([]byte, error) {
	if cluster == nil {
		return nil, fmt.Errorf("unexpected nil object")
	}
	if cluster.Spec.ConfigTemplate == nil {
		return nil, fmt.Errorf("unexpected nil config template")
	}
	templateName := *cluster.Spec.ConfigTemplate
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

func (c *Config) RenderObjects(cluster *v1alpha1.TestClusterGKE) (*unstructured.UnstructuredList, error) {
	objs := &unstructured.UnstructuredList{}

	data, err := c.RenderJSON(cluster)
	if err != nil {
		return nil, err
	}

	if err := objs.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return objs, nil
}
