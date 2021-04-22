package main

import (
    "net/url"
    "os"
    "strings"
    //"log"
    "encoding/json"
    "net/http"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const AWS_REGION = "us-east-1"
var TABLE_NAME string

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(AWS_REGION))

type MyEvent struct {
        ID string `json:"id"`
        Ordem int `json:"ordem"` 
        Linguagem string `json:"linguagem"`
        Horario string `json:"horario"`
        TempoExecucao string `json:"tempo"`
}

// Retorna os dados da multiplicação do ID passado por string de consulta
func RetornaDadosMult(id string) (MyEvent, error) {
input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        "id": {
            S: aws.String(id),
        },
    },
    TableName: aws.String(TABLE_NAME),
  }

  result, err := db.GetItem(input) 
  if err != nil {
    return MyEvent{}, err

  }

  if len(result.Item) == 0 {
    return MyEvent{}, nil
  }


  var dados MyEvent
  err = dynamodbattribute.UnmarshalMap(result.Item, &dados)
  if err != nil {
    return MyEvent{}, err
  }
  return dados, nil
}

func handleRetornaResultados(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  rawId,_ := req.PathParameters["id"]

  // path parameters are typically URL encoded so to get the value
  value, err := url.QueryUnescape(rawId)
  if nil != err {
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusInternalServerError,
      Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
      Body: http.StatusText(http.StatusInternalServerError),
    }, err
  }

  dados, err := RetornaDadosMult(string(value))
  if err != nil {
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusInternalServerError,
      Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
      Body: http.StatusText(http.StatusInternalServerError),
    }, nil
  }

  js, err := json.Marshal(dados)
  if err != nil {
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusInternalServerError,
      Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
      Body: http.StatusText(http.StatusInternalServerError),
    }, nil
  }

  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
    Body: string(js),
  }, nil

}

func route (req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){ 
        if req.HTTPMethod == "GET" {
            return handleRetornaResultados(req)
        }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusMethodNotAllowed,
        Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
        Body:       http.StatusText(http.StatusMethodNotAllowed),
    }, nil

}

func main() {
	tableArn := os.Getenv("TABLEARN")
	s := strings.Split(tableArn, "/")
	TABLE_NAME = s[1]

   lambda.Start(route)
}

/*
input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        "id": {
            S: aws.String("35300b11-246f-494d-8c8f-414a48797786"),
        },
    },
    TableName: aws.String(TABLE_NAME),
  }

  result, err := db.GetItem(input)
  if err != nil {
    return []MyEvent{}, err
  }
  if len(result.Item) == 0 {
    return []MyEvent{}, nil
  }

  var dados[]MyEvent
  err = dynamodbattribute.UnmarshalMap(result.Item, &dados)
  if err != nil {
    return []MyEvent{}, err
  }

  return dados, nil
*/
