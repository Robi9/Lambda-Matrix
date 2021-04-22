import json
import time
import boto3
import os
from datetime import datetime


client = boto3.client('dynamodb')
tableName = os.getenv('TABLEARN').split('/')[1]

def criaMatriz(ordem):
	matriz = []
	for i in range(ordem): 
		matriz.append(1)
	return matriz

def multiplicacaoMatriz(matriz1, matriz2, Ordem):
    matriz3 = []
    print("Calculando o produto de Matrizes...")
    
    for i in range(Ordem):
        for j in range(Ordem):
            val = 0			
            for k in range(Ordem):
                val += matriz1[i*Ordem+k] * matriz2[k*Ordem+j]
            matriz3.append(val)		
	
    return matriz3

def handlerMultiplicaMatriz(body):
    data_e_hora_atuais = datetime.now()
    inicio = int(time.time() * 1000)
    ID = body["id"]
    Ordem = body["ordem"]
    Linguagem = body["linguagem"]

    if not ID or not Ordem:
        return json.dumps({'error': 'Please provide Id and ordem'}), 400

    matriz01 = criaMatriz(Ordem*Ordem)
    matriz02 = criaMatriz(Ordem*Ordem)
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
