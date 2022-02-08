package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/klog/v2"
)

// Serve handles the http portion of a request prior to handing to an admit
// function
func ServeAdmission(w http.ResponseWriter, r *http.Request, admit func(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		message := fmt.Sprintf("contentType=%s, expect application/json", contentType)
		klog.Warning(message)
		http.Error(w, message, http.StatusUnsupportedMediaType)
		return
	}

	klog.V(4).Infof("handling request: %s", body)

	deserializer := Codecs.UniversalDeserializer()
	obj, gvk, err := deserializer.Decode(body, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Request could not be decoded: %v", err)
		klog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if *gvk != admissionv1.SchemeGroupVersion.WithKind("AdmissionReview") {
		msg := fmt.Sprintf("Unsupported group version kind: %v", gvk)
		klog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	requestedAdmissionReview, ok := obj.(*admissionv1.AdmissionReview)
	if !ok {
		klog.Errorf("Expected v1.AdmissionReview but got: %T", obj)
		return
	}

	responseAdmissionReview := &admissionv1.AdmissionReview{}
	responseAdmissionReview.SetGroupVersionKind(*gvk)
	responseAdmissionReview.Response = admit(requestedAdmissionReview)
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
	responseObj := responseAdmissionReview

	klog.V(2).Infof("%s: %s", responseObj.Response.Result.Message, responseObj.Response.UID)
	klog.V(4).Infof("sending response: %v", responseObj)
	respBytes, err := json.Marshal(responseObj)
	if err != nil {
		klog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}
