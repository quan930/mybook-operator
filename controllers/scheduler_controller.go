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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"

	cachev1 "github.com/quan930/mybook-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SchedulerReconciler reconciles a Scheduler object
type SchedulerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.lilq.cn,resources=schedulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.lilq.cn,resources=schedulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.lilq.cn,resources=schedulers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scheduler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SchedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Info("========== start Scheduler ===============>")
	// 获取 MyBook 实例
	scheduler := &cachev1.Scheduler{}
	ctx = context.Background()

	err := r.Get(ctx, req.NamespacedName, scheduler)
	if err != nil {
		if errors.IsNotFound(err) {
			// 对象未找到
			klog.Info("Scheduler resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		klog.Error(err, "Failed to get Scheduler")
		return ctrl.Result{}, err
	}
	klog.Info("Scheduler:", scheduler)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.Scheduler{}).
		Complete(r)
}
