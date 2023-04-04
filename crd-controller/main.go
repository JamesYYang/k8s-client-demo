package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	customclient "k8s-client-demo/crd-controller/client/clientset/versioned"
	customscheme "k8s-client-demo/crd-controller/client/clientset/versioned/scheme"
	custominformers "k8s-client-demo/crd-controller/client/informers/externalversions"

	"k8s.io/client-go/kubernetes/scheme"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

func main() {
	var kubeconfig string
	var master string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()

	utilruntime.Must(customscheme.AddToScheme(scheme.Scheme))

	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	client, err := customclient.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	factory := custominformers.NewSharedInformerFactory(client, time.Hour)
	informer := factory.Crd().V1alpha1().EmptyApps()
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	deployHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}

	informer.Informer().AddEventHandler(deployHandler)

	ctrl := NewController(queue, informer.Lister(), informer.Informer())

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	defer close(stopper)

	stop := make(chan struct{})
	factory.Start(stop)
	go ctrl.Run(1, stop)

	<-stopper
	stop <- struct{}{}

	log.Println("Received signal, exiting program..")
}
