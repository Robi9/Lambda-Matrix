package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gonum.org/v1/gonum/mat"
    "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const AWS_REGION = "us-east-1"

var TABLE_NAME string

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(AWS_REGION))

type MyEvent struct {
	ID            string `json:"id"`
	Ordem         int    `json:"ordem"` //ordem
	Linguagem     string `json:"linguagem"`
	Horario       string `json:"horario"`
	TempoExecucao string `json:"tempo"`
}

func criaMatriz(ordem int) *mat.Dense {
   matriz := mat.NewDense(ordem, ordem, nil)

   for i := 0; i < ordem; i++ {
      for j := 0; j < ordem; j++ {
          matriz.Set(i, j, 1.0)
      }
   }

   return matriz
}

// Grava os dados da multiplicação no dynamo
func GravaDadosMult(dados MyEvent) error { 

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(dados.ID),
			},
			"ordem": {
				N: aws.String(strconv.Itoa(dados.Ordem)),
			},
			"linguagem": {
				S: aws.String(dados.Linguagem),
			},
			"horario": {
				S: aws.String(dados.Horario),
			},
			"tempo": {
				S: aws.String(dados.TempoExecucao),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

func handleMultiplicaMatriz(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()

	var dados MyEvent
	err := json.Unmarshal([]byte(req.Body), &dados)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
			Body:       err.Error(),
		}, nil
	}
	fmt.Println(dados.Ordem)
	dados.Horario = start.Format(("02/01/2006 15:04:05"))

    matriz01 := criaMatriz(dados.Ordem)
    matriz02 := criaMatriz(dados.Ordem)

	matriz03 := multiplicacaoMatriz(matriz01, matriz02, dados.Ordem)
	//fmt.Println(matriz03)

	elapsed := time.Since(start).Milliseconds()
	dados.TempoExecucao = strconv.FormatInt(elapsed, 10)

	log.Printf("Tempo da função GO %s", elapsed)
	fmt.Println(matriz03)
	err = GravaDadosMult(dados)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       "Calculado",
	}, nil
}

func route(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if req.Path == "/gonum" {
		if req.HTTPMethod == "POST" {
			return handleMultiplicaMatriz(req)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       http.StatusText(http.StatusMethodNotAllowed),
	}, nil

}

func multiplicacaoMatriz(matriz1 *mat.Dense, matriz2 *mat.Dense, ordemMatriz int) mat.Dense {
   var matriz3 mat.Dense
   matriz3.Mul(matriz1, matriz2)
   return matriz3
}

func main() {
	tableArn := os.Getenv("TABLEARN")
	s := strings.Split(tableArn, "/")
	TABLE_NAME = s[1]

	lambda.Start(route)
}
