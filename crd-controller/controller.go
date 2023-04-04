package main

import (
	"context"
	"fmt"
	"time"

	emptyappv1alpha1 "k8s-client-demo/crd-controller/apis/emptyapp/v1alpha1"
	clientset "k8s-client-demo/crd-controller/client/clientset/versioned"
	informers "k8s-client-demo/crd-controller/client/informers/externalversions/emptyapp/v1alpha1"
	customlisters "k8s-client-demo/crd-controller/client/listers/emptyapp/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	appsinformers "k8s.io/client-go/informers/apps/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Controller struct {
	kubeclientset  kubernetes.Interface
	emptyclientset clientset.Interface

	emptyLister customlisters.EmptyAppLister
	emptySynced cache.InformerSynced

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface
	// recorder record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	emptyInformer informers.EmptyAppInformer) *Controller {
	controller := &Controller{
		kubeclientset:     kubeclientset,
		emptyclientset:    sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		emptyLister:       emptyInformer.Lister(),
		emptySynced:       emptyInformer.Informer().HasSynced,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	emptyAppHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				controller.queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				controller.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				controller.queue.Add(key)
			}
		},
	}

	emptyInformer.Informer().AddEventHandler(emptyAppHandler)

	return controller
}

func (c *Controller) Run(workers int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	klog.Info("Starting Deployment controller")

	// go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.emptySynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Info("Stopping Deployment controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	k := key.(string)
	fmt.Printf("handler deployment: %s\n", k)
	err := c.procesResource(k)
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}

// procesResource is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *Controller) procesResource(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Foo resource with this namespace/name
	app, err := c.emptyLister.EmptyApps(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("EmptyApp '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	imageName := app.Spec.ImageName
	if imageName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("%s: image name must be specified", key))
		return nil
	}

	deploymentName := fmt.Sprintf("%s-empty-deployment", app.Name)

	deployment, err := c.deploymentsLister.Deployments(app.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {

		deployment, err = c.kubeclientset.AppsV1().Deployments(app.Namespace).
			Create(context.TODO(), newDeployment(app, deploymentName), metav1.CreateOptions{})
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If this number of the replicas on the Foo resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if (app.Spec.Replicas != 0 && app.Spec.Replicas != *deployment.Spec.Replicas) ||
		(app.Spec.ImageName != deployment.Spec.Template.Spec.Containers[0].Image) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(app.Namespace).
			Update(context.TODO(), newDeployment(app, deploymentName), metav1.UpdateOptions{})
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the EmptyApp resource to reflect the
	// current state of the world
	err = c.updateEmptyAppStatus(app, deployment)
	if err != nil {
		return err
	}

	// Note that you also have to check the uid if you have a local controlled resource, which
	// is dependent on the actual instance, to detect that a Pod was recreated with the same name
	fmt.Printf("Sync/Add/Update for EmptyApp %s, Replicas: %d\n", app.GetName(), app.Spec.Replicas)

	return nil
}

func (c *Controller) updateEmptyAppStatus(app *emptyappv1alpha1.EmptyApp, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	appCopy := app.DeepCopy()
	appCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Foo resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.emptyclientset.CrdV1alpha1().EmptyApps(app.Namespace).
		UpdateStatus(context.TODO(), appCopy, metav1.UpdateOptions{})
	return err
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing deployment %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	klog.Infof("Dropping deployment %q out of the queue: %v", key, err)
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
						},
					},
				},
			},
		},
	}
}
