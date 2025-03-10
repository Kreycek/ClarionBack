package movement

import (
	"Clarion/internal/db"
	"Clarion/internal/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllMovements(client *mongo.Client, dbName, collectionName string, page, limit int) ([]any, int, error) {
	collection := db.GetCollection(client, dbName, collectionName)

	// Criar o filtro (por enquanto vazio, pode ser expandido)
	filter := bson.M{}

	// Obter a contagem total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar documentos: %v", err)
	}

	// Definir opções de busca com paginação
	options := options.Find()
	options.SetLimit(int64(limit))
	options.SetSkip(int64((page - 1) * limit))

	// Buscar usuários com paginação
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar movimentos: %v", err)
	}
	defer cursor.Close(context.Background())

	var movements []any
	for cursor.Next(context.Background()) {
		var movement models.Movement
		if err := cursor.Decode(&movement); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar movimento: %v", err)
		}

		// Adiciona os usuários formatados
		movements = append(movements, map[string]any{
			"ID":          movement.ID.Hex(), // Agora o campo ID é uma string
			"CodDaily":    movement.CodDaily,
			"CodDocument": movement.CodDocument,
			"Date":        movement.Date,
			"Active":      movement.Active,
			"Movements":   movement.Movements,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar no cursor: %v", err)
	}

	return movements, int(total), nil
}

func GetMovementByID(client *mongo.Client, dbName, collectionName, movementId string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// Criar filtro para buscar um usuário pelo ID
	fmt.Println("id", movementId)

	objectID, erroId := primitive.ObjectIDFromHex(movementId)
	if erroId != nil {
		log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	}

	filter := bson.M{"_id": objectID}

	// Variável para armazenar o usuário retornado
	var movement models.Movement

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&movement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("plano de contas não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string
	movementId = movement.ID.Hex()

	// Retornar o usuário como um mapa
	movements := map[string]any{
		"ID":          movement.ID.Hex(), // Agora o campo ID é uma string
		"CodDaily":    movement.CodDaily,
		"CodDocument": movement.CodDocument,
		"Movements":   movement.Movements,
		"Date":        movement.Date,
		"Active":      movement.Active,
	}

	fmt.Println("COAData", movements)

	return movements, nil
}

// Função para inserir um usuário na coleção "user"
func InsertMovement(client *mongo.Client, dbName, collectionName string, movement models.Movement) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Criar um contexto para a operação de inserção
	ctx := context.Background()

	// Inserir o documento
	_, err := collection.InsertOne(ctx, movement)
	if err != nil {
		return fmt.Errorf("erro ao inserir plano de contas: %v", err)
	}

	return nil
}

func SearchMovements(client *mongo.Client, dbName, collectionName string, typeSearchDate int, CodDaily, CodDocument, date *string, page, limit int64) ([]any, int64, error) {
	collection := client.Database(dbName).Collection(collectionName)

	// Criando o filtro dinâmico
	filter := bson.M{}
	if CodDaily != nil && *CodDaily != "" {
		filter["codDaily"] = *CodDaily // Apenas atribui diretamente
	}

	if CodDocument != nil && *CodDocument != "" {
		// filter["codDocument"] = bson.M{"$regex": *CodDocument, "$options": "i"}
		filter["codDocument"] = *CodDocument // Apenas atribui diretamente
	}

	var parsedTime time.Time // Declara parsedTime fora do if
	var err error            // Declara err fora do if

	if date != nil && *date != "" {
		// Convertendo string para tempo (assumindo que o formato esperado seja YYYY-MM-DD ou YYYY-MM)
		layoutDay := "2006-01-02"
		layoutMonth := "2006-01"

		if typeSearchDate == 1 {
			// Corrigido para usar a variável declarada fora do if
			parsedTime, err = time.Parse(layoutDay, *date)
			if err != nil {
				return nil, 0, fmt.Errorf("formato de data inválido, use YYYY-MM-DD")
			}

			startOfDay := parsedTime
			endOfDay := startOfDay.Add(24 * time.Hour)

			// Filtro no MongoDB com o tipo time.Time
			filter["date"] = bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			}

		} else if typeSearchDate == 2 {

			trimmedDate := strings.TrimSpace(*date)

			// Verifique se a data tem dia (caso sim, remova)
			if strings.Count(trimmedDate, "-") == 2 { // Caso seja YYYY-MM-DD
				// Apenas mantemos ano e mês (removemos o dia)
				trimmedDate = trimmedDate[:7] // Apenas "YYYY-MM"
			}

			fmt.Println("Data de mês/ano a ser analisada:", trimmedDate)

			// Tenta analisar como YYYY-MM
			parsedTime, err = time.Parse(layoutMonth, trimmedDate)
			if err != nil {
				return nil, 0, fmt.Errorf("formato de data inválido para mês/ano, use YYYY-MM: %v", err)
			}

			firstDay := parsedTime
			lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second) // Último segundo do mês

			filter["date"] = bson.M{
				"$gte": firstDay,
				"$lt":  lastDay,
			}
		} else {
			// Caso nenhum tipo de pesquisa válido
			return nil, 0, fmt.Errorf("tipo de pesquisa de data inválido")
		}
	}

	// Contar total de usuários antes da paginação
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Executa a consulta com paginação
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Processa os resultados
	var movements []any
	for cursor.Next(context.Background()) {
		var movement models.Movement
		if err := cursor.Decode(&movement); err != nil {
			return nil, 0, fmt.Errorf("erro ao decodificar usuário: %v", err)
		}

		movements = append(movements, map[string]any{
			"ID":          movement.ID.Hex(), // Agora o campo ID é uma string
			"CodDaily":    movement.CodDaily,
			"CodDocument": movement.CodDocument,
			"Movements":   movement.Movements,
			"Date":        movement.Date,
			"Active":      movement.Active,
		})
	}

	fmt.Println("Total registros", total)

	// Retorna usuários e total de registros
	return movements, total, nil
}
