package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Event struct {
	AccessToken string `json:"accessToken"`
	Gtype       string `json:"type"`
}

type Record struct {
	Age    int `json:"age"`
	Height int `json:"height"`
	Income int `json:"income"`
}

type Keystruct struct {
	Userid string
}

func handler(ctx context.Context, e Event) ([]Record, error) {

	fmt.Println("Event:", e)
	fmt.Println("Type: ", e.Gtype)
	// Create the config, session and Dynamodb client -- common for all type options
	config := &aws.Config{
		Region: aws.String("us-east-2"),
	}
	sess := session.Must(session.NewSession(config))
	dbc := dynamodb.New(sess)
	//create cognito service provider for this session
	cisp := cognitoidentityprovider.New(sess)

	accessToken := e.AccessToken
	fmt.Printf("TOKEN: %v", accessToken)

	if e.Gtype == "all" { //  GET ALL RECORDS
		// "all" will use Scan method.
		// will require a ScanInput structure called si here, containing the table name we want
		// will return a ScanOutput structure called so here
		si := &dynamodb.ScanInput{
			TableName: aws.String("compare-yourself"),
		}
		so, err := dbc.Scan(si) // Read the table
		if err != nil {
			fmt.Println("Scan Error", err)
			return []Record{}, err
		}

		// so object contains a property Items which is a []map[string]*AttributeValue `type:"list"` each slice contains the dynamodb "attribute value" syntax
		// we need to unmarshall this into our object  we will place the result in a []Record

		var data []Record
		err = dynamodbattribute.UnmarshalListOfMaps(so.Items, &data)
		if err != nil {
			fmt.Println("Error unmarshalling data")
			return []Record{}, err
		}

		// Now if we want to return a string containing a JSON object we need to Marshal the data structure in to a JSON string.

		// DONT MARSHAL  out, err := json.Marshal(data)
		// if err != nil {
		// 	fmt.Println("Error Marshalling output")
		// 	return []Record{}, err
		// }
		//DEBUG	fmt.Printf("%s", out)

		return data, nil

	} else if e.Gtype == "single" {
		// First need to create GetUserInput struct and then call cisp.GetUser to get the current user

		getui := &cognitoidentityprovider.GetUserInput{
			AccessToken: aws.String(accessToken),
		}
		getuo, err := cisp.GetUser(getui)
		if err != nil {
			fmt.Println(err)
			return []Record{}, err
		}
		fmt.Println(getuo)
		userID := getuo.UserAttributes[0].Value
		// We have current user in userID.  Now need to set up keyvalyre and dynamodb structs
		// "single will use GetItem method.
		//  We will need a GetItemInput structure containing the filename and key."
		// key will need to be marshalled into dynamodb attribute value
		//Create keyvalue and marshal into dynamodb attribute value form
		//userID := "28e09c22-f67a-4e34-bb22-51f3ddcf2da1"
		keyval := Keystruct{Userid: *userID}

		av, err := dynamodbattribute.MarshalMap(keyval)
		if err != nil {
			fmt.Println("Unable to marshal key structure")
			return []Record{}, err
		}
		// Create GetItemInput structure
		gi := &dynamodb.GetItemInput{
			TableName: aws.String("compare-yourself"),
			Key:       av,
		}
		// Get the item; put result dynamodb structure in gout
		gout, err := dbc.GetItem(gi)
		if err != nil {
			fmt.Println("GetItem failure")
			return []Record{}, err
		}
		// gout.Item contains a map[string]*AttributeValue and needs to be dynamodb-unmarshalled
		if gout.Item == nil { // there was no record found
			err = fmt.Errorf("Record not found for user %s", *userID)
			fmt.Println(err.Error())
			return []Record{}, err
		}
		fmt.Println("GOUT:", gout) // otherwise we got a record
		var r Record               //record to unmarshal into
		err = dynamodbattribute.UnmarshalMap(gout.Item, &r)

		fmt.Println(r)
		rs := []Record{r}
		return rs, nil

	} else if e.Gtype == "test" {

		r := Record{42, 24, 14}
		// DONT MARSHAL out, err := json.Marshal(r)
		// if err != nil {
		// 	return []byte("test error marshalling json"), err
		// }
		// fmt.Println("TestOUT", string(out))
		rs := []Record{r}
		return rs, nil
	}

	return []Record{}, fmt.Errorf("wrong value %s", e.Gtype)
}

func main() {
	lambda.Start(handler)
}
