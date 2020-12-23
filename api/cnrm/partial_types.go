package cnrm

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	clustersv1alpha2 "github.com/isovalent/gke-test-cluster-operator/api/v1alpha2"
)

type PartialContainerCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec struct {
		Location   string `json:"location"`
		MasterAuth struct {
			ClusterCACertificate string `json:"clusterCaCertificate"`
			ClientKey            string `json:"clientKey"`
			ClientCertificate    string `json:"clientCertificate"`
		} `json:"masterAuth"`
	} `json:"spec,omitempty"`

	Status struct {
		Endpoint string `json:"endpoint"`
	} `json:"status,omitempty"`
}

func ParsePartialContainerCluster(obj *unstructured.Unstructured) (*PartialContainerCluster, error) {
	if kind := obj.GetKind(); kind != "ContainerCluster" {
		return nil, fmt.Errorf("given object is %q, not \"ContainerCluster\"", kind)
	}
	data, err := json.Marshal(obj.Object)
	if err != nil {
		return nil, err
	}

	pcc := &PartialContainerCluster{}
	if err := json.Unmarshal(data, pcc); err != nil {
		return nil, err
	}
	return pcc, nil
}

type PartialStatus struct {
	Conditions clustersv1alpha2.CommonConditions `json:"conditions,omitempty"`
}

func ParsePartialStatus(obj *unstructured.Unstructured) (*PartialStatus, error) {
	statusObj, ok := obj.Object["status"]
	if !ok {
		// ignore objects without status,
		// presumably it just hasn't been populated yet
		return nil, nil
	}

	data, err := json.Marshal(statusObj)
	if err != nil {
		return nil, err
	}

	status := &PartialStatus{}
	if err := json.Unmarshal(data, status); err != nil {
		return nil, err
	}
	return status, nil
}

func (c *PartialStatus) HasReadyCondition() bool {
	return c.Conditions.HaveReadyCondition()
}
