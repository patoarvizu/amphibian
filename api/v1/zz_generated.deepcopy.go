// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteConfig) DeepCopyInto(out *RemoteConfig) {
	*out = *in
	out.Workspaces = in.Workspaces
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteConfig.
func (in *RemoteConfig) DeepCopy() *RemoteConfig {
	if in == nil {
		return nil
	}
	out := new(RemoteConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3Config) DeepCopyInto(out *S3Config) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3Config.
func (in *S3Config) DeepCopy() *S3Config {
	if in == nil {
		return nil
	}
	out := new(S3Config)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Target) DeepCopyInto(out *Target) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Target.
func (in *Target) DeepCopy() *Target {
	if in == nil {
		return nil
	}
	out := new(Target)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformState) DeepCopyInto(out *TerraformState) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformState.
func (in *TerraformState) DeepCopy() *TerraformState {
	if in == nil {
		return nil
	}
	out := new(TerraformState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TerraformState) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformStateList) DeepCopyInto(out *TerraformStateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TerraformState, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformStateList.
func (in *TerraformStateList) DeepCopy() *TerraformStateList {
	if in == nil {
		return nil
	}
	out := new(TerraformStateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TerraformStateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformStateSpec) DeepCopyInto(out *TerraformStateSpec) {
	*out = *in
	out.RemoteConfig = in.RemoteConfig
	out.S3Config = in.S3Config
	out.Target = in.Target
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformStateSpec.
func (in *TerraformStateSpec) DeepCopy() *TerraformStateSpec {
	if in == nil {
		return nil
	}
	out := new(TerraformStateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformStateStatus) DeepCopyInto(out *TerraformStateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformStateStatus.
func (in *TerraformStateStatus) DeepCopy() *TerraformStateStatus {
	if in == nil {
		return nil
	}
	out := new(TerraformStateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkspacesConfig) DeepCopyInto(out *WorkspacesConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkspacesConfig.
func (in *WorkspacesConfig) DeepCopy() *WorkspacesConfig {
	if in == nil {
		return nil
	}
	out := new(WorkspacesConfig)
	in.DeepCopyInto(out)
	return out
}
