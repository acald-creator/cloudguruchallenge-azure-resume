package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type CounterEntity struct {
	aztables.Entity
	Id    string
	Count int32
}

func getClient() *aztables.Client {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT environment variable not found")
	}

	tableName, ok := os.LookupEnv("AZURE_TABLE_NAME")
	if !ok {
		panic("AZURE_TABLE_NAME environment variable not found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, tableName)
	client, err := aztables.NewClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}
	return client
}

func addEntity(client *aztables.Client) {
	myEntity := CounterEntity{
		Entity: aztables.Entity{
			PartitionKey: "1",
		},
		Count: 1,
		Id:    "Count",
	}

	marshalled, err := json.Marshal(myEntity)
	if err != nil {
		panic(err)
	}

	_, err = client.AddEntity(context.TODO(), marshalled, nil)
	if err != nil {
		panic(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
	}
	fmt.Fprint(w, message)
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/GetVisitorCounter", helloHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))

	fmt.Println("Authenticating...")
	client := getClient()

	fmt.Println("Adding an entity to the table...")
	addEntity(client)

}
