# Test Cluster Operator for GKE

This operator provides an API-driven cluster provisioning for integration and performance testing of software that integrates very
deeply with Kubernetes APIs, cloud provider APIs and other elements of Kubernetes cluster infrastructure.

> This project was developed at [Isovalent][] for very specific needs and later open-sourced as-is. At present, a few code changes
> will be needed to remove assumption about how it's deployed, and deployment docs are also needed.
> If you wish to use this operator and contribute to the project, please join the #testing channel on [Cilium & eBPF Slack][slack].

## Motivation

> NB: Current implementation focuses on GKE (and GCP), most of the ideas described below would apply to any managed Kubernetes
> providers, and some are even more general.


It is relatively easy to test a workload on Kubernetes, whether it's just an application comprised of multiple components,
or a basic operator that glues together a handful of Kubernetes APIs.

Cluster-sharing is one option, that is when application under test is deployed into one or more namespaces in a large shared cluster.
Another option is to run setup small test clusters using something like [kind][] or [k3s][]/[k3d][] within a CI environment.

If the application under test depends on non-namespaced resources, cluster-sharing is still possible with [VirtualCluster][].
That way instances of the application under test can be isolated from one another, but only if Kubernetes API boundaries are fully
respected. It implies that a large underlying cluster will still be used, but virtually divided into small "pretend-clusters".
However, that will only work if the application doesn't make assumptions about cloud provider APIs and doesn't attempt non-trivial
modes of access to the underlying host OS or the network infrastructure.

When the application under test interacts with multiple Kubernetes APIs, and presumes cluster-wide access, and even attempts to
interact with the underlying host OS or the network infrastructure, any kind of cluster-sharing setup may prove challenging to use.
It may also be deemed unrepresentative of the clusters that end-users run. Additionally, testing integrations with cloud provider
APIs may have other implications. Applications that enable core functionality of Kubernetes cluster often fall into this category,
e.g. CNI implementations, storage controllers, service meshes, ingress controller, etc. Cluster-sharing becomes just not viable for
some of these use-cases. And something like kind or k3d is of very limited use.

All of the above applies to testing of applications by developers that are directly responsible for the application itself. End-users
may need to test off-the-shelf application they are deploying also. Quite commonly in large organisation, an operations team will
assemble a bundle of Kubernetes addons that defines a platform that their organisation relies on. The operations team may not be able
to make direct changes to source code of some of the components in order to improved testability for cluster-sharing, or they just
won't have the confidence in testing those components in a shared cluster. Even if one version was easily testable in a shared
cluster, it may change in the future. While testing on kind or k3s remains an option, it may be undesirable due to cloud provider
integration that needs to be tested also, and could be just unrepresentative of the deployment target. Therefore, the operations team
may have strong preference to test in a cluster that is provisioned in exactly the same way as the deployment target and has mostly
identical or comparable properties.

These are just some of the use-cases that illustrate a need for getting a dedicated cluster for running integrations or performance
tests, one that matches deployment target as closely as possible.

What does it take to obtain a cluster in GKE? Technically, it's possible to simply write a script that calls `gcloud` commands, or
relies on something like Terraform or use API client to provision a cluster. This approach inevitably adds a lot of complexity to
the CI job by inheriting all the different failure modes there are to the provisioning and destruction processes, it needs to carry
any additional infrastructure configuration (e.g. metric & log gathering), widens access scopes etc. Aside from all of the steps
that take time and are hard to optimise, it is possible to have a pool of pre-built clusters, yet make the script even more complex.
It is hard to maintain complex scripts of this kind long-term, as by nature scripts don't offers a clear contract (especially the
shell scripts). The lack of contract makes it too easy for anyone to tweak a shell script for an ad-hoc use-case without adding any
tests. Over time, script evolution is hard to unwind, especially in a context where many developers contribute to the project.
In contrast, an API offers many advantages - it's a contract, and the implementation can be optimised more easily.

### Architectural goals of this project

- Test Cluster API
  - enables developer and CI jobs to request clusters for running tests in a consistent and well-defined manner
  - provider abstraction that will enable future optimisations, e.g. pooling of pre-built cluster
- Asynchronous workflow
  - avoid heavy-lifting logic in CI jobs that doesn't directly relate to building binaries or executing tests
  - avoid polling for status
    - once cluster is ready, launch a test runner job inside the management cluster, report the results back to GitHub
- Enable support multiple test cluster templates
  - do not assume there is only one type of test cluster configuration that's going to be used for all purposes
  - allow for pooling pre-built clusters base on commonly used templates
- Include a common set of components in each test cluster
  - Prometheus
  - Log exporter for CI

### You may ask...

**How is this different from something like [Cluster API][]?**

The Test Cluster API is aimed to be much more high-level and shouldn't need to expose as many parameters as Cluster API does,
in fact, it can be implemented on top of Cluster API. The initial implementation targets GKE, and relies on [Config Connector][],
which is similar to Cluster API in spirit.

**What about other providers?**

