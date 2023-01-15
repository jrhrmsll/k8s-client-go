package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"k8s.io/client-go/discovery"
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
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("the kubeconfig cannot be loaded: %v\n", err)
		}
	}

	discovery, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	version, err := discovery.ServerVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)
}
