package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/axarus/vectrag/internal/application"
	infrahttp "github.com/axarus/vectrag/internal/infrastructure/http"
	"github.com/spf13/cobra"
)

var developCmd = &cobra.Command{
	Use:   "develop",
	Short: "Run VectraG in development mode",
	Long: `The develop command starts VectraG with development-friendly settings and tools, allowing you to:

- Create, update, and delete content models
- Experiment with schema changes safely
- Iterate quickly during local development

This mode is intended for local environments and early-stage development.`,
	Run: func(cmd *cobra.Command, args []string) {
		basePort := 51987
		host := "localhost"

		svc := application.NewDevelopService(
			infrahttp.ListenerProvider{},
			infrahttp.ServerStarter{},
			infrahttp.AdminHandlerProvider{},
			infrahttp.APIRoutesProvider{},
		)

		url, shutdown, err := svc.Start(basePort, host)
		if err != nil {
			fmt.Println("Error starting server:", err)
			return
		}
		fmt.Println("Server started at", url)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("Shutting down server...")
		_ = shutdown(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(developCmd)
}
