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
    "context"
    "fmt"
    "github.com/briantopping/freeipa-operator/api"
    freeipav1alpha1 "github.com/briantopping/freeipa-operator/api/v1alpha1"
    template2 "github.com/briantopping/freeipa-operator/controllers/template"
    "github.com/go-logr/logr"
    k8serr "k8s.io/apimachinery/pkg/api/errors"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/tools/record"
    "reflect"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IpaClusterReconciler reconciles a IpaCluster object
type IpaClusterReconciler struct {
    Client   client.Client
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

    // Fetch the IpaCluster cluster
	r.Log.Info("Starting reconcile", "Name", req.Name, "Namespace", req.Namespace, "Kind", "ipaclusters.freeipa.coglative.com")
    if e, done := r.operatorHousekeeping(req, cluster); done {
        return ctrl.Result{}, e
    }

    statusByGuid := make(map[string]freeipav1alpha1.ServerStatus)
    if e, done := r.buildStatus(cluster, statusByGuid); done {
        return ctrl.Result{}, e
    }

    for _, podStatus := range statusByGuid {
        podStatus.Dispatch(r)
    }

    // FIXME: copy the podStatus to the cluster variable for storage
    err := r.Client.Update(context.Background(), cluster)
    return ctrl.Result{}, err

    //// recover cache from Status
    //statusByName := cluster.Status.ServerStatus
    //statusByGuid := make(map[string]freeipav1alpha1.ServerStatus)
    //for _, status := range statusByName {
    //    if status.GUID != nil {
    //        statusByGuid[*status.GUID] = status
    //    }
    //}
    //
    //// cache initial pod contents
    //var anyErr, err error
    //initial := make(map[string]runtime.Object)
    //listObj := &corev1.PodList{}
    //err = r.Client.List(context.Background(), listObj, client.MatchingLabels{"freeipa-operator-parent": string(cluster.ObjectMeta.GetUID())})
    //for _, pod := range listObj.Items {
    //    guid := string(pod.GetUID())
    //    if _, ok := statusByGuid[guid]; !ok {
    //        missing[guid] = &pod
    //    }
    //}
    //
    //// iterate on all the serverStatus specs. Iterate first to create or update all servers...
    //var reason string
    //for _, serverStatus := range cluster.Status.ServerStatus {
    //    if initial[serverStatus.ResourceVersion] != nil {
    //        delete(initial, serverStatus.ResourceVersion)
    //        err = r.update(serverStatus)
    //        reason = "Update"
    //    } else {
    //        err = r.create(cluster, &serverStatus)
    //        reason = "Create"
    //    }
    //
    //    if err == nil {
    //        serverStatus.ResourceVersion = cluster.ObjectMeta.ResourceVersion
    //        r.Recorder.Event(cluster, "Normal", reason, fmt.Sprintf("%s serverStatus %s success", reason, serverStatus.ServerName))
    //    } else {
    //        anyErr = err
    //        r.Recorder.Event(cluster, "Warning", reason, fmt.Sprintf("%s serverStatus %s failed, error is \"%s\"", reason, serverStatus.ServerName, err.Error()))
    //    }
    //}
    //
    //// ... now delete whatever remains
    //for _, podObj := range initial {
    //    getObjectDescription(podObj)
    //    if delErr := r.Client.Delete(context.Background(), podObj); delErr == nil {
    //        r.Recorder.Event(cluster, "Warning", "Delete", fmt.Sprintf("Deleted reconciled Pod %s", getObjectDescription(podObj)))
    //    } else {
    //        r.Recorder.Event(cluster, "Warning", "Delete", fmt.Sprintf("Delete reconciled Pod %s failed, error is \"%s\"", getObjectDescription(podObj), delErr.Error()))
    //    }
    //}
    //
    ////
    //if anyErr != nil {
    //    return ctrl.Result{}, anyErr
    //}
    ////if !cluster.IsSubmitted() {
    ////    if err := r.submit(cluster); err != nil {
    ////        return ctrl.Result{}, fmt.Errorf("error when submitting run: %v", err)
    ////    }
    ////    r.Recorder.Event(cluster, "Normal", "Submitted", "Object is submitted")
    ////}
    ////
    ////if cluster.IsSubmitted() {
    ////    if err := r.refresh(cluster); err != nil {
    ////        return ctrl.Result{}, fmt.Errorf("error when refreshing run: %v", err)
    ////    }
    ////    r.Recorder.Event(cluster, "Normal", "Refreshed", "Object is refreshed")
    ////}
    //
    ////return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *IpaClusterReconciler) operatorHousekeeping(req ctrl.Request) (error, bool, api.StateReconciler) {
    cluster := &unstructured.Unstructured{}
    if err := r.Client.Get(context.Background(), req.NamespacedName, cluster); err != nil {
        if k8serr.IsNotFound(err) {
            return nil, true, nil
        }
        return err, true, nil
    }
    if cluster.IsBeingDeleted() {
        if err := r.removeFinalizer(cluster); err != nil {
            return fmt.Errorf("error when handling finalizer: %v", err), true, nil
        }
        r.Recorder.Event(cluster, "Normal", "Deleted", "Object finalizer is deleted")
        return nil, true, nil
    }
    if !cluster.HasFinalizer(freeipav1alpha1.IpaClusterFinalizerName) {
        if err := r.addFinalizer(cluster); err != nil {
            return fmt.Errorf("error when adding finalizer: %v", err), true, nil
        }
        r.Recorder.Event(cluster, "Normal", "Added", "Object finalizer is added")
        return nil, true, nil
    }
    return nil, false, nil
}

