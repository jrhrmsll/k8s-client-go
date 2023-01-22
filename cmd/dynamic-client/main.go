package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := rest.InClusterConfig()

	if err != nil {
		// fallback to kubeconfig
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		kubeconfig := path.Join(home, ".kube/config")
		if env, ok := os.LookupEnv("KUBECONFIG"); ok {
			kubeconfig = env
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("the kubeconfig cannot be loaded: %v\n", err)
		}
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "programming.kubernetes.dev",
		Version:  "v1",
		Resource: "faultcollections",
	}

	faultCollection, err := dc.Resource(gvr).Namespace("default").Get(context.TODO(), "ecommerce", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	description, ok, err := unstructured.NestedString(faultCollection.Object, "spec", "description")
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Println("nested field 'spec.description' not found")
	}

	fmt.Println(description)
}
