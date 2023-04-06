package main

import (
	"context"
	emptyappv1alpha1 "k8s-client-demo/crd-controller/apis/emptyapp/v1alpha1"

	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func (c *Controller) ensureService(app *emptyappv1alpha1.EmptyApp) (*corev1.Service, error) {
	svcName := app.Name

	svc, err := c.servicesLister.Services(app.Namespace).Get(svcName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {

		klog.Infof("create service for empty app: %s\n", app.Name)
		svc, err = c.kubeclientset.CoreV1().Services(app.Namespace).
			Create(context.TODO(), newService(app, svcName), metav1.CreateOptions{})

		return svc, err
	} else if err != nil {
		return nil, err
	}

	if app.Spec.SvcPort != 0 && app.Spec.Replicas != *&svc.Spec.Ports[0].Port {
		klog.Infof("service spec changed for empty app: %s\n", app.Name)
		svc, err = c.kubeclientset.CoreV1().Services(app.Namespace).
			Update(context.TODO(), newService(app, svcName), metav1.UpdateOptions{})
	}

	if err != nil {
		return nil, err
	}

	return svc, nil
}

func newService(app *emptyappv1alpha1.EmptyApp, svcName string) *corev1.Service {
	labels := map[string]string{
		"app":        "emptyapp",
		"controller": app.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, emptyappv1alpha1.SchemeGroupVersion.WithKind("EmptyApp")),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("http"),
					Port:       app.Spec.SvcPort,
				},
			},
		},
	}
}
