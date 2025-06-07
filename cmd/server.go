package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"reminder-service/internal/handlers"
	"reminder-service/internal/services"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the reminder service server",
	Long:  "Start the HTTP server for the reminder service with API endpoints",
	Run:   runServer,
}

var (
	port string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
}

func runServer(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	
	// Database connection
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/reminder_service?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Initialize services
	reminderService := services.NewReminderService(pool)

	// Initialize handlers
	reminderHandler := handlers.NewReminderHandler(reminderService)

	// Setup router
	router := mux.NewRouter()

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Register reminder routes
	reminderHandler.RegisterRoutes(router)

	// Use environment variable if port flag not explicitly set
	if !cmd.Flags().Changed("port") {
		if envPort := os.Getenv("PORT"); envPort != "" {
			port = envPort
		}
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}