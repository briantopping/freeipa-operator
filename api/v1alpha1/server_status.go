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

package v1alpha1

import "github.com/briantopping/freeipa-operator/api"

type ServerStatus struct {
    State           api.ServerState
    GUID            *string // The ObjectMeta generation of the parent for identifying items to reconcile
    ResourceVersion string  `json:"resourceVersion,omitempty"`
}

//func (podStatus ServerStatus) Dispatch(reconciler api.StateReconciler) {
//    switch podStatus.State {
//    case api.StatusCreating:
//        reconciler.HandleCreating(podStatus)
//    case api.StatusReCreating:
//        reconciler.HandleReCreating(podStatus)
//    case api.StatusPeering:
//        reconciler.HandlePeering(podStatus)
//    case api.StatusRunning:
//        reconciler.HandleRunning(podStatus)
//    case api.StatusDegraded:
//        reconciler.HandleDegraded(podStatus)
//    case api.StatusStopping:
//        reconciler.HandleStopping(podStatus)
//    }
//
//}

func (podStatus ServerStatus) HandleCreating() {}
func (podStatus ServerStatus) HandleReCreating() {}
func (podStatus ServerStatus) HandlePeering()  {}
func (podStatus ServerStatus) HandleRunning()  {}
func (podStatus ServerStatus) HandleDegraded()  {}
func (podStatus ServerStatus) HandleStopping()  {}
