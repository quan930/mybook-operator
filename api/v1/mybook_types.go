/*
Copyright 2022.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyBookSpec defines the desired state of MyBook
type MyBookSpec struct {
	//UUID book's UUID
	//+kubebuilder:validation:len=36
	UUID string `json:"uuid"`

	//Name book's Name
	//+kubebuilder:validation:max=32
	Name string `json:"name"`

	//Price book's Price
	//+kubebuilder:validation:Minimum=0
	Price int `json:"price"`
}

// MyBookStatus defines the observed state of MyBook
type MyBookStatus struct {
	//History update history
	History []MyBookSpec `json:"history"`
}

//+kubebuilder:object:root=true

// MyBook is the Schema for the mybooks API
//+kubebuilder:subresource:status
type MyBook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyBookSpec   `json:"spec,omitempty"`
	Status MyBookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyBookList contains a list of MyBook
type MyBookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyBook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyBook{}, &MyBookList{})
}
