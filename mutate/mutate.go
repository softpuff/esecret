package mutate

import (
	"encoding/json"
	"fmt"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Mutate(body []byte) ([]byte, error) {
	admReview := v1beta1.AdmissionReview{}

	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("Unmarshalling request failed: %v", err)
	}

	var pod *corev1.Pod

	ar := admReview.Request

	resp := v1beta1.AdmissionResponse{}

	if ar != nil {
		if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("Unable to unmarshal pod json: %v", err)
		}
	}

	resp.Allowed = true
	resp.UID = ar.UID
	_, ok := pod.Labels["entelo-secret"]
	if ok {
		secretName := fmt.Sprintf("%s-secret", pod.Name)

		p := []map[string]string{}

		// scrEnvSource := &corev1.SecretEnvSource{
		// 	LocalObjectReference: corev1.LocalObjectReference{
		// 		Name: secretName,
		// 	},
		// }
		for _, c := range pod.Spec.Containers {
			patch := map[string]string{
				"op":    "replace",
				"path":  fmt.Sprintf("/spec/containers/%s/envFrom/secretRef", c.Name),
				"value": secretName,
			}
			// c.EnvFrom = corev1.EnvFromSource{
			// 	SecretRef: scrEnvSource,
			// }
			p = append(p, patch)

		}

		resp.Patch = json.Marshal(p)

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp

		responseBody, err := json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
		return responseBody, nil
	}
	return nil, nil

}
