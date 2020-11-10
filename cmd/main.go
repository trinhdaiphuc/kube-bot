package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/infracloudio/msbotbuilder-go/core"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/joho/godotenv"
)

var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
		return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
	},
}

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	core.Adapter
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	activity, err := ht.Adapter.ParseRequest(ctx, req)
	if err != nil {
		fmt.Println("Failed to parse request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ht.Adapter.ProcessActivity(ctx, activity, customHandler)
	if err != nil {
		fmt.Println("Failed to process request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Request processed successfully.")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	setting := core.AdapterSetting{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	fmt.Printf("Setting %+v\n", setting)

	adapter, err := core.NewBotAdapter(setting)
	if err != nil {
		log.Fatal("Error creating adapter: ", err)
	}

	httpHandler := &HTTPHandler{adapter}
	port := os.Getenv("PORT")
	http.HandleFunc("/api/messages", httpHandler.processMessage)
	fmt.Println("Starting server on port", port)
	http.ListenAndServe(":"+port, nil)
}
