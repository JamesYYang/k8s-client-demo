package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	customclient "k8s-client-demo/crd-controller/client/clientset/versioned"
	customscheme "k8s-client-demo/crd-controller/client/clientset/versioned/scheme"
	custominformers "k8s-client-demo/crd-controller/client/informers/externalversions"

	kubeinformers "k8s.io/client-go/informers"

	"k8s.io/client-go/kubernetes/scheme"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	client, err := customclient.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	factory := custominformers.NewSharedInformerFactory(client, time.Hour)

	ctrl := NewController(kubeClient, client,
		kubeInformerFactory.Apps().V1().Deployments(),
		factory.Crd().V1alpha1().EmptyApps())

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	defer close(stopper)

	stop := make(chan struct{})
	kubeInformerFactory.Start(stop)
	factory.Start(stop)
	go ctrl.Run(1, stop)

	<-stopper
	stop <- struct{}{}

	klog.Infoln("Received signal, exiting program..")
}
