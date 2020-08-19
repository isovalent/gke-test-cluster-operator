package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/isovalent/gke-test-cluster-management/operator/pkg/requester"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

func main() {

	image := flag.String("image", "", "image to test - this should be ether full image name, or path to a file containing the name")
	namespace := flag.String("namespace", "", "namespace to use")

	namePrefix := flag.String("name-prefix", "test-", "name prefix for the test cluster")

	project := flag.String("project", requester.DefaultProject, "GCP project")
	managementCluster := flag.String("management-cluster", requester.DefaultManagementCluster, "name of the management cluster")

	initManifest := flag.String("init-manifest", "", "path to manifest to use to initialise the cluster")

	flag.Parse()

	if image == nil || *image == "" {
		log.Fatal("--image must be set")
	}

	if namespace == nil || *namespace == "" {
		log.Fatal("--namespace must be set")
	}

	maybeReadImageFromFile(image)

	ctx := context.Background()

	name := *namePrefix + utilrand.String(5)

	tcr, err := requester.NewTestClusterRequest(ctx, *project, *managementCluster, *namespace, name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully authenticated to management cluster %q in GPC project %q\n", *managementCluster, *project)

	if initManifest != nil && *initManifest != "" {
		err = tcr.CreateRunnerConfigMap(ctx, *initManifest)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("configmap created with init manifiest %q\n", *initManifest)
	}

	err = tcr.CreateTestCluster(ctx, "basic", *image, flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully requested a test cluster %q in namespace %q\n", name, *namespace)

	err = tcr.MaybeSendInitialGitHubStatusUpdate(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// maybeReadImageFromFile will attempt to treat &image as a path to a file
// and read the first line, it will not fail, but if it succeeds the first
// line of the file's contents (if non-empty) will be store in &image
func maybeReadImageFromFile(image *string) {
	imageFileInfo, err := os.Stat(*image)
	if os.IsNotExist(err) || imageFileInfo.IsDir() {
		return
	}
	data, err := ioutil.ReadFile(*image)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return
	}
	if lines[0] == "" {
		return
	}
	*image = lines[0]
}
