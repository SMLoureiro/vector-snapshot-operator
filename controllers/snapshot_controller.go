package controllers

import (
    "context"
    "time"

    snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
    "github.com/SMLoureiro/vector-snapshot-operator/internal/drivers"
    "github.com/SMLoureiro/vector-snapshot-operator/internal/storage"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "k8s.io/client-go/tools/record"
)

type SnapshotReconciler struct {
    client.Client
    Scheme   *runtime.Scheme
    Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=snapshots.yourorg.io,resources=snapshots;snapshotpolicies;snapshotstorages,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=snapshots.yourorg.io,resources=snapshots/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods;secrets;events,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *SnapshotReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx).WithValues("snapshot", req.NamespacedName)

    var snap snapshotv1alpha1.Snapshot
    if err := r.Get(ctx, req.NamespacedName, &snap); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }
    if snap.Status.Phase == "Succeeded" || snap.Status.Phase == "Failed" {
        return ctrl.Result{}, nil
    }

    // Load policy
    var pol snapshotv1alpha1.SnapshotPolicy
    if err := r.Get(ctx, client.ObjectKey{Name: snap.Spec.PolicyRef}, &pol); err != nil {
        snap.Status.Phase = "Failed"
        snap.Status.Message = "policy not found: " + err.Error()
        _ = r.Status().Update(ctx, &snap)
        return ctrl.Result{}, nil
    }

    // Load storage (namespaced to Snapshot namespace)
    var store snapshotv1alpha1.SnapshotStorage
    if err := r.Get(ctx, client.ObjectKey{Name: pol.Spec.StorageRef, Namespace: snap.Namespace}, &store); err != nil {
        snap.Status.Phase = "Failed"
        snap.Status.Message = "storage not found: " + err.Error()
        _ = r.Status().Update(ctx, &snap)
        return ctrl.Result{}, nil
    }

    // Mark running
    snap.Status.Phase = "Running"
    snap.Status.Started = &metav1.Time{Time: time.Now()}
    _ = r.Status().Update(ctx, &snap)

    // Build driver + discover
    drv, err := drivers.NewDriver(ctx, pol, r.Client)
    if err != nil {
        snap.Status.Phase = "Failed"
        snap.Status.Message = "driver init: " + err.Error()
        _ = r.Status().Update(ctx, &snap)
        return ctrl.Result{}, nil
    }

    targets, err := drv.DiscoverTargets(ctx, pol.Spec.Selector)
    if err != nil {
        snap.Status.Phase = "Failed"
        snap.Status.Message = "discover: " + err.Error()
        _ = r.Status().Update(ctx, &snap)
        return ctrl.Result{}, nil
    }

    snap.Status.Total = int32(len(targets))
    _ = r.Status().Update(ctx, &snap)

    up, err := storage.NewUploader(ctx, &store)
    if err != nil {
        snap.Status.Phase = "Failed"
        snap.Status.Message = "storage init: " + err.Error()
        _ = r.Status().Update(ctx, &snap)
        return ctrl.Result{}, nil
    }

    var uris []string
    for i, t := range targets {
        uri, err := drv.SnapshotTarget(ctx, t, up)
        if err == nil {
            uris = append(uris, uri)
        }
        snap.Status.Completed = int32(i + 1)
        _ = r.Status().Update(ctx, &snap)
    }

    if len(uris) == 0 {
        snap.Status.Phase = "Failed"
        snap.Status.Message = "no successful snapshots"
        snap.Status.Ended = &metav1.Time{Time: time.Now()}
        _ = r.Status().Update(ctx, &snap)
        return ctrl.Result{}, nil
    }

    snap.Status.URIs = uris
    snap.Status.Phase = "Succeeded"
    snap.Status.Ended = &metav1.Time{Time: time.Now()}
    _ = r.Status().Update(ctx, &snap)
    logger.Info("snapshot completed", "count", len(uris))

    return ctrl.Result{}, nil
}

func (r *SnapshotReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&snapshotv1alpha1.Snapshot{}).
        Complete(r)
}
