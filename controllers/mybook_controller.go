/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	cachev1 "github.com/quan930/mybook-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MyBookReconciler reconciles a MyBook object
type MyBookReconciler struct {
	Log logr.Logger
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.lilq.cn,resources=mybooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.lilq.cn,resources=mybooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.lilq.cn,resources=mybooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyBook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
//+kubebuilder:rbac:groups=cache.lilq.cn,resources=mybook,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.lilq.cn,resources=mybook/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
func (r *MyBookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("========== start ===============>")
	// 获取 MyBook 实例
	mybook := &cachev1.MyBook{}
	ctx = context.Background()
	klog.Info("hello", "world2")

	err := r.Get(ctx, req.NamespacedName, mybook)
	if err != nil {
		if errors.IsNotFound(err) {
			// 对象未找到
			klog.Info("MyBook resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		klog.Error(err, "Failed to get MyBook")
		return ctrl.Result{}, err
	}
	klog.Info("MyBook:", mybook)

	klog.Info("deployment ........ init =>")
	found := &v1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: "mybook-operator-mybook-server", Namespace: "mybook-operator-system"}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForMyBook(mybook)

		klog.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			klog.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		klog.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}
	klog.Info("deployment ........ finish =>")

	// Update status.Nodes if needed
	if contains(mybook.Status.History, mybook.Spec) {
		return ctrl.Result{}, nil
	}

	mybook.Status.History = append(mybook.Status.History, mybook.Spec)
	err = r.Status().Update(ctx, mybook)
	if err != nil {
		klog.Error(err, "Failed to update Memcached status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyBookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//控制器监视的资源
	return ctrl.NewControllerManagedBy(mgr).
		//将 MyBook 类型指定为要监视的主要资源
		For(&cachev1.MyBook{}).
		//将 Deployments 类型指定为要监视的辅助资源。对于每个部署类型的添加/更新/删除事件，事件处理程序会将每个事件映射到Request部署所有者的协调
		Owns(&v1.Deployment{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

func contains(array interface{}, val interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					return true
				}
			}
		}
	}
	return false
}

//deploymentForMyBook 部署服务
func (r *MyBookReconciler) deploymentForMyBook(m *cachev1.MyBook) *v1.Deployment {
	ls := labelsForMyBook(m.Name)
	replicas := int32(2)

	dep := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mybook-operator-mybook-server",
			Namespace: "mybook-operator-system",
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: v12.PodSpec{
					Containers: []v12.Container{{
						Image: "gcr.io/google-samples/node-hello:1.0",
						Name:  "book-server",
						//Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []v12.ContainerPort{{
							ContainerPort: 8080,
							Name:          "nginx",
						}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

func labelsForMyBook(name string) map[string]string {
	return map[string]string{"app": "mybook", "mybook_cr": name}
}
