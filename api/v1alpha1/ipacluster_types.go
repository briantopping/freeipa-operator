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

package v1alpha1

import (
    //"github.com/briantopping/freeipa-operator/controllers"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const IpaClusterFinalizerName = "ipacluster.finalizers.coglative.com"

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
    UIDStart int32 `json:"uidStart,omitempty"`
	// The instantiation parameters for the nodes
	Servers []Server `json:"servers"`
}

type Server struct {
	// The FQDN of the server
	ServerName string `json:"serverName"`
	// The name of the secret for a node type, defaults to "ipa-server-secrets"
    // +optional
	SecretName string `json:"secretName,omitempty"`
	// Whether to create a DNS server / replica on this node, defaults to `false`
    // +optional
	DnsEnable bool `json:"dnsEnable,omitempty"`
	// Whether to create a CA server / replica on this node, defaults to `false`
    // +optional
	CaEnable bool `json:"caEnable,omitempty"`
	// Whether to create a NTP server / replica on this node, defaults to `false`
    // +optional
	NtpEnable bool `json:"ntpEnable,omitempty"`
	// Name of the storage class to use. Will try default storage class if omitted
    // +optional
	StorageClassName string `json:"storageClassName,omitempty"`
	// Size of the storage allocation
    // +optional
	Capacity string `json:"capacity,omitempty"`
	// The LB address of a node
    // +optional
	LbAddress string `json:"address,omitempty"`
	// The externalTrafficPolicy of the LoadBalancer Service
    // +optional
	ExternalTrafficPolicy string `json:"externalTrafficPolicy,omitempty"`
}

// IpaClusterStatus defines the observed state of IpaCluster
type IpaClusterStatus struct {
    ServerStatus map[string]ServerStatus `json:"serverStatus"`
}

// +kubebuilder:object:root=true

// IpaCluster is the Schema for the ipaclusters API
type IpaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   *IpaClusterSpec   `json:"spec,omitempty"`
	Status *IpaClusterStatus `json:"status,omitempty"`
}

func (c *IpaCluster) IsBeingDeleted() bool {
    return !c.ObjectMeta.DeletionTimestamp.IsZero()
}

func (c *IpaCluster) HasFinalizer(finalizerName string) bool {
    return containsString(c.ObjectMeta.Finalizers, finalizerName)
}

func (c *IpaCluster) AddFinalizer(finalizerName string) {
    c.ObjectMeta.Finalizers = append(c.ObjectMeta.Finalizers, finalizerName)
}

func (c *IpaCluster) RemoveFinalizer(finalizerName string) {
    c.ObjectMeta.Finalizers = removeString(c.ObjectMeta.Finalizers, finalizerName)
}

// +kubebuilder:object:root=true

// IpaClusterList contains a list of IpaCluster
type IpaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpaCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpaCluster{}, &IpaClusterList{})
}

func containsString(slice []string, s string) bool {
    for _, item := range slice {
        if item == s {
            return true
        }
    }
    return false
}

func removeString(slice []string, s string) (result []string) {
    for _, item := range slice {
        if item == s {
            continue
        }
        result = append(result, item)
    }
    return
}
