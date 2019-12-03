package testwebhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s-webhook-test/pkg/api/model"
	v1 "k8s.io/api/admission/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
)

type MutatingWebhookHandler struct {
}

var PatchTypeJSONPatch = v1.PatchTypeJSONPatch

func (s MutatingWebhookHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	resp := v1.AdmissionReview{
		Response: &v1.AdmissionResponse{},
	}
	reqReview := v1.AdmissionReview{
		Request: &v1.AdmissionRequest{},
	}

	defer func() {
		resp.Response.UID = reqReview.Request.UID
		resp.Response.PatchType = &PatchTypeJSONPatch
		resptmp := reqReview
		resptmp.Request = nil
		resptmp.Response = resp.Response

		resBody, _ := json.Marshal(resptmp)
		writer.Write(resBody)
		writer.Header().Set("content-type", "application/json")
		return
	}()
	if req.Header.Get("Content-Type") != runtime.ContentTypeJSON || req.Method != "POST" {
		return
	}

	// 读取Body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("读取请求体error : ", err)
		return
	}
	fmt.Println("[DEBUG] : ", string(reqBody))

	err = json.Unmarshal(reqBody, &reqReview)
	if err != nil {
		fmt.Println("请求数据解析error：", err)
		return
	}

	pod := v12.Pod{}
	err = json.Unmarshal(reqReview.Request.Object.Raw, &pod)
	if err != nil {
		fmt.Println("资源对象不为Pod：", err)
		return
	}

	var patches []model.PatchOperation

	// 添加 annotation 注入patch
	patches = append(patches, injectAnnotations(&pod))

	resp.Response.Allowed = true
	resp.Response.Patch, _ = json.Marshal(patches)
}
