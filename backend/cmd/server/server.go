package server

import (
	"ariadne/api"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

const defaultPort = "8080"

var port string

var ServerCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start API server",
	Example: "ariadne server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", defaultPort, "Port to run server on")
}

func start() {
	router := api.NewRouter()

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}