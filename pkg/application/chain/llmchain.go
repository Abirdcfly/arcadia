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

package chain

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	langchaingoschema "github.com/tmc/langchaingo/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/kubeagi/arcadia/api/app-node/chain/v1alpha1"
	"github.com/kubeagi/arcadia/pkg/application/base"
)

type LLMChain struct {
	chains.LLMChain
	base.BaseNode
}

func NewLLMChain(baseNode base.BaseNode) *LLMChain {
	return &LLMChain{
		chains.LLMChain{},
		baseNode,
	}
}

func (l *LLMChain) Run(ctx context.Context, cli dynamic.Interface, args map[string]any) (map[string]any, error) {
	v1, ok := args["llm"]
	if !ok {
		return args, errors.New("no llm")
	}
	llm, ok := v1.(llms.LanguageModel)
	if !ok {
		return args, errors.New("llm not llms.LanguageModel")
	}
	v2, ok := args["prompt"]
	if !ok {
		return args, errors.New("no prompt")
	}
	prompt, ok := v2.(prompts.FormatPrompter)
	if !ok {
		return args, errors.New("prompt not prompts.FormatPrompter")
	}
	v3, ok := args["_history"]
	if !ok {
		return args, errors.New("no history")
	}
	history, ok := v3.(langchaingoschema.ChatMessageHistory)
	if !ok {
		return args, errors.New("history not memory.ChatMessageHistory")
	}

	ns := base.GetAppNamespace(ctx)
	instance := &v1alpha1.LLMChain{}
	obj, err := cli.Resource(schema.GroupVersionResource{Group: v1alpha1.GroupVersion.Group, Version: v1alpha1.GroupVersion.Version, Resource: "llmchains"}).
		Namespace(l.Ref.GetNamespace(ns)).Get(ctx, l.Ref.Name, metav1.GetOptions{})
	if err != nil {
		return args, fmt.Errorf("cant find the chain in cluster: %w", err)
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), instance)
	if err != nil {
		return args, err
	}
	options := getChainOptions(instance.Spec.CommonChainConfig)

	chain := chains.NewLLMChain(llm, prompt)
	chain.Memory = getMemory(llm, instance.Spec.Memory, history)
	l.LLMChain = *chain
	var out string
	if needStream, ok := args["_need_stream"].(bool); ok && needStream {
		options = append(options, chains.WithStreamingFunc(stream(args)))
		out, err = chains.Predict(ctx, l.LLMChain, args, options...)
	} else {
		if len(options) > 0 {
			out, err = chains.Predict(ctx, l.LLMChain, args, options...)
		} else {
			out, err = chains.Predict(ctx, l.LLMChain, args)
		}
	}
	if err == nil {
		args["_answer"] = out
	}
	return args, err
}
