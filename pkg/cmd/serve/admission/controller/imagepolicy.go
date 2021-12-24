package controller

import (
	"net/http"
	"time"

	"github.com/qqbuby/kuberos/pkg/cmd/serve/admission"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func imagePolicy(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("always-allow-with-delay sleeping for 5 seconds")
	time.Sleep(5 * time.Second)
	klog.V(2).Info("calling always-allow")
	reviewResponse := admissionv1.AdmissionResponse{}
	reviewResponse.Allowed = true
	reviewResponse.Result = &metav1.Status{Message: "this webhook allows all requests"}
	return &reviewResponse
}

func ServeImagePolicy(w http.ResponseWriter, r *http.Request) {
	admission.Serve(w, r, imagePolicy)
}
