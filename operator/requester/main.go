package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/isovalent/gke-test-cluster-management/operator/pkg/requester"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

func main() {

	image := flag.String("image", "", "name of the image to use for driving the tests")
	imageTagFilePath := flag.String("image-tag-from-file", "", "read the name of the test image from a file")

	namespace := flag.String("namespace", "", "namespace to use")

	namePrefix := flag.String("name-prefix", "test-", "name prefix for the test cluster")

	project := flag.String("project", requester.DefaultProject, "GCP project")
	managementCluster := flag.String("management-cluster", requester.DefaultManagementCluster, "name of the management cluster")

	initManifest := flag.String("init-manifest", "", "path to manifest to use to initialise the cluster")

	description := flag.String("description", "", "definition of the purpose of this cluster")

	waitForCluster := flag.Bool("wait", false, "once cluster has been requested, wait for it to become ready")

	waitTimeout := flag.Duration("wait-timeout", requester.DefaultTimeout, "how long to wait for cluster")

	debug := flag.Bool("debug", false, "enable interactive test debug mode with 'kubectl exec'")

	flag.Parse()

	if namespace == nil || *namespace == "" {
		log.Fatal("--namespace must be set")
	}

	if image == nil || imageTagFilePath == nil {
		if description == nil || *description == "" {
			log.Fatal("--description must be set when neither --image nor --image-from-file are set")
		}
		log.Println("cluster will be created without a test job since image was not given")
	}

	if imageTagFilePath != nil {
		if err := readImageFromFile(image, imageTagFilePath); err != nil {
			log.Fatalf("cannot parse image tag from file: %s", err)
		}
	}

	log.Printf("will use image %q", *image)

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
			log.Println("authentication failed in interactive mode")
			log.Fatal("to fix this, run 'gcloud auth application-default login'")
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

	if *debug {
		log.Printf("requesting cluster in debug mode")
		err = tcr.CreateTestCluster(ctx, nil, description, image, "/bin/sh", "-c", "while :; do sleep 120; done")
	} else {
		err = tcr.CreateTestCluster(ctx, nil, description, image, flag.Args()...)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully requested a test cluster %q in namespace %q\n", name, *namespace)

	err = tcr.MaybeSendInitialGitHubStatusUpdate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if *waitForCluster {
		log.Printf("waiting for test cluster %q in namespace %s to become ready", name, *namespace)
		cluster, err := tcr.WaitForTestCluster(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("test cluster %q in namespace %q is ready", name, *namespace)
		log.Printf("for credentials run:\ngcloud container clusters get-credentials --region %s --project %s %s", *cluster.Spec.Region, *project, name)
	}

	if *debug {
		log.Print("Debug mode detected, exec into pod using following commands")
		if !*waitForCluster { // get-credentials instruction not printed yet
			log.Printf("gcloud container clusters list --uri | grep %s | xargs gcloud container clusters get-credentials --zone x", name)
		}
		log.Printf(`
kubectl get jobs -n %[1]s -o json | \
jq -r '.items[].metadata.name' | \
grep %[2]s | \
xargs -I "{}" kubectl get pods -n %[1]s -l job-name="{}" -o json | \
jq -r '.items[].metadata.name' | \
xargs -I "{}" kubectl exec -n %[1]s -ti "{}" sh`, *namespace, name)
	}

	log.Printf("for cluster cleanup run following command against management cluster:\nkubectl delete tcg %s -n %s", name, *namespace)
}

// readImageFromFile will read image name from filePath and either set image or return an error;
// the file is expected to contain image name on the first line
func readImageFromFile(image, filePath *string) error {
	imageFileInfo, err := os.Stat(*filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%q does not exist", *filePath)
	}

	if imageFileInfo.IsDir() {
		return fmt.Errorf("%q is a directory", *filePath)
	}

	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return fmt.Errorf("%q must have at least one line", *filePath)
	}
	parsedImage := strings.TrimSpace(lines[0])
	if parsedImage == "" {
		return fmt.Errorf("first line in %q is empty", *filePath)
	}
	*image = parsedImage
	return nil
}
