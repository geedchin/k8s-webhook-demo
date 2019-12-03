package main

import (
	"k8s-webhook-test/pkg/api/testwebhook"
	"net/http"
)

func main() {

	certPath := "file/server.crt"
	keyPath := "file/server-key.pem"

	server := &http.ServeMux{}

	server.Handle("/mutatingwebhook", testwebhook.MutatingWebhookHandler{})

	http.ListenAndServeTLS(":1443", certPath, keyPath, server)
}
