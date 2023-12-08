/*
Copyright 2023 KubeAGI.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	node "github.com/kubeagi/arcadia/api/app-node"
	"github.com/kubeagi/arcadia/api/base/v1alpha1"
)

// RetrievalQAChainSpec defines the desired state of RetrievalQAChain
type RetrievalQAChainSpec struct {
	v1alpha1.CommonSpec `json:",inline"`

	CommonChainConfig `json:",inline"`

	Input  RetrievalQAChainInput `json:"input"`
	Output Output                `json:"output"`
}

type RetrievalQAChainInput struct {
	LLMChainInput `json:",inline"`
	Retriever     node.RetrieverRef `json:"retriever"`
}

// RetrievalQAChainStatus defines the observed state of RetrievalQAChain
type RetrievalQAChainStatus struct {
	// ObservedGeneration is the last observed generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ConditionedStatus is the current status
	v1alpha1.ConditionedStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RetrievalQAChain is the Schema for the RetrievalQAChains API
type RetrievalQAChain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RetrievalQAChainSpec   `json:"spec,omitempty"`
	Status RetrievalQAChainStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RetrievalQAChainList contains a list of RetrievalQAChain
type RetrievalQAChainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RetrievalQAChain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RetrievalQAChain{}, &RetrievalQAChainList{})
}
