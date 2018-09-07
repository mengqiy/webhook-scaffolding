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

	creaturesv1alpha1 "github.com/mengqiy/example-crd-apis/pkg/apis/creatures/v1alpha1"
	crewv1alpha1 "github.com/mengqiy/webhook-scaffolding/pkg/apis/crew/v1alpha1"
	apitypes "k8s.io/apimachinery/pkg/types"
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

// mutateKrakenFn add an annotation to the given pod
func (h *Handler) mutateKrakenFn(ctx context.Context, k *creaturesv1alpha1.Kraken) error {
	firstmate := &crewv1alpha1.Firstmate{}
	err := h.Client.Get(context.Background(), apitypes.NamespacedName{Namespace: "default", Name: "firstmate-sample"}, firstmate)
	if err != nil {
		return err
	}

	v := firstmate.Spec.Foo
	anno := k.GetAnnotations()
	if anno == nil {
		anno = map[string]string{}
	}
	anno["foo"] = v + v
	k.SetAnnotations(anno)
	return nil
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &Handler{}

// Mutator changes a field in a CR.
func (h *Handler) Handle(ctx context.Context, req types.Request) types.Response {
	k := &creaturesv1alpha1.Kraken{}

	err := h.Decoder.Decode(req, k)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := k.DeepCopy()

	err = h.mutateKrakenFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(k, copy)
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
