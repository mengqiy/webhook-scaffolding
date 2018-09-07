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
	"fmt"
	"net/http"

	crewv1alpha1 "github.com/mengqiy/webhook-scaffolding/pkg/apis/crew/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

// Handler annotates Pods
type Handler struct {
	Client  client.Client
	Decoder types.Decoder
}

// mutateFirstMateFn add an annotation to the given pod
func mutateFirstMateFn(ctx context.Context, fm *crewv1alpha1.Firstmate) error {
	v, ok := ctx.Value(admission.StringKey("foo")).(string)
	if !ok {
		return fmt.Errorf("the value associated with %v is expected to be a string", "foo")
	}
	fm.Spec.Foo = v
	return nil
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &Handler{}

// Handler changes a field in a CR.
func (a *Handler) Handle(ctx context.Context, req types.Request) types.Response {
	fm := &crewv1alpha1.Firstmate{}

	err := a.Decoder.Decode(req, fm)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := fm.DeepCopy()

	err = mutateFirstMateFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(fm, copy)
}

var _ inject.Client = &Handler{}

// InjectClient injects the client into the Handler
func (h *Handler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &Handler{}

// InjectDecoder injects the decoder into the Handler
func (h *Handler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
