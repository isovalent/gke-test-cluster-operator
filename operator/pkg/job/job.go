// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"fmt"

	"github.com/isovalent/gke-test-cluster-management/operator/api/v1alpha1"
	"github.com/isovalent/gke-test-cluster-management/operator/pkg/template"
)

type Config struct {
	BaseDirectory string

	jobTemplate *template.Generator
}

func (c *Config) Load() error {
	jobTemplate := &template.Generator{
		InputDirectory: c.BaseDirectory + "/job",
	}

	if err := jobTemplate.CompileAndValidate(); err != nil {
		return fmt.Errorf("unable to load job template from %q: %w", c.BaseDirectory, err)
	}

	c.jobTemplate = jobTemplate

	return nil
}

func (c *Config) ApplyDefaults(defaults *v1alpha1.TestClusterGKE) error {
	jobTemplate, err := c.jobTemplate.WithDefaults(defaults.WithoutTypeMeta())
	if err != nil {
		return err
	}
	c.jobTemplate = jobTemplate
	return nil
}

func (c *Config) RenderJSON(clusterRequest *v1alpha1.TestClusterGKE, actualName string) ([]byte, error) {
	if clusterRequest == nil {
		return nil, fmt.Errorf("unexpected nil object")
	}
	if clusterRequest.Spec.JobSpec == nil {
		return nil, fmt.Errorf("unexpected nil jobSpec")
	}
	if clusterRequest.Spec.JobSpec.RunnerImage == nil || *clusterRequest.Spec.JobSpec.RunnerImage == "" {
		return nil, fmt.Errorf("unexpected nil/empty runnerImage")
	}
	cluster := clusterRequest.DeepCopy()
	cluster.SetName(actualName)
	jobTemplate, err := c.jobTemplate.WithResource(cluster.WithoutTypeMeta())
	if err != nil {
		return nil, err
	}
	return jobTemplate.RenderJSON()
}
