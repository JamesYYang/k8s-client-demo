package main

import (
	"context"
	emptyappv1alpha1 "k8s-client-demo/crd-controller/apis/emptyapp/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func (c *Controller) ensureDeployment(app *emptyappv1alpha1.EmptyApp) (*appsv1.Deployment, error) {
	deploymentName := app.Name

	deployment, err := c.deploymentsLister.Deployments(app.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {

		klog.Infof("create deployment for empty app: %s\n", app.Name)
		deployment, err = c.kubeclientset.AppsV1().Deployments(app.Namespace).
			Create(context.TODO(), newDeployment(app, deploymentName), metav1.CreateOptions{})

		return deployment, err
	} else if err != nil {
		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		return nil, err
	}

	// If this number of the replicas on the EmptyApp resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if (app.Spec.Replicas != 0 && app.Spec.Replicas != *deployment.Spec.Replicas) ||
		(app.Spec.ImageName != deployment.Spec.Template.Spec.Containers[0].Image) {
		klog.Infof("deployment spec changed for empty app: %s\n", app.Name)
		deployment, err = c.kubeclientset.AppsV1().Deployments(app.Namespace).
			Update(context.TODO(), newDeployment(app, deploymentName), metav1.UpdateOptions{})
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

// newDeployment creates a new Deployment for a EmptyApp resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the EmptyApp resource that 'owns' it.
func newDeployment(app *emptyappv1alpha1.EmptyApp, deploymentName string) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "emptyapp",
		"controller": app.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, emptyappv1alpha1.SchemeGroupVersion.WithKind("EmptyApp")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &app.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "empty-server",
							Image: app.Spec.ImageName,
							Ports: []corev1.ContainerPort{
								{
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8000,
									Name:          "http",
								},
							},
						},
					},
				},
			},
		},
	}
}
