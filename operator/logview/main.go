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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	ns := os.Getenv("NAMESPACE")
	if ns == "" {
		log.Fatal("NAMESPACE must be set")
	}

	l := &logview{
		ctx:       context.Background(),
		podClient: clientSet.CoreV1().Pods(ns),
	}

	router := mux.NewRouter()

	router.HandleFunc("/logs/{pod}", l.handleLogs)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

type logview struct {
	ctx       context.Context
	podClient typedcorev1.PodInterface
}

func (l *logview) handleLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pod := vars["pod"]

	logStream, err := l.podClient.GetLogs(pod, &corev1.PodLogOptions{Follow: true}).Stream(l.ctx)
	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "cannot get log stream", http.StatusBadRequest)
		return
	}
	defer logStream.Close()
	defer log.Printf("stopped streaming logs for %q to %q", pod, r.RemoteAddr)

	log.Printf("started streaming logs for %q to %q", pod, r.RemoteAddr)
	if _, err := io.Copy(w, logStream); err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "log stream terminated unexpectedly", http.StatusInternalServerError)
		return
	}
	return
}
