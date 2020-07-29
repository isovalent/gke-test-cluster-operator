// Copyright 2020 Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/go-logr/logr"
)

type ClientLogger struct {
	client.Client
	Log logr.Logger

	created bool
}

func NewClientLogger(mgr manager.Manager, l logr.Logger, name string) ClientLogger {
	return ClientLogger{
		Client: mgr.GetClient(),
		Log:    l.WithName("controllers").WithName(name),
	}
}

func (w *ClientLogger) createOrSkip(obj runtime.Object) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	ctx := context.Background()
	log := w.Log.WithValues("createOrSkip", key)
	client := w.Client

	remoteObj := obj.DeepCopyObject()
	getErr := client.Get(ctx, key, remoteObj)
	if apierrors.IsNotFound(getErr) {
		log.Info("will create", "obj", obj)
		w.created = true
		return client.Create(ctx, obj)
	}
	if getErr == nil {
		log.Info("already exists", "remoteObj", remoteObj)
	}
	return getErr
}
