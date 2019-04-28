package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Gtype  string `json:"type"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
	Income int    `json:"income"`
}

func handler(ctx context.Context, e Event) (string, error) {

	fmt.Println("Event:", e)
	fmt.Println("Type: ", e.Gtype)

	if e.Gtype == "all" {
		return "ALL THE DATA", nil
	} else if e.Gtype == "single" {
		return "JUST MY DATA", nil
	}

	return fmt.Sprintf("wrong type : %v", e.Gtype), nil
}

func main() {
	lambda.Start(handler)
}
