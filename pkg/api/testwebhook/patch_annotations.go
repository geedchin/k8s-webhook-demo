package testwebhook

import (
	"fmt"
	"k8s-webhook-test/pkg/api/model"
	"k8s-webhook-test/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func injectAnnotations(object v1.Object) model.PatchOperation {

	annotations := object.GetAnnotations()
	myWebhookVal := utils.MapGet(annotations, "my-webhook")

	patchOpr := model.PatchOperation{
		Op:    "add",
		Path:  "/metadata/annotations/my-webhook",
		Value: fmt.Sprintf("%s--%s", "exist", myWebhookVal),
	}
	if myWebhookVal == "" {
		patchOpr.Value = "my-def"
	}
	// 如果 annotations为空处理
	if annotations == nil {
		patchOpr.Path = "/metadata/annotations"
		patchOpr.Value = map[string]string{"my-webhook": "my-def"}
	}
	return patchOpr
}
