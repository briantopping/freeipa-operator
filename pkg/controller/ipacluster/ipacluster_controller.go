// Copyright 2019 The FreeIPA Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipacluster

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"text/template"

	freeipav1alpha1 "github.com/briantopping/freeipa-operator/pkg/apis/freeipa/v1alpha1"
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	yamlDecoder "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("ipacluster.controller")

// Add creates a new IpaCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIpaCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ipacluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to IpaCluster
	err = c.Watch(&source.Kind{Type: &freeipav1alpha1.IpaCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Uncomment watch a StatefulSet created by IpaCluster - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &freeipav1alpha1.IpaCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileIpaCluster{}

// ReconcileIpaCluster reconciles a IpaCluster object
type ReconcileIpaCluster struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a IpaCluster object and makes changes based on the state read
// and what is in the IpaCluster.Spec
// Automatically generate RBAC rules to allow the Controller to read and write StatefulSets
// +kubebuilder:rbac:groups=apps,resources=statefulset,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulset/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=freeipa.coglative.com,resources=ipaclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=freeipa.coglative.com,resources=ipaclusters/status,verbs=get;update;patch
func (r *ReconcileIpaCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the IpaCluster instance
	log.Info("Starting reconcile", "Name", request.Name, "Namespace", request.Namespace, "Kind", "ipaclusters.freeipa.coglative.com")
	instance := &freeipav1alpha1.IpaCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if k8serr.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define the desired StatefulSets object
	//terminationGrace := int64(300)
	//item := &appsv1.StatefulSet{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:      instance.Name + "-statefulset",
	//		Namespace: instance.Namespace,
	//	},
	//	Spec: appsv1.StatefulSetSpec{
	//		Selector: &metav1.LabelSelector{
	//			MatchLabels: map[string]string{"statefulset": instance.Name + "-statefulset"},
	//		},
	//		ServiceName: "test",
	//		Template: corev1.PodTemplateSpec{
	//			ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"statefulset": instance.Name + "-statefulset"}},
	//			Spec: corev1.PodSpec{
	//				PriorityClassName: "high-priority",
	//				TerminationGracePeriodSeconds: &terminationGrace,
	//				Containers: []corev1.Container{
	//					{
	//						Name:  "freeipa-server",
	//						Image: "freeipa/freeipa-server:centos-7",
	//						ImagePullPolicy: corev1.PullIfNotPresent,
	//						Args: "ipa-replica-install",
	//
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	list, err := ProcessTemplate(instance) // returns ([]metaV1.Object, error)
	for _, item := range list {
		if !reflect.ValueOf(item).MethodByName("DeepCopyObject").IsValid() {
			return reconcile.Result{}, errors.New("no DeepCopyObject method on object")
		}
		if err := controllerutil.SetControllerReference(instance, item, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		itemObject := reflect.ValueOf(item).MethodByName("DeepCopyObject").Call(nil)[0].Interface().(runtime.Object)
		kind := itemObject.GetObjectKind().GroupVersionKind().Kind

		// Check if the StatefulSet already exists
		found := reflect.New(reflect.TypeOf(itemObject).Elem()).Interface().(runtime.Object)
		err = r.Get(context.TODO(), types.NamespacedName{Name: item.GetName(), Namespace: item.GetNamespace()}, found)
		if err != nil && k8serr.IsNotFound(err) {
			log.Info("Creating object", "type", kind, "namespace", item.GetNamespace(), "name", item.GetName())
			err = r.Create(context.TODO(), itemObject)
			return reconcile.Result{}, err
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Update the found object and write the result back if there are any changes
		itemSpec := reflect.ValueOf(item).Elem().FieldByName("Spec")
		foundSpec := reflect.ValueOf(found).Elem().FieldByName("Spec")
		if !reflect.DeepEqual(itemSpec.Interface(), foundSpec.Interface()) {
			log.Info("Updating object", "type", kind, "namespace", item.GetNamespace(), "name", item.GetName())
			foundSpec.Set(itemSpec)
			err = r.Update(context.TODO(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}
	return reconcile.Result{}, nil
}

func ProcessTemplate(cluster *freeipav1alpha1.IpaCluster) ([]metaV1.Object, error) {
	t, err := template.New("template").Parse(Template)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, cluster)
	if err != nil {
		return nil, err
	}

	reader := yamlDecoder.NewYAMLReader(bufio.NewReaderSize(buf, 4096))
	result := []metaV1.Object(nil)
	for {
		// Read a single YAML object
		b, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Unmarshal the object enough to read the Kind field
		var meta metaV1.TypeMeta
		if err := yaml.Unmarshal(b, &meta); err != nil {
			return nil, err
		}

		switch meta.Kind {
		case "StatefulSet":
			var ss appsv1.StatefulSet
			err = yaml.Unmarshal(b, &ss)
			if err != nil {
				return nil, err
			}
			result = append(result, &ss)

		case "Service":
			var service corev1.Service
			err = yaml.Unmarshal(b, &service)
			if err != nil {
				return nil, err
			}
			result = append(result, &service)
		}
	}
	return result, nil
}
