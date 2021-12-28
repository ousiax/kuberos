package pods

import (
	"fmt"
	"net/http"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/qqbuby/kuberos/pkg/plugins/pods/containers"
	"github.com/qqbuby/kuberos/pkg/plugins/util"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.ServeAdmission(w, r, func(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
		klog.V(2).Info("admitting pods")
		podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
		if ar.Request.Resource != podResource {
			err := fmt.Errorf("expect resource to be %s", podResource)
			klog.Error(err)
			return util.V1AdmissionResponse(err)
		}

		raw := ar.Request.Object.Raw
		pod := corev1.Pod{}
		deserializer := util.Codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
			klog.Error(err)
			return util.V1AdmissionResponse(err)
		}

		registry, ok := os.LookupEnv("REQUIRED_IMAGE_REGISTRY")
		if !ok {
			klog.Info("could not retrieve the value of the environment variable named %s", registry)
			reviewResponse := admissionv1.AdmissionResponse{}
			reviewResponse.Allowed = true
			reviewResponse.Result = &metav1.Status{Message: "this webhook allows all requests"}
			return &reviewResponse
		}

		images := make(map[string]bool)
		for _, c := range pod.Spec.InitContainers {
			images[c.Image] = true
		}
		for _, c := range pod.Spec.Containers {
			images[c.Image] = true
		}

		org := os.Getenv("REQUIRED_IMAGE_ORG")
		for img := range images {
			ref, err := containers.ParseImageRef(img)
			if err != nil {
				return util.V1AdmissionResponse(err)
			}
			if ref.Registry != registry {
				return util.V1AdmissionResponse(fmt.Errorf("%s must be at [%s]", img, registry))
			}
			if org != "" && ref.Org != org {
				return util.V1AdmissionResponse(fmt.Errorf("%s must be at [%s/%s]", img, registry, org))
			}
		}

		return nil
	})
}
