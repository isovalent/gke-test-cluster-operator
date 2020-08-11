module github.com/isovalent/gke-test-cluster-management/operator

go 1.14

require (
	github.com/errordeveloper/imagine v0.0.0-20200811092946-97d806f9167f
	github.com/errordeveloper/kue v0.2.5
	sigs.k8s.io/controller-tools v0.3.0
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
