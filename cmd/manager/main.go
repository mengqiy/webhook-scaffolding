/*

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

package main

import (
	"flag"
	"os"

	externalapis "github.com/mengqiy/example-crd-apis/pkg/apis"
	"github.com/mengqiy/webhook-scaffolding/pkg/apis"
	"github.com/mengqiy/webhook-scaffolding/pkg/controller"
	"github.com/mengqiy/webhook-scaffolding/pkg/webhook"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var log = logf.Log.WithName("example-controller")

func main() {
	flag.Parse()
	logf.SetLogger(logf.ZapLogger(false))
	entryLog := log.WithName("entrypoint")

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		entryLog.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		entryLog.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}
	if err := externalapis.AddToScheme(mgr.GetScheme()); err != nil {
		entryLog.Error(err, "unable add external APIs to scheme")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		entryLog.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	if err := webhook.AddToManager(mgr); err != nil {
		entryLog.Error(err, "unable to register webhooks to the manager")
		os.Exit(1)
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	log.Error(mgr.Start(signals.SetupSignalHandler()), "unable to run manager")
}
