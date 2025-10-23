package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/hs-heilbronn-devsecops/acetlisto/handlers"
	"github.com/hs-heilbronn-devsecops/acetlisto/stores"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")

	store := stores.NewMemoryItemStore()
	r := handlers.New(store)

	port := viper.GetString("PORT")
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), gorillahandlers.LoggingHandler(os.Stdout, r)))
}