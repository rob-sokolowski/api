package lib

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"time"
)

type SecretsWrapper struct {
	client      *secretmanager.Client
	vellumAiKey string
}

func InitSecretsWrapper(ctx context.Context) (*SecretsWrapper, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println("failed to setup fsClient:", err)
		return nil, err
	}

	return &SecretsWrapper{client, ""}, nil
}

func (sw *SecretsWrapper) VellumApiKey() (string, error) {
	if sw.vellumAiKey == "" {
		val, err := sw.fetchSecret(context.Background(), "projects/740743074280/secrets/vellum-ai-key/versions/1")
		if err != nil {
			fmt.Println("failed to fetch vellum-ai-key:", err)
			return "", err
		}
		sw.vellumAiKey = val
	}

	return sw.vellumAiKey, nil
}

func (sw *SecretsWrapper) fetchSecret(ctx context.Context, secretKey string) (string, error) {
	start := time.Now()
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretKey,
	}

	result, err := sw.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", err
	}

	dt := time.Since(start)
	fmt.Println("Fetching secret ", secretKey, " took ", dt)
	secretVal := result.Payload.Data
	return string(secretVal), nil
}
