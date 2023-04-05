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
		AddFunc: controller.enqueueEmptyApp,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueueEmptyApp(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete empty app")
		},
	}

	emptyInformer.Informer().AddEventHandler(emptyAppHandler)

	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a EmptyApp resource then the handler will enqueue that EmptyApp resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

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
	klog.Infof("handler empty app: %s\n", k)
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

	// Get the EmptyApp resource with this namespace/name
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

		klog.Infof("create deployment for empty app: %s\n", app.Name)
		deployment, err = c.kubeclientset.AppsV1().Deployments(app.Namespace).
			Create(context.TODO(), newDeployment(app, deploymentName), metav1.CreateOptions{})
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If this number of the replicas on the EmptyApp resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if (app.Spec.Replicas != 0 && app.Spec.Replicas != *deployment.Spec.Replicas) ||
		(app.Spec.ImageName != deployment.Spec.Template.Spec.Containers[0].Image) {
		klog.Infof("spec changed for empty app: %s\n", app.Name)
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

	return nil
}

func (c *Controller) updateEmptyAppStatus(app *emptyappv1alpha1.EmptyApp, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	appCopy := app.DeepCopy()
	appCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the EmptyApp resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.emptyclientset.CrdV1alpha1().EmptyApps(app.Namespace).
		UpdateStatus(context.TODO(), appCopy, metav1.UpdateOptions{})

	klog.Infof("update empty app: %s status in namespace: %s with available replicas: %d\n",
		appCopy.Name, appCopy.Namespace, appCopy.Status.AvailableReplicas)
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
		klog.Infof("Error syncing empty app %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	klog.Infof("Dropping empty app %q out of the queue: %v", key, err)
}

// enqueueEmptyApp takes a EmptyApp resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than EmptyApp.
func (c *Controller) enqueueEmptyApp(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.queue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the EmptyApp resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that EmptyApp resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {

	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
	}
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a EmptyApp, we should not do anything more
		// with it.
		if ownerRef.Kind != "EmptyApp" {
			return
		}

		klog.Infoln("deployment changed for empty app")

		app, err := c.emptyLister.EmptyApps(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			return
		}

		c.enqueueEmptyApp(app)
		return
	}
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
