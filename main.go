package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "webhook",
	}

	cmd.AddCommand(serverCommand())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
