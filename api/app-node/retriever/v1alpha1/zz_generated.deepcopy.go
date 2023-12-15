//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonRetrieverConfig) DeepCopyInto(out *CommonRetrieverConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonRetrieverConfig.
func (in *CommonRetrieverConfig) DeepCopy() *CommonRetrieverConfig {
	if in == nil {
		return nil
	}
	out := new(CommonRetrieverConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Input) DeepCopyInto(out *Input) {
	*out = *in
	out.KnowledgeBaseRef = in.KnowledgeBaseRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Input.
func (in *Input) DeepCopy() *Input {
	if in == nil {
		return nil
	}
	out := new(Input)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KnowledgeBaseRetriever) DeepCopyInto(out *KnowledgeBaseRetriever) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KnowledgeBaseRetriever.
func (in *KnowledgeBaseRetriever) DeepCopy() *KnowledgeBaseRetriever {
	if in == nil {
		return nil
	}
	out := new(KnowledgeBaseRetriever)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KnowledgeBaseRetriever) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KnowledgeBaseRetrieverList) DeepCopyInto(out *KnowledgeBaseRetrieverList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KnowledgeBaseRetriever, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KnowledgeBaseRetrieverList.
func (in *KnowledgeBaseRetrieverList) DeepCopy() *KnowledgeBaseRetrieverList {
	if in == nil {
		return nil
	}
	out := new(KnowledgeBaseRetrieverList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KnowledgeBaseRetrieverList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KnowledgeBaseRetrieverSpec) DeepCopyInto(out *KnowledgeBaseRetrieverSpec) {
	*out = *in
	out.CommonSpec = in.CommonSpec
	out.Input = in.Input
	in.Output.DeepCopyInto(&out.Output)
	out.CommonRetrieverConfig = in.CommonRetrieverConfig
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KnowledgeBaseRetrieverSpec.
func (in *KnowledgeBaseRetrieverSpec) DeepCopy() *KnowledgeBaseRetrieverSpec {
	if in == nil {
		return nil
	}
	out := new(KnowledgeBaseRetrieverSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KnowledgeBaseRetrieverStatus) DeepCopyInto(out *KnowledgeBaseRetrieverStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KnowledgeBaseRetrieverStatus.
func (in *KnowledgeBaseRetrieverStatus) DeepCopy() *KnowledgeBaseRetrieverStatus {
	if in == nil {
		return nil
	}
	out := new(KnowledgeBaseRetrieverStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Output) DeepCopyInto(out *Output) {
	*out = *in
	in.CommonOrInPutOrOutputRef.DeepCopyInto(&out.CommonOrInPutOrOutputRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Output.
func (in *Output) DeepCopy() *Output {
	if in == nil {
		return nil
	}
	out := new(Output)
	in.DeepCopyInto(out)
	return out
}
