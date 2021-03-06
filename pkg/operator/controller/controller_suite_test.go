/*
Copyright 2018 The CDI Authors.

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

package controller

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"kubevirt.io/containerized-data-importer/pkg/operator/resources/operator"
	"kubevirt.io/containerized-data-importer/tests/reporters"
)

func TestOperator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Controller Suite", reporters.NewReporters())
}

var testenv *envtest.Environment
var cfg *rest.Config
var clientset *kubernetes.Clientset

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	env := &envtest.Environment{}

	var err error
	cfg, err = env.Start()
	Expect(err).NotTo(HaveOccurred())

	clientset, err = kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	opts := envtest.CRDInstallOptions{
		CRDs: []client.Object{operator.NewCdiCrd()},
	}

	crds, err := envtest.InstallCRDs(cfg, opts)
	Expect(err).NotTo(HaveOccurred())
	err = envtest.WaitForCRDs(cfg, crds, envtest.CRDInstallOptions{})
	Expect(err).NotTo(HaveOccurred())

	// Prevent the metrics listener being created
	metrics.DefaultBindAddress = "0"

	testenv = env

	close(done)
}, 60)

var _ = AfterSuite(func() {
	if testenv == nil {
		return
	}

	testenv.Stop()

	// Put the DefaultBindAddress back
	metrics.DefaultBindAddress = ":8080"
})