This is something that authors of this project are planning on exploring, albeit it may not be done as part of the same project
to begin with. One of the ideas is to create a generic provider based on either Terraform or Cluster API, possibly both.

## How it works

There is a management cluster that runs on GKE, it has Config Connector, Cert Manager and Contour along with the GKE Test Cluster
operator ("the operator" from here onwards).

User creates a CR similar to this:

```YAML
apiVersion: clusters.ci.cilium.io/v1alpha2
kind: TestClusterGKE
metadata:
  name: testme-1
  namespace: test-clusters

spec:
  configTemplate: basic
  jobSpec:
    runner:
      image: cilium/cilium-test:8cfdbfe
      command:
      - /usr/local/bin/test-gke.sh
  machineType: n1-standard-4
  nodes: 2
  project: cilium-ci
  location: europe-west2-b
  region: europe-west2
```

The operator renders various objects for Config Connector and other APIs as defined in `basic` template, it substitutes the given
parameters, i.e. `machineType`, `nodes` etc, and then it creates all of these objects and monitors the cluster until it's ready.

Once the test cluster is ready, it deploys the job using the given image and command, and ensures the job is authenticated to run
against the test cluster. The job runs inside management cluster. The test cluster is deleted upon job completion.

The template is defined using [CUE][] and can define any Kubernetes objects, such as Config Connector objects that define additional
GCP resources or some other objects in the management cluster to support test execution. That being said, the implementation
currently expects to find exactly one `ContainerCluster` as part of the template and it's not fully generalised.

As part of test cluster provisioning, Prometheus is deployed in the test cluster and metrics are federated to the Prometheus server
in the management cluster, so all metrics from all test runs can be accessed centrally. In the future other components can be added
as needed.

## Example 2

Here is what a `TestClusterGKE` object may look like with additional fields and status.

```YAML
apiVersion: clusters.ci.cilium.io/v1alpha2
kind: TestClusterGKE
metadata:
  name: test-c6v87
  namespace: test-clusters

spec:
  configTemplate: basic
  jobSpec:
    runner:
      command:
      - /usr/local/bin/run_in_test_cluster.sh
      - --prom-name=prom
      - --prom-ns=prom
      - --duration=30m
      configMap: test-c6v87-user
      image: cilium/hubble-perf-test:8cfdbfe
      initImage: quay.io/isovalent/gke-test-cluster-initutil:854733411778d633350adfa1ae66bf11ba658a3f
  location: europe-west2-b
  machineType: n1-standard-4
  nodes: 2
  project: cilium-ci
  region: europe-west2

status:
  clusterName: test-c6v87-fn86p
  conditions:
  - lastTransitionTime: "2020-11-17T09:29:33Z"
    message: All 2 dependencies are ready
    reason: AllDependenciesReady
    status: "True"
    type: Ready
  dependencyConditions:
    ContainerCluster:test-clusters/test-c6v87-fn86p:
    - lastTransitionTime: "2020-11-17T09:29:22Z"
      message: The resource is up to date
      reason: UpToDate
      status: "True"
      type: Ready
    ContainerNodePool:test-clusters/test-c6v87-fn86p:
    - lastTransitionTime: "2020-11-17T09:29:33Z"
      message: The resource is up to date
      reason: UpToDate
      status: "True"
      type: Ready
```

## Using Test Cluster Requester

There is a simple Go program that serves as a client to the GKE Test Cluster Operator.

It can be use by CI jobs as well as developers.

### Developer Usage

To run this program outside CI, you must ensure that Google Cloud SDK Application credentials are setup correctly, to do so, run:
```
gcloud auth application-default login
```

Run:
```
go run ./requester --namespace=test-clusters-dev --description="<your name and purpose of this cluster>"
```

### CI Usage

This program supports the traditional `GOOGLE_APPLICATION_CREDENTIALS` environment variable, but also for convenience it has
`GCP_SERVICE_ACCOUNT_KEY`  that is expected to contain a base64-encoded JSON service account key (i.e. no need to have the data
written to a file).

For GitHub Actions, it's recommended to use the official image:
```
      - name: Request GKE test cluster
        uses: docker://quay.io/isovalent/gke-test-cluster-requester:ad06d7c2151d012901fc2ddc92406044f2ffba2d
        env:
          GCP_SERVICE_ACCOUNT_KEY: ${{ secrets.GCP_SERVICE_ACCOUNT_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          args: --namespace=... --image=...
```

[VirtualCluster]: https://github.com/kubernetes-sigs/multi-tenancy/tree/master/incubator/virtualcluster
[kind]: https://kind.sigs.k8s.io/
[k3s]: https://k3s.io/
[k3d]: https://k3d.io/
[Cluster API]: https://github.com/kubernetes-sigs/cluster-api
[Config Connector]: https://cloud.google.com/config-connector/docs/overview
[CUE]: https://cuelang.org/

[Isovalent]: https://www.isovalent.com
[slack]: http://slack.cilium.io/
