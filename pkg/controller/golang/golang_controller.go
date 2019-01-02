package golang

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"

	golangv1alpha1 "github.com/craicoverflow/golang-operator/pkg/apis/golang/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_golang")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Golang Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGolang{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("golang-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Golang
	err = c.Watch(&source.Kind{Type: &golangv1alpha1.Golang{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Golang
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &golangv1alpha1.Golang{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileGolang{}

// ReconcileGolang reconciles a Golang object
type ReconcileGolang struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Golang object and makes changes based on the state read
// and what is in the Golang.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileGolang) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", "Request.Name", request.Name)
	reqLogger.Info("Reconciling Golang")

	// Fetch the Memcached instance
	golang := &golangv1alpha1.Golang{}
	err := r.client.Get(context.TODO(), request.NamespacedName, golang)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object is not found, could have been deleted after reconcile request
			reqLogger.Info("Golang resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request
		reqLogger.Error(err, "failed to get Golang")
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: golang.Name, Namespace: golang.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// define a new deployment
		dep := r.deploymentForGo(golang)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Info("Failed to create a new deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
	}

	size := golang.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err := r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the Memcached status with the pod names
	// List the pods for this memcached's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForGolang(golang.Name))
	listOpts := &client.ListOptions{Namespace: golang.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOpts, podList)
	if err != nil {
		reqLogger.Error(err, "failed to list pods", "Golang.Namespace", golang.Namespace)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Pods if needed
	if !reflect.DeepEqual(podNames, golang.Status.Pods) {
		golang.Status.Pods = podNames
		err := r.client.Status().Update(context.TODO(), golang)
		if err != nil {
			reqLogger.Error(err, "failed to update Memcached status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileGolang) deploymentForGo(g *golangv1alpha1.Golang) *appsv1.Deployment {
	labels := labelsForGolang(g.Name)
	replicas := g.Spec.Size

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      g.Name,
			Namespace: g.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: g.Spec.Image,
						Name:  g.Spec.Name,
					}},
				},
			},
		},
	}
	// set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(g, dep, r.scheme)
	return dep
}

func labelsForGolang(name string) map[string]string {
	return map[string]string{"app": "golang", "golang_cr": name}
}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
