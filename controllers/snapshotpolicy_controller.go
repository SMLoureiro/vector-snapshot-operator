package controllers

import (
	"context"
	"time"

	snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
	"github.com/robfig/cron/v3"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const policyFinalizer = "snapshots.loureiro.io/policy-protect"

type SnapshotPolicyReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=snapshots.yourorg.io,resources=snapshotpolicies;snapshots,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=snapshots.yourorg.io,resources=snapshotpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *SnapshotPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("snapshotpolicy", req.NamespacedName)

	var pol snapshotv1alpha1.SnapshotPolicy
	if err := r.Get(ctx, req.NamespacedName, &pol); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Finalizers
	if pol.DeletionTimestamp != nil {
		controllerutil.RemoveFinalizer(&pol, policyFinalizer)
		if err := r.Update(ctx, &pol); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	if !controllerutil.ContainsFinalizer(&pol, policyFinalizer) {
		controllerutil.AddFinalizer(&pol, policyFinalizer)
		if err := r.Update(ctx, &pol); err != nil {
			return ctrl.Result{}, err
		}
	}

	sched, err := cron.ParseStandard(pol.Spec.Schedule)
	if err != nil {
		apimeta.SetStatusCondition(&pol.Status.Conditions, metav1.Condition{
			Type: "Ready", Status: metav1.ConditionFalse, Reason: "BadSchedule", Message: err.Error(),
		})
		_ = r.Status().Update(ctx, &pol)
		return ctrl.Result{}, nil
	}

	now := time.Now()
	var next time.Time
	if pol.Status.NextRun != nil {
		next = pol.Status.NextRun.Time
	} else {
		next = sched.Next(now.Add(-time.Second))
	}
	due := !next.After(now)

	if due {
		snap := snapshotv1alpha1.Snapshot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "snap-" + time.Now().UTC().Format("20060102-150405"),
				Namespace: "default",
				Labels: map[string]string{
					"snapshots.yourorg.io/policy": pol.Name,
				},
			},
			Spec: snapshotv1alpha1.SnapshotSpec{
				PolicyRef: pol.Name,
				Reason:    "scheduled",
			},
		}
		if err := r.Create(ctx, &snap); err != nil {
			logger.Error(err, "create Snapshot")
		} else {
			r.Recorder.Eventf(&pol, "Normal", "SnapshotCreated", "Snapshot %s created", snap.Name)
			pol.Status.LastRun = &metav1.Time{Time: now}
		}
		next = sched.Next(now)
		pol.Status.NextRun = &metav1.Time{Time: next}
		_ = r.Status().Update(ctx, &pol)
	}

	return ctrl.Result{RequeueAfter: time.Until(next)}, nil
}

func (r *SnapshotPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&snapshotv1alpha1.SnapshotPolicy{}).
		Owns(&snapshotv1alpha1.Snapshot{}).
		Complete(r)
}
