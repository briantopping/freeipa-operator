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

package api

import "github.com/briantopping/freeipa-operator/api/v1alpha1"

/*
 * Versioned behaviors for reconciling configuration objects. Theory is to abastractly
 * manage different versions with the same underlying code, this provides a "double dispatch"
 * pattern for proper handling with static typing at the compiler level.
 */
type ServerState string

const (
    StatusCreating   ServerState = "creating"
    StatusReCreating ServerState = "re-creating"
    StatusPeering    ServerState = "peering"
    StatusRunning    ServerState = "running"
    StatusDegraded   ServerState = "degraded"
    StatusStopping   ServerState = "stopping"
    StatusError      ServerState = "error"
)

// +kubebuilder:object:generate=false
type StateReconciler interface {
    HandleCreating()
    HandleReCreating()
    HandlePeering()
    HandleRunning()
    HandleDegraded()
    HandleStopping()
}

