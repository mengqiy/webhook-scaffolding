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

package server

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

var (
	log             = logf.Log.WithName("server")
	webhookBuilders []*builder.WebhookBuilder
	ListOfHandlers  [][]admission.Handler
)

func Add(mgr manager.Manager) error {
	svr, err := webhook.NewServer("foo-admission-server", mgr, webhook.ServerOptions{
		// TODO(user): change the configuration of ServerOptions based on your need.
		Port:    9876,
		CertDir: "/tmp/cert",
		BootstrapOptions: &webhook.BootstrapOptions{
			Secret: &types.NamespacedName{
				Namespace: "system",
				Name:      "foo-admission-server-secret",
			},

			Service: &webhook.Service{
				Namespace: "system",
				Name:      "foo-admission-server-service",
				// Selectors should select the pods that runs this webhook server.
				Selectors: map[string]string{
					"app": "foo-admission-server",
				},
			},
		},
	})
	if err != nil {
		return err
	}

	webhooks := make([]webhook.Webhook, len(webhookBuilders))
	for i, builder := range webhookBuilders {
		wh, err := builder.
			Handlers(ListOfHandlers[i]...).
			WithManager(mgr).
			Build()
		if err != nil {
			return err
		}
		webhooks[i] = wh
	}

	return svr.Register(webhooks...)
}
