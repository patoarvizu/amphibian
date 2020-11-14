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
	"fmt"
	"os"
	"os/exec"
	"time"

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
		Value map[string]interface{} `json:"value"`
	} `json:"outputs"`
}

// +kubebuilder:rbac:groups=terraform.patoarvizu.dev,resources=terraformstates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=terraform.patoarvizu.dev,resources=terraformstates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=watch;list;create;get;update;patch

func (r *TerraformStateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("terraformstate", req.NamespacedName)

	baseDir := "/terraform"
	stateDir := fmt.Sprintf("%s/%s/%s", baseDir, req.Namespace, req.Name)
	err := os.MkdirAll(stateDir, 0777)
	if err != nil {
		return ctrl.Result{}, err
	}

	state := &terraformv1.TerraformState{}
	err = r.Get(ctx, req.NamespacedName, state)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

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
	case "consul":
		dataBody.SetAttributeValue("config", createConsulBackendBody(state.Spec.ConsulConfig))
	}
	dataFile, err := os.Create(fmt.Sprintf("%s/data.tf", stateDir))
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
	outputsFile, err := os.Create(fmt.Sprintf("%s/outputs.tf", stateDir))
	if err != nil {
		return ctrl.Result{}, err
	}
	defer outputsFile.Close()
	outputsFile.Write(outputs.Bytes())

	cmd := exec.Command("terraform", "apply", "-auto-approve")
	cmd.Dir = stateDir
	err = cmd.Run()
	if err != nil {
		return ctrl.Result{}, err
	}

	cmd = exec.Command("terraform", "output", "-json")
	cmd.Dir = stateDir
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
	for k, v := range tfOutputs.Outputs.Value {
		s, ok := v.(string)
		if ok {
			configMapData[k] = s
		} else {
			data, err := json.Marshal(v)
			if err == nil {
				configMapData[k] = fmt.Sprintf("%s", data)
			} else {
				r.Log.Info(fmt.Sprintf("Skipping field %s: %v", k, err))
			}
		}
	}
	configMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Namespace: state.ObjectMeta.Namespace, Name: state.Spec.Target.ConfigMapName}, configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			configMap.ObjectMeta.Namespace = state.ObjectMeta.Namespace
			configMap.ObjectMeta.Name = state.Spec.Target.ConfigMapName
			configMap.Data = configMapData
			err = ctrl.SetControllerReference(state, configMap, r.Scheme)
			if err != nil {
				return ctrl.Result{}, err
			}
			err = r.Create(ctx, configMap)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second * time.Duration(60)}, nil
		}
		return ctrl.Result{}, err
	}
	configMap.Data = configMapData
	err = r.Update(ctx, configMap)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * time.Duration(60)}, nil
}

func (r *TerraformStateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terraformv1.TerraformState{}).
		Complete(r)
}

func createRemoteBackendBody(config terraformv1.RemoteConfig) cty.Value {
	c := make(map[string]cty.Value)
	c["hostname"] = cty.StringVal(config.Hostname)
	c["organization"] = cty.StringVal(config.Organization)
	if len(config.Workspaces.Prefix) > 0 {
		c["workspaces"] = cty.ObjectVal(map[string]cty.Value{
			"prefix": cty.StringVal(config.Workspaces.Prefix),
		})
	} else {
		c["workspaces"] = cty.ObjectVal(map[string]cty.Value{
			"name": cty.StringVal(config.Workspaces.Name),
		})
	}
	return cty.ObjectVal(c)
}

