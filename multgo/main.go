package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/satori/go.uuid"
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

func criaMatriz(ordem int) []float64 {

	var matriz []float64
	for i := 0; i < ordem; i++ {
		matriz = append(matriz, 1)
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

	dados.Horario = start.Format(("02/01/2006 15:04:05"))
	var matriz01, matriz02 []float64

	matriz01 = criaMatriz((dados.Ordem) * (dados.Ordem))
	matriz02 = criaMatriz((dados.Ordem) * (dados.Ordem))

	matriz03 := multiplicacaoMatriz(matriz01, matriz02, dados.Ordem)
	fmt.Println(matriz03)

	elapsed := time.Since(start).Milliseconds()
	dados.TempoExecucao = strconv.FormatInt(elapsed, 10)

	log.Printf("Tempo da função GO %s", elapsed)

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

	if req.Path == "/multgo" {
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

func multiplicacaoMatriz(matriz1 []float64, matriz2 []float64, ordemMatriz int) []float64 {
	var matriz3 []float64
	var val float64
	for i := 0; i < ordemMatriz; i++ {
		for j := 0; j < ordemMatriz; j++ {
			val = 0
			for k := 0; k < ordemMatriz; k++ {
				val += matriz1[i*ordemMatriz+k] * matriz2[k*ordemMatriz+j]
			}
			matriz3 = append(matriz3, val)
		}
	}
	return matriz3
}

func main() {
	tableArn := os.Getenv("TABLEARN")
	s := strings.Split(tableArn, "/")
	TABLE_NAME = s[1]

	lambda.Start(route)
}
