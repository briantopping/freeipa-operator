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
    "bytes"
    "context"
    "errors"
    "k8s.io/client-go/tools/record"
    "reflect"
    "strings"
    "text/template"

    rice "github.com/GeertJohan/go.rice"
    "github.com/go-logr/logr"
    k8serr "k8s.io/apimachinery/pkg/api/errors"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"

    freeipav1alpha1 "github.com/briantopping/freeipa-operator/api/v1alpha1"
)

// IpaClusterReconciler reconciles a IpaCluster object
type IpaClusterReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=apps,resources=statefulset,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulset/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=freeipa.coglative.com,resources=ipaclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=freeipa.coglative.com,resources=ipaclusters/status,verbs=get;update;patch

func (r *IpaClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("ipacluster", req.NamespacedName)

	// Fetch the IpaCluster instance
	r.Log.Info("Starting reconcile", "Name", req.Name, "Namespace", req.Namespace, "Kind", "ipaclusters.freeipa.coglative.com")
	instance := &freeipav1alpha1.IpaCluster{}
	if err := r.Get(context.Background(), req.NamespacedName, instance); err != nil {
		if k8serr.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	list, err := r.processTemplate(instance, "service.goyaml")
	if err != nil {
		r.Log.Error(err, "Could not parse template")
		return reconcile.Result{}, err
	}
	for _, item := range list {
		if !reflect.ValueOf(item).MethodByName("DeepCopyObject").IsValid() {
			return reconcile.Result{}, errors.New("no DeepCopyObject method on object")
		}
        itemObject := item.(metaV1.Object)
        if err := controllerutil.SetControllerReference(instance, itemObject, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}
		kind := item.GetObjectKind()
		key := types.NamespacedName{Name: itemObject.GetName(), Namespace: itemObject.GetNamespace()}

		// Check if the object already exists. We need to pass the object to search for,
		// need to reflect on whatever we have to get it
		found := reflect.New(reflect.TypeOf(itemObject).Elem()).Interface().(runtime.Object)

		err = r.Get(context.Background(), key, found)
		if err != nil && k8serr.IsNotFound(err) {
			r.Log.Info("Creating object", "type", kind, "namespace", itemObject.GetNamespace(), "name", itemObject.GetName())
			err = r.Create(context.Background(), item)
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
			r.Log.Info("Updating object", "type", kind, "namespace", itemObject.GetNamespace(), "name", itemObject.GetName())
			//updateObject(foundSpec, itemSpec)
            foundSpec.Set(itemSpec)
            err = r.Update(context.Background(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// Creates a list of objects from a YAML template. These objects are introspected and applied by the caller
func (r *IpaClusterReconciler) processTemplate(cluster *freeipav1alpha1.IpaCluster, key string) ([]runtime.Object, error) {
	templateBox, err := rice.FindBox("template")
	if err != nil {
		r.Log.Error(err, "Could not open templates box")
	}
	// get file contents as string
	templateString, err := templateBox.String(key)
	if err != nil {
		r.Log.Error(err, "Could not open template")
	}
	t, err := template.New("template").Funcs(funcMaps).Parse(templateString)
	if err != nil {
		return nil, err
	}
	templateBuf := &bytes.Buffer{}
	err = t.Execute(templateBuf, cluster)
	if err != nil {
		return nil, err
	}

    files := strings.Split(templateBuf.String(), "---")
    decode := scheme.Codecs.UniversalDeserializer().Decode
    var objs []runtime.Object
    for _, f := range files {
        if f == "\n" || f == "" {
            // ignore empty cases
            continue
        }

        obj, _, e := decode([]byte(f), nil, nil)

        if e != nil {
            return nil, e
        }
        objs = append(objs, obj)
    }
    return objs, nil
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
