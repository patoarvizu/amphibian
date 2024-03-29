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

package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	terraformv1 "github.com/patoarvizu/amphibian/api/v1"
	// +kubebuilder:scaffold:imports
)

func createRemoteStateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-remote",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "remote",
			RemoteConfig: terraformv1.RemoteConfig{
				Hostname:     "app.terraform.io",
				Organization: "patoarvizu",
				Workspaces: terraformv1.WorkspacesConfig{
					Name: "amphibian-test-state",
				},
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-remote",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createS3StateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-s3",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "s3",
			S3Config: terraformv1.S3Config{
				Bucket: "patoarvizu-terraform-states",
				Key:    "patoarvizu-infra/amphibian/s3-state/terraform.tfstate",
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-s3",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createConsulStateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-consul",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "consul",
			ConsulConfig: terraformv1.ConsulConfig{
				Path:    "state",
				Address: "consul-server.consul:8500",
				Scheme:  "http",
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-consul",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createKubernetesStateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-kubernetes",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "kubernetes",
			KubernetesConfig: terraformv1.KubernetesConfig{
				SecretSuffix:    "state",
				InClusterConfig: true,
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-kubernetes",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createGCSStateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-gcs",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "gcs",
			GCSConfig: terraformv1.GCSConfig{
				Bucket: "patoarvizu-terraform-states",
				Prefix: "patoarvizu-infra/amphibian/gcs-state",
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-gcs",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createPostgresStateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pg",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "pg",
			PostgresConfig: terraformv1.PostgresConfig{
				ConnStr: "postgres://postgres:postgres123@postgresql.pg:5432/terraform_backend?sslmode=disable",
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-pg",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createArtifactoryStateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-artifactory",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "artifactory",
			ArtifactoryConfig: terraformv1.ArtifactoryConfig{
				Url:     "http://localhost:8082/artifactory",
				Repo:    "example-repo-local",
				Subpath: "/",
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-artifactory",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createEtcdV3StateConfig(targetType string) (*terraformv1.TerraformState, error) {
	s := &terraformv1.TerraformState{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TerraformState",
			APIVersion: "terraform.patoarvizu.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-etcdv3",
			Namespace: "default",
		},
		Spec: terraformv1.TerraformStateSpec{
			Type: "etcdv3",
			EtcdV3Config: terraformv1.EtcdV3Config{
				Endpoints: []string{"etcdv3.etcdv3:2379"},
				Username:  "root",
				Password:  "root123",
			},
			Target: terraformv1.Target{
				Type: targetType,
				Name: "test-etcdv3",
			},
		},
	}
	err := k8sClient.Create(context.TODO(), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:  []string{filepath.Join("..", "..", "config", "crd", "bases")},
		UseExistingCluster: newTrue(),
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = terraformv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = Describe("With the controller running", func() {
	var (
		state *terraformv1.TerraformState
		err   error
	)
	When("Deploying a TerraformState object with 'remote' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createRemoteStateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'remote' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createRemoteStateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 's3' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createS3StateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 's3' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createS3StateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'consul' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createConsulStateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'consul' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createConsulStateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'kubernetes' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createKubernetesStateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'kubernetes' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createKubernetesStateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'gcs' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createGCSStateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'gcs' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createGCSStateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'pg' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createPostgresStateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'pg' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createPostgresStateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'etcdv3' config and target type 'configmap'", func() {
		It("Should create the target ConfigMap", func() {
			state, err = createEtcdV3StateConfig("configmap")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetConfigMap(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	When("Deploying a TerraformState object with 'etcdv3' config and target type 'secret'", func() {
		It("Should create the target Secret", func() {
			state, err = createEtcdV3StateConfig("secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(state).ToNot(BeNil())
			err = validateStateTargetSecret(state)
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Delete(context.TODO(), state)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func validateStateTargetConfigMap(s *terraformv1.TerraformState) error {
	configMap := &corev1.ConfigMap{}
	err := wait.Poll(time.Second*2, time.Second*60, func() (done bool, err error) {
		err = k8sClient.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: s.Spec.Target.Name}, configMap)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	if configMap.Data["hello"] != "world" {
		return errors.New("ConfigMap data doesn't match remote state")
	}
	jsonMap := make(map[string]string)
	err = json.Unmarshal([]byte(configMap.Data["map"]), &jsonMap)
	if err != nil {
		return err
	}
	if jsonMap["a"] != "b" || jsonMap["x"] != "y" {
		return errors.New("ConfigMap data doesn't match remote state")
	}
	jsonList := []string{}
	err = json.Unmarshal([]byte(configMap.Data["list"]), &jsonList)
	if err != nil {
		return err
	}
	if jsonList[0] != "a" || jsonList[1] != "b" || jsonList[2] != "c" {
		return errors.New("ConfigMap data doesn't match remote state")
	}
	return nil
}

func validateStateTargetSecret(s *terraformv1.TerraformState) error {
	secret := &corev1.Secret{}
	err := wait.Poll(time.Second*2, time.Second*60, func() (done bool, err error) {
		err = k8sClient.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: s.Spec.Target.Name}, secret)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	if string(secret.Data["hello"]) != "world" {
		return errors.New("Secret data doesn't match remote state")
	}
	jsonMap := make(map[string]string)
	err = json.Unmarshal([]byte(secret.Data["map"]), &jsonMap)
	if err != nil {
		return err
	}
	if jsonMap["a"] != "b" || jsonMap["x"] != "y" {
		return errors.New("Secret data doesn't match remote state")
	}
	jsonList := []string{}
	err = json.Unmarshal([]byte(secret.Data["list"]), &jsonList)
	if err != nil {
		return err
	}
	if jsonList[0] != "a" || jsonList[1] != "b" || jsonList[2] != "c" {
		return errors.New("Secret data doesn't match remote state")
	}
	return nil
}

func newTrue() *bool {
	b := true
	return &b
}
