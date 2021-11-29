package v1alpha1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/isovalent/gke-test-cluster-operator/api/v1alpha2"
)

// ConvertTo converts this TestClusterGKE to the Hub version.
func (src *TestClusterGKE) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.TestClusterGKE)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Project = src.Spec.Project
	dst.Spec.ConfigTemplate = src.Spec.ConfigTemplate
	dst.Spec.Location = src.Spec.Location
	dst.Spec.Region = src.Spec.Region
	dst.Spec.MultiZone = new(bool)
	*dst.Spec.MultiZone = false
	dst.Spec.KubernetesVersion = src.Spec.KubernetesVersion

	dst.Spec.JobSpec.Runner.Image = src.Spec.JobSpec.Runner.Image
	dst.Spec.JobSpec.Runner.Command = src.Spec.JobSpec.Runner.Command
	dst.Spec.JobSpec.Runner.InitImage = src.Spec.JobSpec.Runner.InitImage
	dst.Spec.JobSpec.Runner.Env = src.Spec.JobSpec.Runner.Env
	dst.Spec.JobSpec.Runner.ConfigMap = src.Spec.JobSpec.Runner.ConfigMap

	dst.Spec.JobSpec.ImagesToTest = src.Spec.JobSpec.ImagesToTest

	dst.Spec.MachineType = src.Spec.MachineType
	dst.Spec.Nodes = src.Spec.Nodes

	dst.Status.ClusterName = src.Status.ClusterName

	if len(src.Status.Conditions) > 0 {
		dstCondtions := v1alpha2.CommonConditions{}
		for _, srcCondition := range src.Status.Conditions {
			dstCondition := v1alpha2.CommonCondition{
				Type:               srcCondition.Type,
				Status:             srcCondition.Status,
				LastTransitionTime: srcCondition.LastTransitionTime,
				Reason:             srcCondition.Reason,
				Message:            srcCondition.Message,
			}
			dstCondtions = append(dstCondtions, dstCondition)
		}
		dst.Status.Dependencies = map[string]v1alpha2.CommonConditions{
			fmt.Sprintf("ContainerCluster:%s/%s", src.Namespace, src.Name): dstCondtions,
		}
	}

	readinessStatus := "False"
	readinessReason := "DependenciesNotReady"
	readinessMessage := "Some dependencies are not ready yet"

	if dst.Status.AllDependeciesReady() {
		readinessStatus = "True"
		readinessReason = "AllDependenciesReady"
		readinessMessage = fmt.Sprintf("All %d dependencies are ready", len(dst.Status.Dependencies))
	}

	dst.Status.Conditions = v1alpha2.CommonConditions{{
		Type:               "Ready",
		Status:             readinessStatus,
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             readinessReason,
		Message:            readinessMessage,
	}}

	return nil
}

// ConvertFrom converts from Hub version to this TestClusterGKE.
func (dst *TestClusterGKE) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.TestClusterGKE)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Project = src.Spec.Project
	dst.Spec.ConfigTemplate = src.Spec.ConfigTemplate
	dst.Spec.Location = src.Spec.Location
	dst.Spec.Region = src.Spec.Region

	dst.Spec.KubernetesVersion = src.Spec.KubernetesVersion

	dst.Spec.JobSpec.Runner.Image = src.Spec.JobSpec.Runner.Image
	dst.Spec.JobSpec.Runner.Command = src.Spec.JobSpec.Runner.Command
	dst.Spec.JobSpec.Runner.InitImage = src.Spec.JobSpec.Runner.InitImage
	dst.Spec.JobSpec.Runner.Env = src.Spec.JobSpec.Runner.Env
	dst.Spec.JobSpec.Runner.ConfigMap = src.Spec.JobSpec.Runner.ConfigMap

	dst.Spec.JobSpec.ImagesToTest = src.Spec.JobSpec.ImagesToTest

	dst.Spec.MachineType = src.Spec.MachineType
	dst.Spec.Nodes = src.Spec.Nodes

	dst.Status.ClusterName = src.Status.ClusterName

	if srcConditions, ok := src.Status.Dependencies[fmt.Sprintf("ContainerCluster:%s/%s", src.Namespace, src.Name)]; ok {
		for _, srcCondition := range srcConditions {
			dstCondition := v1alpha2.CommonCondition{
				Type:               srcCondition.Type,
				Status:             srcCondition.Status,
				LastTransitionTime: srcCondition.LastTransitionTime,
				Reason:             srcCondition.Reason,
				Message:            srcCondition.Message,
			}
			dst.Status.Conditions = append(dst.Status.Conditions, dstCondition)
		}
	}

	return nil
}
