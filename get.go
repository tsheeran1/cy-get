package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Event struct {
	Gtype  string `json:"type"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
	Income int    `json:"income"`
}

type Record struct {
	Userid string `json:"Userid"`
	Age    string `json:"age"`
	Height string `json:"height"`
	Income string `json:"income"`
}

func handler(ctx context.Context, e Event) (string, error) {

	fmt.Println("Event:", e)
	fmt.Println("Type: ", e.Gtype)
	// Create the Dynamodb client -- common for all type options
	config := &aws.Config{
		Region: aws.String("us-east-2"),
	}
	sess := session.Must(session.NewSession(config))
	dbc := dynamodb.New(sess)

	if e.Gtype == "all" {
		// "all" will use Scan method.
		// will require a ScanInput structure called si here.
		// will return a ScanOutput structure called so here
		si := &dynamodb.ScanInput{
			TableName: aws.String("compare-yourself"),
		}
		so, err := dbc.Scan(si) // Read the table
		if err != nil {
			fmt.Println("Scan Error", err)
			return "", err
		}
		fmt.Println("Scan Output", so)

		// so object contains a property Items which is a []map[string]*AttributeValue `type:"list"` each slice contains the dynamodb "attribute value" syntax
		// we need to unmarshall this into our object  we will place the result in a []Record

		data := []Record{}
		err = dynamodbattribute.UnmarshalListOfMaps(so.Items, &data)
		if err != nil {
			return "Error unmarshalling data", err
		}

		// Now if we want to return a string containing a JSON object we need to Marshal the data structure in to a JSON string.

		out, err := json.Marshal(data)
		if err != nil {
			return "Error Marshalling output", err
		}
		
		return fmt.Sprintf("%s", string(out)), nil

	} else if e.Gtype == "single" {

		return "JUST MY DATA", nil
	}

	return fmt.Sprintf("wrong type : %v", e.Gtype), nil
}

func main() {
	lambda.Start(handler)
}
