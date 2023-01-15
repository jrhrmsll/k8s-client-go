package main

import (
	"context"
	"log"
	"os"
	"path"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, 10*time.Second)

	jobInformer := informerFactory.Batch().V1().Jobs()

	jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Println("ADD handler for Job: ", obj.(*batchv1.Job).Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Println("UPDATE handler for Job: ", newObj.(*batchv1.Job).Name)
		},
		DeleteFunc: func(obj interface{}) {
			log.Println("DELETE handler for Job: ", obj.(*batchv1.Job).Name)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	informerFactory.Start(ctx.Done())
	informerFactory.WaitForCacheSync(ctx.Done())

	time.Sleep(5 * time.Minute)
}
