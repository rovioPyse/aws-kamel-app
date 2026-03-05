package main

import (
	"context"

	"aws-kamel-app/services/lambdas/latest-state-writer/internal/processor"
	"github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Status string `json:"status"`
}

func handle(ctx context.Context, event map[string]any) (response, error) {
	if err := processor.Process(ctx, event); err != nil {
		return response{}, err
	}

	return response{Status: "ok"}, nil
}

func main() {
	lambda.Start(handle)
}
