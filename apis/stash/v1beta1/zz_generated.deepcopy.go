// +build !ignore_autogenerated

/*
Copyright 2019 The Stash Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	apiv1 "kmodules.xyz/offshoot-api/api/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBatch) DeepCopyInto(out *BackupBatch) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBatch.
func (in *BackupBatch) DeepCopy() *BackupBatch {
	if in == nil {
		return nil
	}
	out := new(BackupBatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupBatch) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBatchList) DeepCopyInto(out *BackupBatchList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BackupBatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBatchList.
func (in *BackupBatchList) DeepCopy() *BackupBatchList {
	if in == nil {
		return nil
	}
	out := new(BackupBatchList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupBatchList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBatchSpec) DeepCopyInto(out *BackupBatchSpec) {
	*out = *in
	if in.BackupConfigurationTemplates != nil {
		in, out := &in.BackupConfigurationTemplates, &out.BackupConfigurationTemplates
		*out = make([]BackupConfigurationTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Repository = in.Repository
	in.RetentionPolicy.DeepCopyInto(&out.RetentionPolicy)
	if in.BackupHistoryLimit != nil {
		in, out := &in.BackupHistoryLimit, &out.BackupHistoryLimit
		*out = new(int32)
		**out = **in
	}
	if in.Hooks != nil {
		in, out := &in.Hooks, &out.Hooks
		*out = new(Hooks)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBatchSpec.
func (in *BackupBatchSpec) DeepCopy() *BackupBatchSpec {
	if in == nil {
		return nil
	}
	out := new(BackupBatchSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBlueprint) DeepCopyInto(out *BackupBlueprint) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBlueprint.
func (in *BackupBlueprint) DeepCopy() *BackupBlueprint {
	if in == nil {
		return nil
	}
	out := new(BackupBlueprint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupBlueprint) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBlueprintList) DeepCopyInto(out *BackupBlueprintList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BackupBlueprint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBlueprintList.
func (in *BackupBlueprintList) DeepCopy() *BackupBlueprintList {
	if in == nil {
		return nil
	}
	out := new(BackupBlueprintList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupBlueprintList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupBlueprintSpec) DeepCopyInto(out *BackupBlueprintSpec) {
	*out = *in
	in.RepositorySpec.DeepCopyInto(&out.RepositorySpec)
	in.Task.DeepCopyInto(&out.Task)
	in.RetentionPolicy.DeepCopyInto(&out.RetentionPolicy)
	in.RuntimeSettings.DeepCopyInto(&out.RuntimeSettings)
	in.TempDir.DeepCopyInto(&out.TempDir)
	if in.InterimVolumeTemplate != nil {
		in, out := &in.InterimVolumeTemplate, &out.InterimVolumeTemplate
		*out = new(v1.PersistentVolumeClaim)
		(*in).DeepCopyInto(*out)
	}
	if in.BackupHistoryLimit != nil {
		in, out := &in.BackupHistoryLimit, &out.BackupHistoryLimit
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupBlueprintSpec.
func (in *BackupBlueprintSpec) DeepCopy() *BackupBlueprintSpec {
	if in == nil {
		return nil
	}
	out := new(BackupBlueprintSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfiguration) DeepCopyInto(out *BackupConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfiguration.
func (in *BackupConfiguration) DeepCopy() *BackupConfiguration {
	if in == nil {
		return nil
	}
	out := new(BackupConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfigurationList) DeepCopyInto(out *BackupConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BackupConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfigurationList.
func (in *BackupConfigurationList) DeepCopy() *BackupConfigurationList {
	if in == nil {
		return nil
	}
	out := new(BackupConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfigurationSpec) DeepCopyInto(out *BackupConfigurationSpec) {
	*out = *in
	in.BackupConfigurationTemplateSpec.DeepCopyInto(&out.BackupConfigurationTemplateSpec)
	out.Repository = in.Repository
	in.RetentionPolicy.DeepCopyInto(&out.RetentionPolicy)
	if in.BackupHistoryLimit != nil {
		in, out := &in.BackupHistoryLimit, &out.BackupHistoryLimit
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfigurationSpec.
func (in *BackupConfigurationSpec) DeepCopy() *BackupConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(BackupConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfigurationTemplate) DeepCopyInto(out *BackupConfigurationTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfigurationTemplate.
func (in *BackupConfigurationTemplate) DeepCopy() *BackupConfigurationTemplate {
	if in == nil {
		return nil
	}
	out := new(BackupConfigurationTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupConfigurationTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfigurationTemplateSpec) DeepCopyInto(out *BackupConfigurationTemplateSpec) {
	*out = *in
	in.Task.DeepCopyInto(&out.Task)
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(BackupTarget)
		(*in).DeepCopyInto(*out)
	}
	in.RuntimeSettings.DeepCopyInto(&out.RuntimeSettings)
	in.TempDir.DeepCopyInto(&out.TempDir)
	if in.InterimVolumeTemplate != nil {
		in, out := &in.InterimVolumeTemplate, &out.InterimVolumeTemplate
		*out = new(v1.PersistentVolumeClaim)
		(*in).DeepCopyInto(*out)
	}
	if in.Hooks != nil {
		in, out := &in.Hooks, &out.Hooks
		*out = new(Hooks)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfigurationTemplateSpec.
func (in *BackupConfigurationTemplateSpec) DeepCopy() *BackupConfigurationTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(BackupConfigurationTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupInvokerRef) DeepCopyInto(out *BackupInvokerRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupInvokerRef.
func (in *BackupInvokerRef) DeepCopy() *BackupInvokerRef {
	if in == nil {
		return nil
	}
	out := new(BackupInvokerRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSession) DeepCopyInto(out *BackupSession) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSession.
func (in *BackupSession) DeepCopy() *BackupSession {
	if in == nil {
		return nil
	}
	out := new(BackupSession)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupSession) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSessionList) DeepCopyInto(out *BackupSessionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BackupSession, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSessionList.
func (in *BackupSessionList) DeepCopy() *BackupSessionList {
	if in == nil {
		return nil
	}
	out := new(BackupSessionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackupSessionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSessionSpec) DeepCopyInto(out *BackupSessionSpec) {
	*out = *in
	out.Invoker = in.Invoker
	if in.BackupConfiguration != nil {
		in, out := &in.BackupConfiguration, &out.BackupConfiguration
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSessionSpec.
func (in *BackupSessionSpec) DeepCopy() *BackupSessionSpec {
	if in == nil {
		return nil
	}
	out := new(BackupSessionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupSessionStatus) DeepCopyInto(out *BackupSessionStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = (*in).DeepCopy()
	}
	if in.TotalHosts != nil {
		in, out := &in.TotalHosts, &out.TotalHosts
		*out = new(int32)
		**out = **in
	}
	if in.Stats != nil {
		in, out := &in.Stats, &out.Stats
		*out = make([]HostBackupStats, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupSessionStatus.
func (in *BackupSessionStatus) DeepCopy() *BackupSessionStatus {
	if in == nil {
		return nil
	}
	out := new(BackupSessionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupTarget) DeepCopyInto(out *BackupTarget) {
	*out = *in
	out.Ref = in.Ref
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupTarget.
func (in *BackupTarget) DeepCopy() *BackupTarget {
	if in == nil {
		return nil
	}
	out := new(BackupTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EmptyDirSettings) DeepCopyInto(out *EmptyDirSettings) {
	*out = *in
	if in.SizeLimit != nil {
		in, out := &in.SizeLimit, &out.SizeLimit
		x := (*in).DeepCopy()
		*out = &x
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EmptyDirSettings.
func (in *EmptyDirSettings) DeepCopy() *EmptyDirSettings {
	if in == nil {
		return nil
	}
	out := new(EmptyDirSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileStats) DeepCopyInto(out *FileStats) {
	*out = *in
	if in.TotalFiles != nil {
		in, out := &in.TotalFiles, &out.TotalFiles
		*out = new(int)
		**out = **in
	}
	if in.NewFiles != nil {
		in, out := &in.NewFiles, &out.NewFiles
		*out = new(int)
		**out = **in
	}
	if in.ModifiedFiles != nil {
		in, out := &in.ModifiedFiles, &out.ModifiedFiles
		*out = new(int)
		**out = **in
	}
	if in.UnmodifiedFiles != nil {
		in, out := &in.UnmodifiedFiles, &out.UnmodifiedFiles
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileStats.
func (in *FileStats) DeepCopy() *FileStats {
	if in == nil {
		return nil
	}
	out := new(FileStats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Function) DeepCopyInto(out *Function) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Function.
func (in *Function) DeepCopy() *Function {
	if in == nil {
		return nil
	}
	out := new(Function)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Function) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionList) DeepCopyInto(out *FunctionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Function, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionList.
func (in *FunctionList) DeepCopy() *FunctionList {
	if in == nil {
		return nil
	}
	out := new(FunctionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FunctionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionRef) DeepCopyInto(out *FunctionRef) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionRef.
func (in *FunctionRef) DeepCopy() *FunctionRef {
	if in == nil {
		return nil
	}
	out := new(FunctionRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FunctionSpec) DeepCopyInto(out *FunctionSpec) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]v1.ContainerPort, len(*in))
		copy(*out, *in)
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeDevices != nil {
		in, out := &in.VolumeDevices, &out.VolumeDevices
		*out = make([]v1.VolumeDevice, len(*in))
		copy(*out, *in)
	}
	if in.RuntimeSettings != nil {
		in, out := &in.RuntimeSettings, &out.RuntimeSettings
		*out = new(apiv1.ContainerRuntimeSettings)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FunctionSpec.
func (in *FunctionSpec) DeepCopy() *FunctionSpec {
	if in == nil {
		return nil
	}
	out := new(FunctionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hooks) DeepCopyInto(out *Hooks) {
	*out = *in
	if in.PreBackup != nil {
		in, out := &in.PreBackup, &out.PreBackup
		*out = new(v1.Handler)
		(*in).DeepCopyInto(*out)
	}
	if in.PostBackup != nil {
		in, out := &in.PostBackup, &out.PostBackup
		*out = new(v1.Handler)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hooks.
func (in *Hooks) DeepCopy() *Hooks {
	if in == nil {
		return nil
	}
	out := new(Hooks)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostBackupStats) DeepCopyInto(out *HostBackupStats) {
	*out = *in
	if in.Snapshots != nil {
		in, out := &in.Snapshots, &out.Snapshots
		*out = make([]SnapshotStats, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostBackupStats.
func (in *HostBackupStats) DeepCopy() *HostBackupStats {
	if in == nil {
		return nil
	}
	out := new(HostBackupStats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostRestoreStats) DeepCopyInto(out *HostRestoreStats) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostRestoreStats.
func (in *HostRestoreStats) DeepCopy() *HostRestoreStats {
	if in == nil {
		return nil
	}
	out := new(HostRestoreStats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Param) DeepCopyInto(out *Param) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Param.
func (in *Param) DeepCopy() *Param {
	if in == nil {
		return nil
	}
	out := new(Param)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSession) DeepCopyInto(out *RestoreSession) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSession.
func (in *RestoreSession) DeepCopy() *RestoreSession {
	if in == nil {
		return nil
	}
	out := new(RestoreSession)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RestoreSession) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSessionList) DeepCopyInto(out *RestoreSessionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RestoreSession, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSessionList.
func (in *RestoreSessionList) DeepCopy() *RestoreSessionList {
	if in == nil {
		return nil
	}
	out := new(RestoreSessionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RestoreSessionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSessionSpec) DeepCopyInto(out *RestoreSessionSpec) {
	*out = *in
	out.Repository = in.Repository
	in.Task.DeepCopyInto(&out.Task)
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(RestoreTarget)
		(*in).DeepCopyInto(*out)
	}
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]Rule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.RuntimeSettings.DeepCopyInto(&out.RuntimeSettings)
	in.TempDir.DeepCopyInto(&out.TempDir)
	if in.InterimVolumeTemplate != nil {
		in, out := &in.InterimVolumeTemplate, &out.InterimVolumeTemplate
		*out = new(v1.PersistentVolumeClaim)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSessionSpec.
func (in *RestoreSessionSpec) DeepCopy() *RestoreSessionSpec {
	if in == nil {
		return nil
	}
	out := new(RestoreSessionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreSessionStatus) DeepCopyInto(out *RestoreSessionStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = (*in).DeepCopy()
	}
	if in.TotalHosts != nil {
		in, out := &in.TotalHosts, &out.TotalHosts
		*out = new(int32)
		**out = **in
	}
	if in.Stats != nil {
		in, out := &in.Stats, &out.Stats
		*out = make([]HostRestoreStats, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreSessionStatus.
func (in *RestoreSessionStatus) DeepCopy() *RestoreSessionStatus {
	if in == nil {
		return nil
	}
	out := new(RestoreSessionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreTarget) DeepCopyInto(out *RestoreTarget) {
	*out = *in
	out.Ref = in.Ref
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]v1.PersistentVolumeClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreTarget.
func (in *RestoreTarget) DeepCopy() *RestoreTarget {
	if in == nil {
		return nil
	}
	out := new(RestoreTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	if in.TargetHosts != nil {
		in, out := &in.TargetHosts, &out.TargetHosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Snapshots != nil {
		in, out := &in.Snapshots, &out.Snapshots
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SnapshotStats) DeepCopyInto(out *SnapshotStats) {
	*out = *in
	in.FileStats.DeepCopyInto(&out.FileStats)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SnapshotStats.
func (in *SnapshotStats) DeepCopy() *SnapshotStats {
	if in == nil {
		return nil
	}
	out := new(SnapshotStats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetRef) DeepCopyInto(out *TargetRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetRef.
func (in *TargetRef) DeepCopy() *TargetRef {
	if in == nil {
		return nil
	}
	out := new(TargetRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Task) DeepCopyInto(out *Task) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Task.
func (in *Task) DeepCopy() *Task {
	if in == nil {
		return nil
	}
	out := new(Task)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Task) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskList) DeepCopyInto(out *TaskList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Task, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskList.
func (in *TaskList) DeepCopy() *TaskList {
	if in == nil {
		return nil
	}
	out := new(TaskList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TaskList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRef) DeepCopyInto(out *TaskRef) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRef.
func (in *TaskRef) DeepCopy() *TaskRef {
	if in == nil {
		return nil
	}
	out := new(TaskRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskSpec) DeepCopyInto(out *TaskSpec) {
	*out = *in
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]FunctionRef, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskSpec.
func (in *TaskSpec) DeepCopy() *TaskSpec {
	if in == nil {
		return nil
	}
	out := new(TaskSpec)
	in.DeepCopyInto(out)
	return out
}
