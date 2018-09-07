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
	crewv1alpha1 "github.com/mengqiy/webhook-scaffolding/pkg/apis/crew/v1alpha1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
	"sigs.k8s.io/controller-runtime/pkg/webhook/types"
)

var Builder *builder.WebhookBuilder

func init() {
	Builder = builder.NewWebhookBuilder().
		Name("validatingfirstmates.k8s.io").
		Type(types.WebhookTypeValidating).
		Path("/validating-firstmates").
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		ForType(&crewv1alpha1.Firstmate{}).
		Handlers(&Handler{})
}
