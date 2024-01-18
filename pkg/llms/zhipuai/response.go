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
	"encoding/json"

	"github.com/kubeagi/arcadia/pkg/llms"
)

type EmbeddingResponse struct {
	Code    int            `json:"code"`
	Data    *EmbeddingData `json:"data"`
	Msg     string         `json:"msg"`
	Success bool           `json:"success"`
}

func (embeddingResp *EmbeddingResponse) Unmarshall(bytes []byte) error {
	return json.Unmarshal(embeddingResp.Bytes(), embeddingResp)
}

func (embeddingResp *EmbeddingResponse) Type() llms.LLMType {
	return llms.ZhiPuAI
}

func (embeddingResp *EmbeddingResponse) Bytes() []byte {
	bytes, err := json.Marshal(embeddingResp)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (embeddingResp *EmbeddingResponse) String() string {
	return string(embeddingResp.Bytes())
}

type Response struct {
	Code    int    `json:"code"`
	Data    *Data  `json:"data"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

func (response *Response) Unmarshal(bytes []byte) error {
	return json.Unmarshal(response.Bytes(), response)
}

func (response *Response) Type() llms.LLMType {
	return llms.ZhiPuAI
}

func (response *Response) Bytes() []byte {
	bytes, err := json.Marshal(response)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (response *Response) String() string {
	return string(response.Bytes())
}

type Data struct {
	// for async
	RequestID string `json:"request_id,omitempty"`
	TaskID    string `json:"id,omitempty"`
	// The request creation time, a Unix timestamp in seconds.
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices,omitempty"`
	Usage   Usage    `json:"usage,omitempty"`
	// for async
	TaskStatus string `json:"task_status,omitempty"`
}

type EmbeddingData struct {
	RequestID  string `json:"request_id,omitempty"`
	TaskID     string `json:"task_id,omitempty"`
	TaskStatus string `json:"task_status,omitempty"`
	Usage      Usage  `json:"usage,omitempty"`

	Embedding []float32 `json:"embedding,omitempty"` // Vectorized texts, length 1024.
}

type EmbeddingText struct {
	Prompt    string `json:"prompt,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type Usage struct {
	TotalTokens      int `json:"total_tokens,omitempty"`
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
}

type Choice struct {
	Index int `json:"index"`
	// The reason for the termination of the model's reasoning.
	// `stop` represents the natural end of reasoning or a trigger stop word.
	// `tool_calls` represents the model hit function.
	// `length` represents the maximum length of tokens reached.
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

const (
	CodeConcurrencyHigh = 1302 // 您当前使用该 API 的并发数过高，请降低并发，或联系客服增加限额
	CodefrequencyHigh   = 1303 // 您当前使用该 API 的频率过高，请降低频率，或联系客服增加限额
	CodeTimesHigh       = 1305 // 当前 API 请求过多，请稍后重试
)