func createS3BackendBody(config terraformv1.S3Config) cty.Value {
	c := make(map[string]cty.Value)
	c["bucket"] = cty.StringVal(config.Bucket)
	c["key"] = cty.StringVal(config.Key)
	if len(config.Region) > 0 {
		c["region"] = cty.StringVal(config.Region)
	}
	if len(config.AccessKey) > 0 {
		c["access_key"] = cty.StringVal(config.AccessKey)
	}
	if len(config.SecretKey) > 0 {
		c["secret_key"] = cty.StringVal(config.SecretKey)
	}
	if len(config.IAMEndpoint) > 0 {
		c["iam_endpoint"] = cty.StringVal(config.IAMEndpoint)
	}
	if config.MaxRetries > 0 {
		c["max_retries"] = cty.NumberIntVal(config.MaxRetries)
	}
	if len(config.Profile) > 0 {
		c["profile"] = cty.StringVal(config.Profile)
	}
	if len(config.SharedCredentialsFile) > 0 {
		c["shared_credentials_file"] = cty.StringVal(config.SharedCredentialsFile)
	}
	if len(config.STSEndpoint) > 0 {
		c["sts_endpoint"] = cty.StringVal(config.STSEndpoint)
	}
	if len(config.Token) > 0 {
		c["token"] = cty.StringVal(config.Token)
	}
	if config.AssumeRoleDurationSeconds > 0 {
		c["assume_role_duration_seconds"] = cty.NumberIntVal(config.AssumeRoleDurationSeconds)
	}
	if len(config.AssumeRolePolicy) > 0 {
		c["assume_role_policy"] = cty.StringVal(config.AssumeRolePolicy)
	}
	if len(config.AssumeRolePolicyARNs) > 0 {
		c["assume_role_policy_arns"] = cty.ListVal(createValueList(config.AssumeRolePolicyARNs))
	}
	if len(config.AssumeRoleTags) > 0 {
		c["assume_role_tags"] = cty.MapVal(createValueMap(config.AssumeRoleTags))
	}
	if len(config.AssumeRoleTransitiveTagKeys) > 0 {
		c["assume_role_transitive_tag_keys"] = cty.ListVal(createValueList(config.AssumeRoleTransitiveTagKeys))
	}
	if len(config.ExternalID) > 0 {
		c["external_id"] = cty.StringVal(config.ExternalID)
	}
	if len(config.RoleARN) > 0 {
		c["role_arn"] = cty.StringVal(config.RoleARN)
	}
	if len(config.SessionName) > 0 {
		c["session_name"] = cty.StringVal(config.SessionName)
	}
	if len(config.Endpoint) > 0 {
		c["endpoint"] = cty.StringVal(config.Endpoint)
	}
	if len(config.KMSKeyID) > 0 {
		c["kms_key_id"] = cty.StringVal(config.KMSKeyID)
	}
	if len(config.SSECustomerKey) > 0 {
		c["sse_customer_key"] = cty.StringVal(config.SSECustomerKey)
	}
	if len(config.WorkspaceKeyPrefix) > 0 {
		c["workspace_key_prefix"] = cty.StringVal(config.WorkspaceKeyPrefix)
	}
	c["skip_credentials_validation"] = cty.BoolVal(config.SkipCredentialsValidation)
	c["skip_region_validation"] = cty.BoolVal(config.SkipRegionValidation)
	c["skip_metadata_api_check"] = cty.BoolVal(config.SkipMetadataAPICheck)
	c["force_path_style"] = cty.BoolVal(config.ForcePathStyle)
	return cty.ObjectVal(c)
}

func createConsulBackendBody(config terraformv1.ConsulConfig) cty.Value {
	c := make(map[string]cty.Value)
	c["path"] = cty.StringVal(config.Path)
	if len(config.AccessToken) > 0 {
		c["access_token"] = cty.StringVal(config.AccessToken)
	}
	if len(config.Address) > 0 {
		c["address"] = cty.StringVal(config.Address)
	}
	if len(config.Scheme) > 0 {
		c["scheme"] = cty.StringVal(config.Scheme)
	}
	if len(config.Datacenter) > 0 {
		c["datacenter"] = cty.StringVal(config.Datacenter)
	}
	if len(config.HTTPAuth) > 0 {
		c["http_auth"] = cty.StringVal(config.HTTPAuth)
	}
	if len(config.CAFile) > 0 {
		c["ca_file"] = cty.StringVal(config.CAFile)
	}
	if len(config.CertFile) > 0 {
		c["cert_file"] = cty.StringVal(config.CertFile)
	}
	if len(config.KeyFile) > 0 {
		c["key_file"] = cty.StringVal(config.KeyFile)
	}
	return cty.ObjectVal(c)
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
