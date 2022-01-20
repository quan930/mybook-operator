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
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	klog.Info("hello world!!!!")
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
					fmt.Println("in")
					return true
				}
			}
		}
	}
	return false
}
