package controllers

import (
	"context"
	"time"

	freeipav1alpha1 "github.com/briantopping/freeipa-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const timeout = time.Second * 30
const interval = time.Second * 1

var _ = Describe("IpaCluster Controller", func() {
	It("can reconcile", func() {
		// Create the IpaCluster object and expect the Reconcile and StatefulSet to be created
		By("Expecting to create custom resource")
		instance := &freeipav1alpha1.IpaCluster{
			ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
			Spec: freeipav1alpha1.IpaClusterSpec{
				RealmName:  "EXAMPLE.COM",
				DomainName: "example.com",
				Servers: []freeipav1alpha1.Server{{
					ServerName: "server01.example.com",
					LbAddress:  "192.168.10.1",
				}},
			}}
		Expect(k8sClient.Create(context.Background(), instance)).Should(Succeed())

		//defer k8sClient.Delete(context.Background(), instance)

		By("Expecting to see custom resource created")
		Eventually(func() error {
			ipac := &freeipav1alpha1.IpaCluster{}
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo", Namespace: "default"}, ipac)
		}, timeout, interval).Should(Succeed())

		By("Expecting to see StatefulSet created")
		Eventually(func() error {
			ss := &appsv1.StatefulSet{}
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo-statefulset", Namespace: "default"}, ss)
		}, timeout, interval).Should(Succeed())

		By("Expecting to see Service created")
		service := &corev1.Service{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo-service", Namespace: "default"}, service)
		}, timeout, interval).Should(gomega.Succeed())

		// Delete the StatefulSet and expect Reconcile to be called for StatefulSet deletion
		By("Expecting to delete successfully")
		Eventually(func() error {
			f := &appsv1.StatefulSet{}
			_ = k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo-statefulset", Namespace: "default"}, f)
			return k8sClient.Delete(context.Background(), f)
		}, timeout, interval).Should(Succeed())

		ss := &appsv1.StatefulSet{}
		By("Expecting to delete finish")
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo-statefulset", Namespace: "default"}, ss)
		}, timeout, interval).ShouldNot(Succeed())

		// Manually delete StatefulSet since GC isn't enabled in the test control plane
		Eventually(func() error {
			return k8sClient.Delete(context.Background(), ss)
		}, timeout, interval).Should(gomega.MatchError("statefulsets.apps \"foo-statefulset\" not found"))

		// Manually delete services since GC isn't enabled in the test control plane
		deleteService(k8sClient, "foo-service")
		deleteService(k8sClient, "foo-service-0a")
		deleteService(k8sClient, "foo-service-0b")
	})
})

func deleteService(c client.Client, name string) bool {
	service := &corev1.Service{}
	key := types.NamespacedName{Name: name, Namespace: "default"}
	Eventually(func() error {
		return c.Get(context.Background(), key, service)
	}, timeout, interval).Should(gomega.Succeed())
	return Eventually(func() error {
		return c.Delete(context.Background(), service)
	}, timeout, interval).Should(gomega.MatchError("services \"" + name + "\" not found"))
}
