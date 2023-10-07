package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const (
	webhookPattern = "/hook/"
)

type server struct {
	port int
}

func serverCommand() *cobra.Command {
	s := &server{}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the webhook server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.Validate(); err != nil {
				return err
			}
			return s.Serve()
		},
	}

	cmd.Flags().IntVarP(&s.port, "port", "p", 8080, "port to serve the webhook server")

	return cmd
}

func (s *server) Validate() error {
	if s.port < 0 || s.port > 65535 {
		return fmt.Errorf("port %d is invalid, must be between 1-65535", s.port)
	}
	return nil
}

func (s *server) Serve() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	m := http.NewServeMux()
	// TODO: Include metrics for received audit log events as well.
	m.Handle(webhookPattern, s.webhookHandler())

	srv := &http.Server{Addr: fmt.Sprintf(":%d", s.port), Handler: m}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed running server: %v", err)
		}
	}()

	log.Println("Started webhook server")
	<-done
	log.Println("Shutting down webhook server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutting down server gracefully failed: %w", err)
	}

	log.Println("Webhook server shutdown finished, exiting now...")

	return nil
}

func (s *server) webhookHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Body == nil {
			writeResponse(writer, http.StatusBadRequest, "No body received from web hook request")
			return
		}

		// TODO: Use v1.AuditMessage to create metrics based on the interaction, method etc.
		dec := json.NewDecoder(request.Body)
		var unstructured map[string]interface{}
		if err := dec.Decode(&unstructured); err != nil {
			writeResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		prettyJSON, err := json.MarshalIndent(unstructured, "", "    ")
		if err != nil {
			writeResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		log.Printf("Received webhook message from hook path %q:\n%s\n", request.URL.Path, prettyJSON)

		writer.WriteHeader(http.StatusOK)
	}
}

func writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(msg))
}
