import json
import time
import boto3
import os
import numpy as np
from datetime import datetime


client = boto3.client('dynamodb')
tableName = os.getenv('TABLEARN').split('/')[1]

def criaMatriz(ordem):
    matriz = np.ones((ordem,ordem), np.float64)	
    return matriz

def multiplicacaoMatriz(matriz1, matriz2, ordem):
    print("Calculando o produto de Matrizes...")
    matriz3 = np.matmul(matriz1, matriz2)
    return matriz3

def handlerMultiplicaMatriz(body):
    data_e_hora_atuais = datetime.now()
    inicio = int(time.time() * 1000)
    ID = body["id"]
    Ordem = body["ordem"]
    Linguagem = body["linguagem"]

    if not ID or not Ordem:
        return json.dumps({'error': 'Please provide Id and ordem'}), 400

    matriz01 = criaMatriz(Ordem)
    matriz02 = criaMatriz(Ordem)
    matriz03 = multiplicacaoMatriz(matriz01, matriz02, Ordem)
	
    fim = int(time.time() * 1000)
    TempoExecucao = (fim-inicio)
    Horario = data_e_hora_atuais.strftime("%d/%m/%Y %H:%M:%S") 	

    print(matriz03)
    print(tableName)
    resp = client.put_item(
    	TableName=tableName,
    	Item={
    		'id': {'S': ID },
    		'ordem': {'S': str(Ordem) },
    		'linguagem': {'S': Linguagem },
    		'horario': {'S': str(Horario) },
    		'tempo': {'S': str(TempoExecucao) }
    	}
    )

    return json.dumps({
		'id': ID,
		'ordem': Ordem,
		'linguagem': Linguagem,
		'horario': Horario,
		'tempo': TempoExecucao
	})

def lambda_handler(event, context): 
    body = json.loads(event["body"])
    result = handlerMultiplicaMatriz(body)
    return {
        'statusCode': 200,
        'body': json.dumps(result)
    }

'''
JSON PARA TESTE
{
    "id" : "1",
    "ordem" : 10,
    "linguagem" : "python",
    "horario" : "",
    "tempo" : ""
}

'''
