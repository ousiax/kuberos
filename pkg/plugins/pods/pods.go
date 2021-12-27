package pods

import (
	"fmt"
	"net/http"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

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

		images := make(map[string]bool)
		for _, c := range pod.Spec.InitContainers {
			images[c.Image] = true
		}
		for _, c := range pod.Spec.Containers {
			images[c.Image] = true
		}

		klog.V(2).Infof("Pod: %s, Images:\n", pod.Name)
		for k := range images {
			klog.V(2).Infof("\t%s", k)
			parts := strings.Split(k, "/")
			for _, p := range parts {
				klog.V(2).Infof("\t\t: %s\n", p)
			}
		}

		klog.V(2).Infof("Pod: %s, containers: %s", &pod.Name, &pod.Spec.Containers)

		reviewResponse := admissionv1.AdmissionResponse{}
		reviewResponse.Allowed = true
		reviewResponse.Result = &metav1.Status{Message: "this webhook allows all requests"}
		return &reviewResponse
	})
}
