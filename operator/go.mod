module github.com/isovalent/gke-test-cluster-management/operator

go 1.14

require (
	cuelang.org/go v0.2.1
	github.com/errordeveloper/kue v0.1.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/gomega v1.8.1
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.3.0
)
