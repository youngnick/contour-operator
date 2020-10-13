package contour

import (
	"testing"

	operatorv1alpha1 "github.com/projectcontour/contour-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testJobImageName = "test-job:main"

func checkJobHasEnvVar(t *testing.T, job *batchv1.Job, name string) {
	t.Helper()

	for _, envVar := range job.Spec.Template.Spec.Containers[0].Env {
		if envVar.Name == name {
			return
		}
	}
	t.Errorf("job is missing environment variable %q", name)
}

func checkJobHasContainer(t *testing.T, job *batchv1.Job, name string) *corev1.Container {
	t.Helper()

	for _, container := range job.Spec.Template.Spec.Containers {
		if container.Name == name {
			return &container
		}
	}
	t.Errorf("job is missing container %q", name)
	return nil
}

func checkContainerHasImage(t *testing.T, container *corev1.Container, image string) {
	t.Helper()

	if container.Image == image {
		return
	}
	t.Errorf("job is missing image %q", image)
}

func TestDesiredJob(t *testing.T) {
	ctr := &operatorv1alpha1.Contour{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-job",
			Namespace: "test-job-ns",
		},
	}

	job, err := DesiredJob(ctr, testJobImageName)
	if err != nil {
		t.Errorf("invalid job: %w", err)
	}

	container := checkJobHasContainer(t, job, jobContainerName)
	checkContainerHasImage(t, container, testJobImageName)
	checkJobHasEnvVar(t, job, jobNsEnvVar)
}