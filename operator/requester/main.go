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

	description := flag.String("description", "", "definition of the purpose of this cluster")

	waitForCluster := flag.Bool("wait", false, "this flag specifies whether command should wait for cluster to be ready")

	waitTimeout := flag.Duration("wait-timeout", requester.DefaultTimeout, "how long to wait for cluster")

	flag.Parse()

	if namespace == nil || *namespace == "" {
		log.Fatal("--namespace must be set")
	}

	if image == nil || *image == "" {
		if description == nil || *description == "" {
			log.Fatal("--description must be set when --image is not set")
		}
		log.Println("cluster will be created without a test job since --image was not set")
	}

	maybeReadImageFromFile(image)
	var ctx context.Context
	var cancel context.CancelFunc

	if *waitForCluster {
		ctx = context.Background()
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), *waitTimeout)
		defer cancel()
	}

	name := *namePrefix + utilrand.String(5)

	tcr, err := requester.NewTestClusterRequest(ctx, *project, *managementCluster, *namespace, name)
	if err != nil {
		if os.Getenv("GCP_SERVICE_ACCOUNT_KEY") == "" && os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
			log.Fatal("Authentication failed in interactive mode. Run: gcloud auth application-default login")
		}
		log.Fatal(err)
	}
	log.Printf("successfully authenticated to management cluster %q in GCP project %q\n", *managementCluster, *project)

	if initManifest != nil && *initManifest != "" {
		err = tcr.CreateRunnerConfigMap(ctx, *initManifest)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("configmap created with init manifiest %q\n", *initManifest)
	}

	err = tcr.CreateTestCluster(ctx, nil, description, image, flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully requested a test cluster %q in namespace %q\n", name, *namespace)

	err = tcr.MaybeSendInitialGitHubStatusUpdate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if *waitForCluster {
		log.Printf("waiting for cluster")
		log.Printf("For credentials run:\ngcloud container clusters list --uri | grep %s | xargs gcloud container clusters get-credentials --zone x", name)
		err = tcr.WaitForTestCluster(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster ready.")
	}
}

// maybeReadImageFromFile will attempt to treat &image as a path to a file
// and read the first line, it will not fail, but if it succeeds the first
// line of the file's contents (if non-empty) will be store in &image
func maybeReadImageFromFile(image *string) {
	if image == nil || *image == "" {
		return
	}
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
