// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

// based on https://gist.github.com/ahmetb/548059cdbf12fb571e4e2f1e29c48997

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Fatal(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	l := &promview{
		ctx:           context.Background(),
		serviceClient: clientSet.CoreV1().Services("prom"),
	}

	router := mux.NewRouter()

	router.HandleFunc("/metrics", l.handleLogs)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

type promview struct {
	ctx           context.Context
	serviceClient typedcorev1.ServiceInterface
}

func (l *promview) handleLogs(w http.ResponseWriter, r *http.Request) {
	_, err := l.serviceClient.Get(l.ctx, "prom", metav1.GetOptions{})
	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "cannot get service", http.StatusBadRequest)
		return
	}

	metricStream, err := l.serviceClient.ProxyGet("http", "prom", "prom", "/metrics", nil).Stream(l.ctx)
	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "cannot get metric stream", http.StatusBadRequest)
		return
	}
	defer metricStream.Close()
	defer log.Printf("stopped streaming metric to %q", r.RemoteAddr)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("started streaming metrics to %q", r.RemoteAddr)
	if _, err := io.Copy(w, metricStream); err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "log stream terminated unexpectedly", http.StatusInternalServerError)
		return
	}
	return
}
