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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FileDistributionConfigSpec defines the desired state of FileDistributionConfig
type FileDistributionConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	SecretName         string `json:"secretName,omitempty"`
	FileName           string `json:"fileName,omitempty"`
	Destination        string `json:"destination,omitempty"`
	FilePermissions    string `json:"filepermissions,omitempty"`
	RescheduleInterval int    `json:"rescheduleInterval,omitempty"`
}

// FileDistributionConfigStatus defines the observed state of FileDistributionConfig
type FileDistributionConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FileDistributionConfig is the Schema for the filedistributionconfigs API
type FileDistributionConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FileDistributionConfigSpec   `json:"spec,omitempty"`
	Status FileDistributionConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FileDistributionConfigList contains a list of FileDistributionConfig
type FileDistributionConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FileDistributionConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FileDistributionConfig{}, &FileDistributionConfigList{})
}
