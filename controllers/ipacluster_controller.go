/*

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
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	yamlDecoder "k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	freeipav1alpha1 "github.com/briantopping/freeipa-operator/api/v1alpha1"
)

var log = logf.Log.WithName("ipacluster.controller")

// IpaClusterReconciler reconciles a IpaCluster object
type IpaClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=statefulset,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulset/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=freeipa.coglative.com,resources=ipaclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=freeipa.coglative.com,resources=ipaclusters/status,verbs=get;update;patch

func (r *IpaClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("ipacluster", req.NamespacedName)

	// Fetch the IpaCluster instance
	log.Info("Starting reconcile", "Name", req.Name, "Namespace", req.Namespace, "Kind", "ipaclusters.freeipa.coglative.com")
	instance := &freeipav1alpha1.IpaCluster{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if k8serr.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the req.
		return reconcile.Result{}, err
	}

	list, err := ProcessTemplate(instance)
	if err != nil {
	    r.Log.Error(err, "Could not parse template")
        return reconcile.Result{}, err
    }
	for _, item := range list {
		if !reflect.ValueOf(item).MethodByName("DeepCopyObject").IsValid() {
			return reconcile.Result{}, errors.New("no DeepCopyObject method on object")
		}
		if err := controllerutil.SetControllerReference(instance, item, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}
		itemObject := reflect.ValueOf(item).MethodByName("DeepCopyObject").Call(nil)[0].Interface().(runtime.Object)
		kind := itemObject.GetObjectKind().GroupVersionKind().Kind
		key := types.NamespacedName{Name: item.GetName(), Namespace: item.GetNamespace()}

		// Check if the object already exists
		found := reflect.New(reflect.TypeOf(itemObject).Elem()).Interface().(runtime.Object)
		err = r.Get(context.TODO(), key, found)
		if err != nil && k8serr.IsNotFound(err) {
			log.Info("Creating object", "type", kind, "namespace", item.GetNamespace(), "name", item.GetName())
			err = r.Create(context.TODO(), itemObject)
			if err != nil {
				return reconcile.Result{}, err
			} else {
				continue
			}
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Update the found object and write the result back if there are any changes
		itemSpec := reflect.ValueOf(item).Elem().FieldByName("Spec")
		foundSpec := reflect.ValueOf(found).Elem().FieldByName("Spec")
		if !reflect.DeepEqual(itemSpec.Interface(), foundSpec.Interface()) {
			log.Info("Updating object", "type", kind, "namespace", item.GetNamespace(), "name", item.GetName())
			updateObject(foundSpec, itemSpec)
			err = r.Update(context.TODO(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// Update object fields in a manner that respects immutables
func updateObject(dest reflect.Value, source reflect.Value) {
	name := dest.Type().Name()
	switch name {
	case "ServiceSpec":
		clusterIP := dest.FieldByName("ClusterIP").Interface().(string)
		dest.Set(source)
		dest.FieldByName("ClusterIP").SetString(clusterIP)
	default:
		dest.Set(source)
	}
}

// Creates a list of objects from a YAML template. These objects are introspected and applied by the caller
func ProcessTemplate(cluster *freeipav1alpha1.IpaCluster) ([]metaV1.Object, error) {
	templateBox, err := rice.FindBox(".")
	if err != nil {
		log.Error(err, "Could not open templates box")
	}
	// get file contents as string
	templateString, err := templateBox.String("service.tmpl")
	if err != nil {
		log.Error(err, "Could not open template")
	}
	t, err := template.New("template").Funcs(funcMaps).Parse(templateString)
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

var funcMaps = template.FuncMap{
	"seq": func(c int) []interface{} {
		return make([]interface{}, c)
	},
}

func (r *IpaClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&freeipav1alpha1.IpaCluster{}).
		Complete(r)
}
