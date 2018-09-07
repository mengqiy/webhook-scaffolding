package validating

import (
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

var (
	Builders   = map[string]*builder.WebhookBuilder{}
	HandlerMap = map[string][]admission.Handler{}
)
