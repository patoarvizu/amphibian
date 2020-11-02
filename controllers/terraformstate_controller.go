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
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/go-logr/logr"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	terraformv1 "github.com/patoarvizu/amphibian/api/v1"
	"github.com/zclconf/go-cty/cty"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TerraformStateReconciler reconciles a TerraformState object
type TerraformStateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type terraformOutputs struct {
	Outputs struct {
		Value map[string]string `json:"value"`
	} `json:"outputs"`
}

// +kubebuilder:rbac:groups=terraform.patoarvizu.dev,resources=terraformstates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=terraform.patoarvizu.dev,resources=terraformstates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=watch;list;create;get;update;patch

func (r *TerraformStateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("terraformstate", req.NamespacedName)

	state := &terraformv1.TerraformState{}
	err := r.Get(ctx, req.NamespacedName, state)

	if state.Spec.Type == "remote" {
		f, err := os.Create("/terraform/.terraformrc")
		if err != nil {
			return ctrl.Result{}, err
		}
		defer f.Close()
		terraformConfig := hclwrite.NewEmptyFile()
		body := terraformConfig.Body()
		credentialsBlock := body.AppendNewBlock("credentials", []string{"app.terraform.io"})
		credentialsBody := credentialsBlock.Body()
		credentialsBody.SetAttributeValue("token", cty.StringVal(os.Getenv("TERRAFORM_CLOUD_TOKEN")))
		f.Write(terraformConfig.Bytes())
	}

	data := hclwrite.NewEmptyFile()
	rootBody := data.Body()
	dataBlock := rootBody.AppendNewBlock("data", []string{"terraform_remote_state", "remote"})
	dataBody := dataBlock.Body()
	dataBody.SetAttributeValue("backend", cty.StringVal(state.Spec.Type))
	switch state.Spec.Type {
	case "remote":
		dataBody.SetAttributeValue("config", createRemoteBackendBody(state.Spec.RemoteConfig))
	case "s3":
		dataBody.SetAttributeValue("config", createS3BackendBody(state.Spec.S3Config))
	}
	dataFile, err := os.Create("/terraform/data.tf")
	if err != nil {
		return ctrl.Result{}, err
	}
	defer dataFile.Close()
	dataFile.Write(data.Bytes())

	outputs := hclwrite.NewEmptyFile()
	outputsRootBody := outputs.Body()
	outputBlock := outputsRootBody.AppendNewBlock("output", []string{"outputs"})
	outputBody := outputBlock.Body()
	outputBody.SetAttributeTraversal("value", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "data",
		},
		hcl.TraverseAttr{
			Name: "terraform_remote_state",
		},
		hcl.TraverseAttr{
			Name: "remote",
		},
		hcl.TraverseAttr{
			Name: "outputs",
		},
	})
	outputsFile, err := os.Create("/terraform/outputs.tf")
	if err != nil {
		return ctrl.Result{}, err
	}
	defer outputsFile.Close()
	outputsFile.Write(outputs.Bytes())

	cmd := exec.Command("terraform", "apply", "-auto-approve")
	cmd.Dir = "/terraform"
	err = cmd.Run()
	if err != nil {
		return ctrl.Result{}, err
	}

	cmd = exec.Command("terraform", "output", "-json")
	cmd.Dir = "/terraform"
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return ctrl.Result{}, err
	}

	var tfOutputs terraformOutputs
	err = json.Unmarshal(out.Bytes(), &tfOutputs)
	if err != nil {
		return ctrl.Result{}, err
	}

	configMapData := make(map[string]string)
	configMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Namespace: state.Spec.Target.NamespaceName, Name: state.Spec.Target.ConfigMapName}, configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			configMap.ObjectMeta.Namespace = state.Spec.Target.NamespaceName
			configMap.ObjectMeta.Name = state.Spec.Target.ConfigMapName
			for k, v := range tfOutputs.Outputs.Value {
				configMapData[k] = v
			}
			configMap.Data = configMapData
			err = r.Create(ctx, configMap)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, err
	}
	for k, v := range tfOutputs.Outputs.Value {
		configMapData[k] = v
	}
	configMap.Data = configMapData
	err = r.Update(ctx, configMap)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TerraformStateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terraformv1.TerraformState{}).
		Complete(r)
}

