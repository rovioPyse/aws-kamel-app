package main

import (
	"context"

	"alarm-processor/internal/processor"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Status string `json:"status"`
}

func handle(ctx context.Context, event events.CloudWatchEvent) (response, error) {
	payload := map[string]any{
		"id":     event.ID,
		"source": event.Source,
		"detail": event.Detail,
	}

	if err := processor.Process(ctx, payload); err != nil {
		return response{}, err
	}

	return response{Status: "ok"}, nil
}

func main() {
	lambda.Start(handle)
}
