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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkspacesConfig struct {
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

type RemoteConfig struct {
	Hostname     string           `json:"hostname"`
	Organization string           `json:"organization"`
	Token        string           `json:"token,omitempty"`
	Workspaces   WorkspacesConfig `json:"workspaces"`
}

type Target struct {
	NamespaceName string `json:"namespace,omitempty"`
	ConfigMapName string `json:"configMapName,omitempty"`
}

type S3Config struct {
	Bucket                      string            `json:"bucket"`
	Key                         string            `json:"key"`
	Region                      string            `json:"region,omitempty"`
	AccessKey                   string            `json:"accessKey,omitempty"`
	SecretKey                   string            `json:"secretKey,omitempty"`
	IAMEndpoint                 string            `json:"iamEndpoint,omitempty"`
	MaxRetries                  int64             `json:"maxRetries,omitempty"`
	Profile                     string            `json:"profile,omitempty"`
	SharedCredentialsFile       string            `json:"sharedCredentialsFile,omitempty"`
	SkipCredentialsValidation   bool              `json:"skipCredentialsValidation,omitempty"`
	SkipRegionValidation        bool              `json:"skipRegionValidation,omitempty"`
	SkipMetadataAPICheck        bool              `json:"skipMetadataAPICheck,omitempty"`
	STSEndpoint                 string            `json:"stsEndpoint,omitempty"`
	Token                       string            `json:"token,omitempty"`
	AssumeRoleDurationSeconds   int64             `json:"assumeRoleDurationSeconds,omitempty"`
	AssumeRolePolicy            string            `json:"assumeRolePolicy,omitempty"`
	AssumeRolePolicyARNs        []string          `json:"assumeRolePolicyARNs,omitempty"`
	AssumeRoleTags              map[string]string `json:"assumeRoleTags,omitempty"`
	AssumeRoleTransitiveTagKeys []string          `json:"assumeRoleTransitiveTagKeys,omitempty"`
	ExternalID                  string            `json:"externalID,omitempty"`
	RoleARN                     string            `json:"roleARN,omitempty"`
	SessionName                 string            `json:"sessionName,omitempty"`
	Endpoint                    string            `json:"endpoint,omitempty"`
	ForcePathStyle              bool              `json:"forcePathStyle,omitempty"`
	KMSKeyID                    string            `json:"kmsKeyID,omitempty"`
	SSECustomerKey              string            `json:"sseCustomerKey,omitempty"`
	WorkspaceKeyPrefix          string            `json:"workspaceKeyPrefix,omitempty"`
}

type TerraformStateSpec struct {
	Type         string       `json:"type"`
	RemoteConfig RemoteConfig `json:"remoteConfig,omitempty"`
	S3Config     S3Config     `json:"s3Config,omitempty"`
	Target       Target       `json:"target,omitempty"`
}

type TerraformStateStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=terraformstates,scope=Namespaced,shortName=tfs
type TerraformState struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TerraformStateSpec   `json:"spec,omitempty"`
	Status TerraformStateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type TerraformStateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TerraformState `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TerraformState{}, &TerraformStateList{})
}
