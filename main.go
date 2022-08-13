package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func main() {
	endpoint := os.Getenv("AZURE_COSMOS_ENDPOINT")
	if endpoint == "" {
		log.Fatal("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key := os.Getenv("AZURE_COSMOS_KEY")
	if key == "" {
		log.Fatal("AZURE_COSMOS_KEY could not be found")
	}

	var databaseName = "VisitorCounter"
	var containerName = "Counter"
	var partitionKey = "/countId"

	item := struct {
		ID           string `json:"id"`
		CountId      int
		CreationDate string
	}{
		ID:           "1",
		CountId:      0,
		CreationDate: "",
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		log.Fatal("Failed to create a credential: ", err)
	}

	// Create a CosmosDB client
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatal("Failed to create cosmos db client: ", err)
	}

	err = createDatabase(client, databaseName)
	if err != nil {
		log.Printf("createDatabase failed: %s\n", err)
	}

	err = createContainer(client, databaseName, containerName, partitionKey)
	if err != nil {
		log.Printf("createContainer failed: %s\n", err)
	}

	err = createItem(client, databaseName, containerName, item.CountId, item)
	if err != nil {
		log.Printf("createItem failed: %s\n", err)
	}

	err = readItem(client, databaseName, containerName, item.CountId, item.ID)
	if err != nil {
		log.Printf("readItem failed: %s\n", err)
	}

	err = deleteItem(client, databaseName, containerName, item.CountId, item.ID)
	if err != nil {
		log.Printf("deleteItem failed: %s\n", err)
	}
}

func createDatabase(client *azcosmos.Client, databaseName string) error {
	//	databaseName := "adventureworks"

	databaseProperties := azcosmos.DatabaseProperties{ID: databaseName}

	// This is a helper function that swallows 409 errors
	errorIs409 := func(err error) bool {
		var responseErr *azcore.ResponseError
		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
	}
	ctx := context.TODO()
	databaseResp, err := client.CreateDatabase(ctx, databaseProperties, nil)

	switch {
	case errorIs409(err):
		log.Printf("Database [%s] already exists\n", databaseName)
	case err != nil:
		return err
	default:
		log.Printf("Database [%v] created. ActivityId %s\n", databaseName, databaseResp.ActivityID)
	}
	return nil
}

func createContainer(client *azcosmos.Client, databaseName, containerName, partitionKey string) error {
	//	databaseName = VisitorCounter
	//	containerName = Counter
	//	partitionKey = "/countId"

	databaseClient, err := client.NewDatabase(databaseName)
	if err != nil {
		return err
	}

	// creating a container
	containerProperties := azcosmos.ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{partitionKey},
		},
	}

	// this is a helper function that swallows 409 errors
	errorIs409 := func(err error) bool {
		var responseErr *azcore.ResponseError
		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
	}

	// setting options upon container creation
	throughputProperties := azcosmos.NewManualThroughputProperties(400) //defaults to 400 if not set
	options := &azcosmos.CreateContainerOptions{
		ThroughputProperties: &throughputProperties,
	}
	ctx := context.TODO()
	containerResponse, err := databaseClient.CreateContainer(ctx, containerProperties, options)

	switch {
	case errorIs409(err):
		log.Printf("Container [%s] already exists\n", containerName)
	case err != nil:
		return err
	default:
		log.Printf("Container [%s] created. ActivityId %s\n", containerName, containerResponse.ActivityID)
	}
	return nil
}

func createItem(client *azcosmos.Client, databaseName, containerName, partitionKey string, item any) error {
	//	databaseName = VisitorCounter
	//	containerName = Counter
	//	partitionKey = "/countId"

	/*	item = struct {
			ID           string `json:"id"`
			CountId      int    `json:"countId"`
			CreationDate string
		}{
			ID:           "1",
			CountId:      0,
			CreationDate: "2014-02-25T00:00:00",
		}
	*/
	// create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client: %s", err)
	}

	// specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	// setting the item options upon creating ie. consistency level
	itemOptions := azcosmos.ItemOptions{
		ConsistencyLevel: azcosmos.ConsistencyLevelSession.ToPtr(),
	}

	// this is a helper function that swallows 409 errors
	errorIs409 := func(err error) bool {
		var responseErr *azcore.ResponseError
		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
	}

	ctx := context.TODO()
	itemResponse, err := containerClient.CreateItem(ctx, pk, b, &itemOptions)

	switch {
	case errorIs409(err):
		log.Printf("Item with partitionkey value %s already exists\n", pk)
	case err != nil:
		return err
	default:
		log.Printf("Status %d. Item %v created. ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)
	}

	return nil
}

func readItem(client *azcosmos.Client, databaseName string, containerName string, partitionKey string, countId int) error {
	//	databaseName = VisitorCounter
	//	containerName = Counter
	//	partitionKey = "/countId"
	//	countId = "1"

	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client: %s", err)
	}

	// Specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Read an item
	ctx := context.TODO()
	itemResponse, err := containerClient.ReadItem(ctx, pk, countId, nil)
	if err != nil {
		return err
	}

	itemResponseBody := struct {
		ID           string `json:"id"`
		CountId      int
		CreationDate string
	}{}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(itemResponseBody, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("Read item with countId %s\n", itemResponseBody.CountId)
	fmt.Printf("%d\n", b)

	log.Printf("Status %d. Item %v read. ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)

	return nil
}

func deleteItem(client *azcosmos.Client, databaseName string, containerName string, partitionKey string, countId string) error {
	//	databaseName = VisitorCounter
	//	containerName = Counter
	//	partitionKey = "/countId"
	//	countId = "1"

	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client:: %s", err)
	}
	// Specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Delete an item
	ctx := context.TODO()

	res, err := containerClient.DeleteItem(ctx, pk, countId, nil)
	if err != nil {
		return err
	}

	log.Printf("Status %d. Item %v deleted. ActivityId %s. Consuming %v Request Units.\n", res.RawResponse.StatusCode, pk, res.ActivityID, res.RequestCharge)

	return nil
}
