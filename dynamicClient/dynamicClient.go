package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

// WebService is crd in my cluster
type WebServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []WebService `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type WebService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec WebServiceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type WebServiceSpec struct {
	Replicas int32  `json:"replicas" yaml:"replicas" `
	HostName string `json:"hostname,omitempty" yaml:"hostname,omitempty" `
	Image    string `json:"image" yaml:"image" `
}

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

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "crd.newegg.org",
		Version:  "v1",
		Resource: "webservices",
	}

	unData, err := client.Resource(gvr).Namespace("bts-common").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	serviceList := &WebServiceList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unData.UnstructuredContent(), serviceList)
	if err != nil {
		klog.Fatal(err)
	}

	for _, svc := range serviceList.Items {
		fmt.Printf("web service name: %s, replica: %d, images: %s\n", svc.ObjectMeta.Name, svc.Spec.Replicas, svc.Spec.Image)
	}

}
