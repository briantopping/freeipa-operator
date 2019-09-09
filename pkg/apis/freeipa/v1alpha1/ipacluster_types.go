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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IpaClusterSpec defines the desired state of IpaCluster
type IpaClusterSpec struct {
	// The Kerberos realm name as in "EXAMPLE.COM", required
	RealmName string `json:"realmName"`
	// The directory service root name as in "example.com", required
	DomainName string `json:"domainName"`
	// A string list of DNS forwarders for name resolution, defaults to no forwarders
	// +optional
	DNSForwarders []string `json:"dnsForwarders,omitempty"`
	// An integer for the start of the UID numbering range, immutable after cluster instantiation, default is defined by FreeIPA
	// +optional
	UIDStart int `json:"uidStart,omitempty"`
	// The instantiation parameters for the nodes
	Servers []Server `json:"servers"`
}

type Server struct {
	// The FQDN of the server
	ServerName string `json:"serverName"`
	// The name of the secret for a node type, defaults to "ipa-server-secrets"
	SecretName string `json:"secretName,omitempty"`
	// Whether to create a DNS server / replica on this node, defaults to `false`
	DnsEnable bool `json:"dnsEnable,omitempty"`
	// Whether to create a CA server / replica on this node, defaults to `false`
	CaEnable bool `json:"caEnable,omitempty"`
	// Whether to create a NTP server / replica on this node, defaults to `false`
	NtpEnable bool `json:"ntpEnable,omitempty"`
	// Name of the storage class to use. Will try default storage class if omitted
	StorageClassName string `json:"storageClassName,omitempty"`
	// Size of the storage allocation
	Capacity string `json:"capacity,omitempty"`
	// The LB address of a node
	LbAddress string `json:"address,omitempty"`
	// The externalTrafficPolicy of the LoadBalancer Service
	ExternalTrafficPolicy string `json:"externalTrafficPolicy,omitempty"`
}

// IpaClusterStatus defines the observed state of IpaCluster
type IpaClusterStatus struct {
	// Human-readable status of the controller
	Status string `json:"status"`
	// Quantity of persistent volumes that are currently generated
	PvQuantity int `json:"pvQuantity"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaCluster is the Schema for the ipaclusters API
// +k8s:openapi-gen=true
type IpaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpaClusterSpec   `json:"spec,omitempty"`
	Status IpaClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaClusterList contains a list of IpaCluster
type IpaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpaCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpaCluster{}, &IpaClusterList{})
}
