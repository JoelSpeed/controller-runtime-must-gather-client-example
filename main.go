package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/library-go/pkg/manifestclient"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	mustGatherRoundTripper, err := manifestclient.NewRoundTripper(filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "openshift", "library-go", "pkg", "manifestclienttest", "testdata", "must-gather-01"))
	if err != nil {
		panic(err)
	}

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := configv1.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := appsv1.AddToScheme(scheme); err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Transport: mustGatherRoundTripper,
	}

	k8sClient, err := client.New(&rest.Config{}, client.Options{
		HTTPClient: httpClient,
		Scheme:     scheme,
	})
	if err != nil {
		panic(err)
	}

	// Core group resource.
	configmap := &corev1.ConfigMap{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: "openshift-config", Name: "cloud-provider-config"}, configmap); err != nil {
		panic(err)
	}
	fmt.Printf("ConfigMap: %+v\n", configmap)

	// Cluster scoped resource.
	infrastructure := &configv1.Infrastructure{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Name: "cluster"}, infrastructure); err != nil {
		panic(err)
	}
	fmt.Printf("Infrastructure: %+v\n", infrastructure)

	// Namespaced resource.
	deployment := &appsv1.Deployment{}
	if err := k8sClient.Get(context.Background(), client.ObjectKey{Namespace: "openshift-machine-api", Name: "machine-api-operator"}, deployment); err != nil {
		panic(err)
	}
	fmt.Printf("Deployment: %+v\n", deployment)
}
