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
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

var webhookBuilders []*builder.WebhookBuilder

func Add(mgr manager.Manager) error {
	return add(mgr, webhookBuilders)
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, builders []*builder.WebhookBuilder) error {
	svr, err := webhook.NewServer("foo-admission-server", mgr, webhook.ServerOptions{
		Port:    9876,
		CertDir: "/tmp/cert",
		KVMap:   map[string]interface{}{"foo": "bar"},
		BootstrapOptions: &webhook.BootstrapOptions{
			Secret: &types.NamespacedName{
				Namespace: "default",
				Name:      "foo-admission-server-secret",
			},

			Service: &webhook.Service{
				Namespace: "default",
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

	webhooks := make([]webhook.Webhook, len(builders))
	for i, builder := range builders {
		wh, err := builder.WithManager(mgr).Build()
		if err != nil {
			return err
		}
		webhooks[i] = wh
	}

	return svr.Register(webhooks...)
}
