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
	// +kubebuilder:validation:Enum={"configmap","secret"}
	Type string `json:"type"`
	Name string `json:"name"`
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

type ConsulConfig struct {
	Path        string `json:"path"`
	AccessToken string `json:"accessToken,omitempty"`
	Address     string `json:"address,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
	Datacenter  string `json:"datacenter,omitempty"`
	HTTPAuth    string `json:"httpAuth,omitempty"`
	CAFile      string `json:"caFile,omitempty"`
	CertFile    string `json:"certFile,omitempty"`
	KeyFile     string `json:"keyFile,omitempty"`
}

type KubernetesConfig struct {
	SecretSuffix    string `json:"secretSuffix"`
	Namespace       string `json:"namespace,omitempty"`
	InClusterConfig bool   `json:"inClusterConfig,omitempty"`
	Host            string `json:"host,omitempty"`
	Insecure        bool   `json:"insecure,omitempty"`
	ConfigPath      string `json:"configPath,omitempty"`
}

type GCSConfig struct {
	Bucket                    string `json:"bucket"`
	Credentials               string `json:"credentials,omitempty"`
	ImpersonateServiceAccount string `json:"impersonateServiceAccount,omitempty"`
	AccessToken               string `json:"accessToken,omitempty"`
	Prefix                    string `json:"prefix,omitempty"`
}

type TerraformStateSpec struct {
	Type             string           `json:"type"`
	RemoteConfig     RemoteConfig     `json:"remoteConfig,omitempty"`
	S3Config         S3Config         `json:"s3Config,omitempty"`
	ConsulConfig     ConsulConfig     `json:"consulConfig,omitempty"`
	KubernetesConfig KubernetesConfig `json:"kubernetesConfig,omitempty"`
	GCSConfig        GCSConfig        `json:"gcsConfig,omitempty"`
	Target           Target           `json:"target"`
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