func createRemoteBackendBody(config terraformv1.RemoteConfig) cty.Value {
	configMap := make(map[string]cty.Value)
	configMap["hostname"] = cty.StringVal(config.Hostname)
	configMap["organization"] = cty.StringVal(config.Organization)
	if len(config.Workspaces.Prefix) > 0 {
		configMap["workspaces"] = cty.ObjectVal(map[string]cty.Value{
			"prefix": cty.StringVal(config.Workspaces.Prefix),
		})
	} else {
		configMap["workspaces"] = cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal(config.Workspaces.Name),
		})
	}
	return cty.ObjectVal(configMap)
}

func createS3BackendBody(config terraformv1.S3Config) cty.Value {
	configMap := make(map[string]cty.Value)
	configMap["bucket"] = cty.StringVal(config.Bucket)
	configMap["key"] = cty.StringVal(config.Key)
	if len(config.Region) > 0 {
		configMap["region"] = cty.StringVal(config.Region)
	}
	if len(config.AccessKey) > 0 {
		configMap["access_key"] = cty.StringVal(config.AccessKey)
	}
	if len(config.SecretKey) > 0 {
		configMap["secret_key"] = cty.StringVal(config.SecretKey)
	}
	if len(config.IAMEndpoint) > 0 {
		configMap["iam_endpoint"] = cty.StringVal(config.IAMEndpoint)
	}
	if config.MaxRetries > 0 {
		configMap["max_retries"] = cty.NumberIntVal(config.MaxRetries)
	}
	if len(config.Profile) > 0 {
		configMap["profile"] = cty.StringVal(config.Profile)
	}
	if len(config.SharedCredentialsFile) > 0 {
		configMap["shared_credentials_file"] = cty.StringVal(config.SharedCredentialsFile)
	}
	if len(config.STSEndpoint) > 0 {
		configMap["sts_endpoint"] = cty.StringVal(config.STSEndpoint)
	}
	if len(config.Token) > 0 {
		configMap["token"] = cty.StringVal(config.Token)
	}
	if config.AssumeRoleDurationSeconds > 0 {
		configMap["assume_role_duration_seconds"] = cty.NumberIntVal(config.AssumeRoleDurationSeconds)
	}
	if len(config.AssumeRolePolicy) > 0 {
		configMap["assume_role_policy"] = cty.StringVal(config.AssumeRolePolicy)
	}
	if len(config.AssumeRolePolicyARNs) > 0 {
		configMap["assume_role_policy_arns"] = cty.ListVal(createValueList(config.AssumeRolePolicyARNs))
	}
	if len(config.AssumeRoleTags) > 0 {
		configMap["assume_role_tags"] = cty.MapVal(createValueMap(config.AssumeRoleTags))
	}
	if len(config.AssumeRoleTransitiveTagKeys) > 0 {
		configMap["assume_role_transitive_tag_keys"] = cty.ListVal(createValueList(config.AssumeRoleTransitiveTagKeys))
	}
	if len(config.ExternalID) > 0 {
		configMap["external_id"] = cty.StringVal(config.ExternalID)
	}
	if len(config.RoleARN) > 0 {
		configMap["role_arn"] = cty.StringVal(config.RoleARN)
	}
	if len(config.SessionName) > 0 {
		configMap["session_name"] = cty.StringVal(config.SessionName)
	}
	if len(config.Endpoint) > 0 {
		configMap["endpoint"] = cty.StringVal(config.Endpoint)
	}
	if len(config.KMSKeyID) > 0 {
		configMap["kms_key_id"] = cty.StringVal(config.KMSKeyID)
	}
	if len(config.SSECustomerKey) > 0 {
		configMap["sse_customer_key"] = cty.StringVal(config.SSECustomerKey)
	}
	if len(config.WorkspaceKeyPrefix) > 0 {
		configMap["workspace_key_prefix"] = cty.StringVal(config.WorkspaceKeyPrefix)
	}
	configMap["skip_credentials_validation"] = cty.BoolVal(config.SkipCredentialsValidation)
	configMap["skip_region_validation"] = cty.BoolVal(config.SkipRegionValidation)
	configMap["skip_metadata_api_check"] = cty.BoolVal(config.SkipMetadataAPICheck)
	configMap["force_path_style"] = cty.BoolVal(config.ForcePathStyle)
	return cty.ObjectVal(configMap)
}

func createValueList(l []string) []cty.Value {
	valueList := []cty.Value{}
	for _, v := range l {
		valueList = append(valueList, cty.StringVal(v))
	}
	return valueList
}

func createValueMap(m map[string]string) map[string]cty.Value {
	valueMap := make(map[string]cty.Value)
	for k, v := range m {
		valueMap[k] = cty.StringVal(v)
	}
	return valueMap
}
