/*
Copyright 2018 The Kubernetes Authors.

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

package mutating

import (
	"context"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "mutate-update-pods"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &FooUpdateHandler{})
}

// FooUpdateHandler annotates Pods
type FooUpdateHandler struct {
	Client  client.Client
	Decoder types.Decoder
}

// mutatePodsFn add an annotation to the given pod
func (h *FooUpdateHandler) mutatePodsFn(ctx context.Context, pod *corev1.Pod) error {
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["example-admission-webhook"] = "foo"
	return nil
}

// Implement admission.FooUpdateHandler so the controller can handle admission request.
var _ admission.Handler = &FooUpdateHandler{}

// FooUpdateHandler adds an annotation to every incoming pods.
func (h *FooUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	pod := &corev1.Pod{}

	err := h.Decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := pod.DeepCopy()

	err = h.mutatePodsFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(pod, copy)
}

var _ inject.Client = &FooUpdateHandler{}

// InjectClient injects the client into the FooUpdateHandler
func (h *FooUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &FooUpdateHandler{}

// InjectDecoder injects the decoder into the FooUpdateHandler
func (h *FooUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
