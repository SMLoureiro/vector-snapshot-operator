package controllers

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	snapshotv1alpha1 "github.com/SMLoureiro/vector-snapshot-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var (
	testEnv *envtest.Environment
	k8sClient client.Client
	ctx context.Context
	cancel context.CancelFunc
	scheme = runtime.NewScheme()
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(snapshotv1alpha1.AddToScheme(scheme))

	testEnv = &envtest.Environment{
		CRDInstallOptions: envtest.CRDInstallOptions{
			Paths: []string{
				filepath.Join("..", "config", "crd", "bases"),
			},
		},
	}

	var err error
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())

	// start manager with our reconciler
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())

	rec := &SnapshotPolicyReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("test"),
	}
	Expect(rec.SetupWithManager(mgr)).To(Succeed())

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		defer GinkgoRecover()
		Expect(mgr.Start(ctx)).To(Succeed())
	}()
})

var _ = AfterSuite(func() {
	cancel()
	Expect(testEnv.Stop()).To(Succeed())
})

var _ = It("creates a Snapshot on due schedule", func() {
	pol := &snapshotv1alpha1.SnapshotPolicy{}
	pol.Name = "every-second"
	pol.Spec.Schedule = "* * * * *" // every minute; fast enough for smoke
	pol.Spec.Retention = 1
	Expect(k8sClient.Create(context.Background(), pol)).To(Succeed())

	// Give the reconciler a moment to run
	time.Sleep(2 * time.Second)

	// A Snapshot should exist (created by the policy reconciler)
	snaps := &snapshotv1alpha1.SnapshotList{}
	Expect(k8sClient.List(context.Background(), snaps)).To(Succeed())
	Expect(len(snaps.Items)).To(BeNumerically(">=", 1))
})
