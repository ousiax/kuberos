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
		podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
		if ar.Request.Resource != podResource {
			err := fmt.Errorf("uid: %s, expect resource to be %s", ar.Request.UID, podResource)
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
		klog.V(2).Infof("admitting pod: %s/%s", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)

		registry, ok := os.LookupEnv("REQUIRED_IMAGE_REGISTRY")
		if !ok {
			klog.Warning("allows the request without admit: could not retrieve the value of the environment variable named %s", registry)
			reviewResponse := admissionv1.AdmissionResponse{}
			reviewResponse.Allowed = true
			reviewResponse.Result = &metav1.Status{Message: "Kuberos: allows the request without admit"}
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
		reviewResponse := admissionv1.AdmissionResponse{}
		for img := range images {
			ref, err := containers.ParseImageRef(img)
			if err != nil {
				reviewResponse.Allowed = false
				reviewResponse.Result = &metav1.Status{
					Message: err.Error(),
				}
				return &reviewResponse
			}
			if ref.Registry != registry {
				reviewResponse.Allowed = false
				reviewResponse.Result = &metav1.Status{
					Message: fmt.Sprintf("Kuberos: pod: %s/%s: %s must be at [%s]", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, img, registry),
				}
				return &reviewResponse
			}
			if org != "" && ref.Org != org {
				reviewResponse.Allowed = false
				reviewResponse.Result = &metav1.Status{
					Message: fmt.Sprintf("Kuberos: pod: %s/%s: %s must be at [%s/%s]", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, img, registry, org),
				}
				return &reviewResponse
			}
		}
		reviewResponse.Allowed = true
		reviewResponse.Result = &metav1.Status{Message: "Kuberos: allows the request"}
		return &reviewResponse
	})
}
