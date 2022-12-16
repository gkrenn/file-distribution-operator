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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	krenndevv1alpha1 "github.com/gkrenn/file-distribution-operator/api/v1alpha1"
	"github.com/pkg/errors"
)

const (
	defaultRescheduleInterval = 5 * time.Minute
)

// FileDistributionConfigReconciler reconciles a FileDistributionConfig object
type FileDistributionConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=krenn.dev,resources=filedistributionconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=krenn.dev,resources=filedistributionconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=krenn.dev,resources=filedistributionconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FileDistributionConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *FileDistributionConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	namespace := req.NamespacedName.Namespace
	config, err := r.getConfig(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	fmt.Println("reconciling config: ", config.Name)

	// todo check if cr is in same namespace as operator

	// todo ensure previous jobs run through => no jobs are left

	// query all nodes and create job for each node
	fmt.Println("creating jobs for all nodes...")

	distributionJob := NewDistributionJob(r.Client, ctx, namespace, *config)
	_, err = distributionJob.setupJobOnAllNodes()
	if err != nil {
		return ctrl.Result{}, err
	}

	// wait till jobs are done or timeout is reached
	fmt.Println("waiting for all jobs to be completed...")
	if err := distributionJob.waitForAllJobsToBeCompleted(); err != nil {
		return ctrl.Result{}, err
	}

	// delete jobs
	fmt.Println("deleting all jobs...")
	if err := distributionJob.deleteAllJobsInNamespace(); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to delete all jobs: %w", err)
	}

	// reschedule after provided interval
	var interval time.Duration
	if config.Spec.RescheduleInterval == 0 {
		interval = defaultRescheduleInterval
	} else {
		interval = time.Duration(config.Spec.RescheduleInterval) * time.Minute
		fmt.Printf("rescheduling in %v minutes \n", config.Spec.RescheduleInterval)
	}
	return ctrl.Result{
		RequeueAfter: interval,
	}, nil
}

func (r FileDistributionConfigReconciler) getConfig(ctx context.Context, req ctrl.Request) (*krenndevv1alpha1.FileDistributionConfig, error) {
	var config krenndevv1alpha1.FileDistributionConfig
	err := r.Get(ctx, req.NamespacedName, &config)
	if err != nil {
		return &config, errors.WithStack(err)
	}
	return &config, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FileDistributionConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&krenndevv1alpha1.FileDistributionConfig{}).
		Complete(r)
}
