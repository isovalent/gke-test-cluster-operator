package cnrm

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
