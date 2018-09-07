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

package validating

import (
	"context"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validate-create-update-pods"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &QuxCreateUpdateHandler{})
}

// QuxCreateUpdateHandler validates Pods
type QuxCreateUpdateHandler struct {
	Client  client.Client
	Decoder types.Decoder
}

func (h *QuxCreateUpdateHandler) validatePodsFn(ctx context.Context, pod *corev1.Pod) (bool, string, error) {
	key := "example-admission-webhook"
	anno, found := pod.Annotations[key]
	switch {
	case !found:
		return found, fmt.Sprintf("failed to find annotation with key: %q", key), nil
	case found && anno == "foo":
		return found, "", nil
	case found && anno != "foo":
		return false,
			fmt.Sprintf("the value associate with key %q is expected to be %q, but got %q", "foo", "foo", anno), nil
	}
	return false, "", nil
}

// Implement admission.QuxCreateUpdateHandler so the controller can handle admission request.
var _ admission.Handler = &QuxCreateUpdateHandler{}

// QuxCreateUpdateHandler admits a pod iff a specific annotation exists.
func (h *QuxCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	pod := &corev1.Pod{}

	err := h.Decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.validatePodsFn(ctx, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

var _ inject.Client = &QuxCreateUpdateHandler{}

// InjectClient injects the client into the QuxCreateUpdateHandler
func (h *QuxCreateUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &QuxCreateUpdateHandler{}

// InjectDecoder injects the decoder into the QuxCreateUpdateHandler
func (h *QuxCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
