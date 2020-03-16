/*
Copyright 2018 The Kubernetes Authors.

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

package main

import (
    "fmt"
    "time"
	"encoding/json"

	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

    examplev1 "website-controller/pkg/apis/smartx.com/v1"
)

const (
	websitePatch1 string = `[
         { "op": "replace", "path": "/spec/deploymentName", "value": "%s" }
     ]`
)

func mutateWebsite(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.Info("mutating custom resource")
    cr := new(examplev1.Website)

	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

    deploymentName := fmt.Sprintf("%s-%d", cr.Spec.DeploymentName, time.Now().Unix())
    klog.Info(deploymentName)
    reviewResponse.Patch = []byte(fmt.Sprintf(websitePatch1,deploymentName))

	pt := v1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	return &reviewResponse
}

func admitWebsite(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.Info("admitting custom resource")

    cr := new(examplev1.Website)

	var raw []byte
	if ar.Request.Operation == v1.Delete {
		raw = ar.Request.OldObject.Raw
	} else {
		raw = ar.Request.Object.Raw
	}
	err := json.Unmarshal(raw, &cr)
	if err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
    //klog.Info(cr)
    //klog.Info(cr.Spec.DeploymentName)
    //klog.Info(*cr.Spec.Replicas)
    deploymentName := cr.Spec.DeploymentName
    replicasPointer := cr.Spec.Replicas
    if deploymentName == ""{
        reviewResponse.Allowed = false
        reviewResponse.Result = &metav1.Status{
            Reason : "Request Website.spec.deploymentName",
        }
        return &reviewResponse
    }else if replicasPointer == nil{
        reviewResponse.Allowed = false
        reviewResponse.Result = &metav1.Status{
            Reason : "Request Website.spec.replicas",
        }
        return &reviewResponse
    }else if *replicasPointer <= 0 || *replicasPointer > 5{
        reviewResponse.Allowed = false
        reviewResponse.Result = &metav1.Status{
            Reason : "Request Website.spec.replicas number n: 0 < n <= 5",
        }
        return &reviewResponse
    }
    return &reviewResponse
}