func (r *IpaClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&freeipav1alpha1.IpaCluster{}).
        Complete(r)
}

func (r *IpaClusterReconciler) removeFinalizer(cluster *freeipav1alpha1.IpaCluster) error {
    cluster.RemoveFinalizer(freeipav1alpha1.IpaClusterFinalizerName)
    return r.Client.Update(context.Background(), cluster)
}

func (r *IpaClusterReconciler) addFinalizer(cluster *freeipav1alpha1.IpaCluster) error {
    cluster.AddFinalizer(freeipav1alpha1.IpaClusterFinalizerName)
    return r.Client.Update(context.Background(), cluster)
}

func (r *IpaClusterReconciler) create(cluster *freeipav1alpha1.IpaCluster, server *freeipav1alpha1.Server) error {
    renderContext := &template2.RenderContext{Key: "pod-set.goyaml", Cluster: cluster, Server: server}
    list, err := renderContext.ProcessTemplate()
    if err != nil {
        return err
    }
    for _, item := range list {
        itemObject := item.(metaV1.Object)
        if err := controllerutil.SetControllerReference(cluster, itemObject, r.Scheme); err != nil {
            return err
        }
        key := types.NamespacedName{Name: itemObject.GetName(), Namespace: itemObject.GetNamespace()}

        // Check if the object already exists. We need to pass the object to search for,
        // need to reflect on whatever we have to get it
        found := reflect.New(reflect.TypeOf(itemObject).Elem()).Interface().(runtime.Object)

        err = r.Client.Get(context.Background(), key, found)
        if err != nil && k8serr.IsNotFound(err) {
            r.Log.Info("Creating object", "id", getObjectDescription(item))
            err = r.Client.Create(context.Background(), item)
            if err != nil {
                return err
            } else {
                continue
            }
        } else if err == nil {
            return fmt.Errorf("object already exists while trying to create %s", getObjectDescription(item))
        } else {
            return err
        }
    }

    return nil
}

func (r *IpaClusterReconciler) update(server freeipav1alpha1.Server) error {
    return nil
}

func (r *IpaClusterReconciler) delete(server freeipav1alpha1.Server) error {
    return nil
}

func (r *IpaClusterReconciler) buildStatus(cluster *freeipav1alpha1.IpaCluster, statuses map[string]freeipav1alpha1.ServerStatus) (error, bool) {
    return nil, false
}

func getObjectDescription(obj runtime.Object) string {
    kind := obj.GetObjectKind().GroupVersionKind().Kind
    unstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
    if err != nil {
        return "unknown"
    }
    return fmt.Sprintf("%s:%s/%s", unstr["metadata"].(map[string]interface{})["namespace"], kind, unstr["metadata"].(map[string]interface{})["name"])
}
