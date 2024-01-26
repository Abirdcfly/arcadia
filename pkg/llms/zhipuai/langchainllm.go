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

package zhipuai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/tmc/langchaingo/callbacks"
	langchainllm "github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"k8s.io/klog/v2"
)

var (
	ErrEmptyResponse = errors.New("no response")
	ErrEmptyPrompt   = errors.New("empty prompt")
)

var (
	_ langchainllm.Model = (*ZhiPuAILLM)(nil)
)

type options struct {
	retryTimes int
}

type Option func(*options)

func WithRetryTimes(retryTimes int) Option {
	return func(o *options) {
		o.retryTimes = retryTimes
	}
}

type ZhiPuAILLM struct {
	c                *ZhiPuAI
	options          *options
	CallbacksHandler callbacks.Handler
}

func NewZhiPuAILLM(apiKey string, opts ...Option) *ZhiPuAILLM {
	z := &ZhiPuAILLM{
		c: NewZhiPuAI(apiKey),
		options: &options{
			// 2 times by default
			retryTimes: 2,
		},
	}
	for _, opt := range opts {
		opt(z.options)
	}
	return z
}

func (z *ZhiPuAILLM) GetNumTokens(text string) int {
	return langchainllm.CountTokens("gpt2", text)
}

func (z *ZhiPuAILLM) Call(ctx context.Context, prompt string, options ...langchainllm.CallOption) (string, error) {
	return langchainllm.GenerateFromSinglePrompt(ctx, z, prompt, options...)
}

func (z *ZhiPuAILLM) GenerateContent(ctx context.Context, messages []langchainllm.MessageContent, options ...langchainllm.CallOption) (*langchainllm.ContentResponse, error) {
	if z.CallbacksHandler != nil {
		z.CallbacksHandler.HandleLLMGenerateContentStart(ctx, messages)
	}
	opts := langchainllm.CallOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	params := DefaultModelParams()
	if opts.TopP > 0 && opts.TopP < 1 {
		params.TopP = float32(opts.TopP)
	}
	if opts.Temperature > 0 && opts.Temperature < 1 {
		params.Temperature = float32(opts.Temperature)
	}
	if opts.Model != "" {
		params.Model = opts.Model
	}
	if len(messages) == 0 {
		return nil, ErrEmptyPrompt
	}
	for _, mc := range messages {
		msgByte, _ := json.Marshal(mc.Parts)
		msg := string(msgByte)
		switch mc.Role {
		case schema.ChatMessageTypeAI:
			params.Prompt = append(params.Prompt, Prompt{Role: Assistant, Content: msg})
		case schema.ChatMessageTypeHuman:
			params.Prompt = append(params.Prompt, Prompt{Role: User, Content: msg})
		default:
			return nil, fmt.Errorf("role %v not supported", mc.Role)
		}
	}
	needStream := opts.StreamingFunc != nil
	if needStream {
		res := bytes.NewBuffer(nil)
		err := z.c.SSEInvoke(params, func(event *sse.Event) {
			if string(event.Event) == "finish" {
				return
			}
			_, _ = res.Write(event.Data)
			_ = opts.StreamingFunc(ctx, event.Data)
		})
		if err != nil {
			return nil, err
		}
		return &langchainllm.ContentResponse{
			Choices: []*langchainllm.ContentChoice{
				{
					Content: res.String(),
				},
			},
		}, nil
	}
	var resp *Response
	var err error
	i := 0
	for {
		i++
		resp, err = z.c.Invoke(params)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, ErrEmptyResponse
		}
		if resp.Data == nil {
			klog.Errorf("empty response: msg:%s code:%d\n", resp.Msg, resp.Code)
			if i <= z.options.retryTimes {
				r := rand.Intn(5)
				klog.Infof("retry[%d], sleep %d seconds, then recall...\n", i, r)
				time.Sleep(time.Duration(r) * time.Second)
				continue
			}
			return nil, ErrEmptyResponse
		}
		if len(resp.Data.Choices) == 0 {
			return nil, ErrEmptyResponse
		}
		break
	}
	choices := make([]*langchainllm.ContentChoice, 1)
	generationInfo := make(map[string]any, reflect.ValueOf(resp.Data.Usage).NumField())
	generationInfo["TotalTokens"] = resp.Data.Usage.TotalTokens
	var s string
	if err := json.Unmarshal([]byte(resp.Data.Choices[0].Content), &s); err != nil {
		return nil, err
	}
	choices[0].Content = s
	response := &langchainllm.ContentResponse{Choices: choices}
	if z.CallbacksHandler != nil {
		z.CallbacksHandler.HandleLLMGenerateContentEnd(ctx, response)
	}
	return response, nil
}
