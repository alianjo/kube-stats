/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"intodevops.com/kube-stats/api/v1alpha1"
	kubestatsv1alpha1 "intodevops.com/kube-stats/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatsReconciler reconciles a Stats object
type StatsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kube-stats.intodevops.com,resources=stats,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kube-stats.intodevops.com,resources=stats/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kube-stats.intodevops.com,resources=stats/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Stats object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *StatsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	stats := v1alpha1.Stats{}
	if err := r.Get(ctx, req.NamespacedName, &stats); err != nil {
		fmt.Println("\nStats", stats, "\t does not exist")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println("\nStats", stats, "\t exist")

	fmt.Println("this is the node name: ", stats.Spec.NodeName)
	node := corev1.Node{}
	if err := r.Get(ctx, client.ObjectKey{Name: stats.Spec.NodeName}, &node); err != nil {
		fmt.Println("\nNode", stats.Spec.NodeName, "\t does not exist")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println("\nNode", stats.Spec.NodeName, "\t exist")
	Replicas := int32(1)
	// Deploy exporter on the node
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stats-exporter",
			Namespace: stats.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "stats-exporter",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "stats-exporter",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "stats-exporter",
							Image: "alianjo178/k8s-exporter:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9100,
									Name:          "http",
								},
							},
						},
					},
					NodeName: stats.Spec.NodeName,
				},
			},
		},
	}

	if err := r.Create(ctx, deployment); err != nil {
		return ctrl.Result{}, err
	}
	fmt.Println("\n Deployment has been created")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubestatsv1alpha1.Stats{}).
		Complete(r)
}
