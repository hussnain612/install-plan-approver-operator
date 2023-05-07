package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
)

func (r *SubscriptionReconciler) handleSubscription(ctx context.Context, log logr.Logger, subscription *v1alpha1.Subscription) (ctrl.Result, error) {
	// Check installPlanApproval
	if subscription.Spec.InstallPlanApproval == v1alpha1.ApprovalAutomatic {
		log.Info("InstallPlanApproval is set to Automatic, skipping...")
		return ctrl.Result{}, nil
	}

	// Check startingCSV
	if subscription.Spec.StartingCSV == "" {
		log.Info("StartingCSV is not provided, skipping...")
		return ctrl.Result{}, nil
	}

	desiredCSV := subscription.Spec.StartingCSV

	// Check subscription state
	if subscription.Status.State != v1alpha1.SubscriptionStateUpgradePending {
		log.Info("Subscription state is not UpgradePending, skipping...")
		return ctrl.Result{}, nil
	}

	// Get install plan
	installPlan := &v1alpha1.InstallPlan{}

	installPlanName := subscription.Status.InstallPlanRef.Name
	installPlanNamespace := subscription.Status.InstallPlanRef.Namespace

	// Updated log
	installPlanNamespacedName := types.NamespacedName{}
	installPlanNamespacedName.Name = installPlanName
	installPlanNamespacedName.Namespace = installPlanNamespace
	log = r.Log.WithValues("installPlan", installPlanNamespacedName)

	err := r.Get(ctx, installPlanNamespacedName, installPlan)
	if err != nil {
		if errors.IsNotFound(err) {
			// InstallPlan is not created yet
			log.Info("InstallPlan is not created yet, reconciling after 10 sec")
			return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, nil
		}
		log.Error(err, "failed to get installPlan")
		return ctrl.Result{}, err
	}

	// Check approved state
	if installPlan.Spec.Approved {
		log.Info("InstallPlan is already approved, skipping...")
		return ctrl.Result{}, nil
	}

	installPlanBasePatch := client.MergeFrom(installPlan.DeepCopy())
	approved := false

	// Approve InstallPlan if desiredCSV is mentioned
	for _, installPlanName := range installPlan.Spec.ClusterServiceVersionNames {
		if installPlanName == desiredCSV {
			approved = true
			log.Info("Approving installPlan...")
			installPlan.Spec.Approved = true
			err = r.Patch(ctx, installPlan, installPlanBasePatch)
			if err != nil {
				log.Error(err, "failed to patch installPlan")
				return ctrl.Result{}, err
			}
			log.Info("InstallPlan successfully approved")
		}
	}

	if !approved {
		log.Info("InstallPlan does not contain currentCSV present in subscription")
	}

	return ctrl.Result{}, nil
}
