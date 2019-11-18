package controllers

import (
	"context"
	"time"

	freeipav1alpha1 "github.com/briantopping/freeipa-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const timeout = time.Second * 30
const interval = time.Second * 1

var numServers int32

var _ = Describe("IpaCluster Controller", func() {
	Context("Run without any existing resources", func() {
		It("can create the custom resource", func() {
			// Create the IpaCluster object and expect the Reconcile and StatefulSet to be created
			By("Creating the custom resource")
			instance := &freeipav1alpha1.IpaCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
				Spec: &freeipav1alpha1.IpaClusterSpec{
					RealmName:  "EXAMPLE.COM",
					DomainName: "example.com",
					Servers: []freeipav1alpha1.Server{{
						ServerName: "server01.example.com",
						LbAddress:  "192.168.10.1",
					}, {
						ServerName: "server02.example.com",
						LbAddress:  "192.168.10.2",
					}},
				}}
			Expect(k8sClient.Create(context.Background(), instance)).Should(Succeed())
			numServers = int32(len(instance.Spec.Servers))

			//defer k8sClient.Delete(context.Background(), instance)

			By("Expecting to see custom resource created")
			Eventually(func() error {
				ipac := &freeipav1alpha1.IpaCluster{}
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo", Namespace: "default"}, ipac)
			}, timeout, interval).Should(Succeed())
		})
	})

	Context("With a IPA CR, watch for child pods", func() {
		It("can create child pods", func() {
			Expect(func() error {
				ipac := &freeipav1alpha1.IpaCluster{}
				Eventually(func() error {
					return k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo", Namespace: "default"}, ipac)
				}, timeout, interval).Should(Succeed())
				listObj := &corev1.PodList{}
				Eventually(func() int32 {
					Expect(k8sClient.List(context.Background(), listObj)).To(Succeed())
					var count int32 = 0
					for _, pod := range listObj.Items {
						for _, owner := range pod.ObjectMeta.GetOwnerReferences() {
							if owner.UID == ipac.UID {
								count += 1
								break
							}
						}
					}
					return count
				}, timeout, interval).Should(Equal(numServers))
				return nil
			}()).To(Succeed())
		})
	})
	//
	//Context("With a IPA CR and master resources created", func() {
	//    It("can delete all resources generated for IPA CR", func() {
	//
	//        // Manually delete StatefulSet since GC isn't enabled in the test control plane
	//        Eventually(func() error {
	//            ipac := &freeipav1alpha1.IpaCluster{}
	//            _ = k8sClient.Get(context.Background(), types.NamespacedName{Name: "foo", Namespace: "default"}, ipac)
	//            return k8sClient.Delete(context.Background(), ipac)
	//        }, timeout, interval).Should(Succeed())
	//
	//        // Manually delete services since GC isn't enabled in the test control plane
	//        deleteService(k8sClient, "foo-service")
	//        deleteService(k8sClient, "foo-service-0a")
	//        deleteService(k8sClient, "foo-service-0b")
	//    })
	//})
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
