package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	w, err := client.AppsV1().Deployments("bts-common").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	fmt.Println("start watch...")
	for {
		select {
		case e := <-w.ResultChan():
			fmt.Printf("event: %v, and object: %v\n", e.Type, e.Object)
		}
	}

}
