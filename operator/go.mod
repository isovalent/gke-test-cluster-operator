module github.com/isovalent/gke-test-cluster-management/operator

go 1.14

require (
	cuelang.org/go v0.2.1
	github.com/errordeveloper/imagine v0.0.0-20200720095627-da3e30ba15e2
	github.com/errordeveloper/kue v0.2.3
	github.com/go-logr/logr v0.1.0
	github.com/onsi/gomega v1.9.0
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/controller-tools v0.3.0
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
