package kubernetes_test

import (
	"testing"

	"github.com/ardaguclu/skipper/dataclients/kubernetes/kubernetestest"
)

func TestIngressFixtures(t *testing.T) {
	kubernetestest.FixturesToTest(t, "testdata/ingress/named-ports")
}
