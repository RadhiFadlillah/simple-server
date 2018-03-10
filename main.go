package main

import (
	"fmt"
	"log"
	"net/http"
	fp "path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "simple-server directory",
		Short: "Run a simple HTTP server to serve content of specified directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get root directory
			rootDir := args[0]

			// Create router
			router := httprouter.New()
			router.GET("/*filename", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				filename := ps.ByName("filename")
				filename = fp.Join(rootDir, filename)
				http.ServeFile(w, r, filename)
			})

			router.PanicHandler = func(w http.ResponseWriter, r *http.Request, arg interface{}) {
				http.Error(w, fmt.Sprint(arg), 500)
			}

			// Serve app
			port, _ := cmd.Flags().GetInt("port")
			addr := fmt.Sprintf(":%d", port)
			log.Println("Server running in port", port)
			log.Fatalln(http.ListenAndServe(addr, router))
		},
	}

	root.Flags().IntP("port", "p", 9000, "Port that used by server")
	root.Execute()
}
