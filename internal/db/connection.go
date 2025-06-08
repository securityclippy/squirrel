package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectWithRetry attempts to connect to the database with exponential backoff
// It will retry for up to maxDuration (5 minutes) before giving up
func ConnectWithRetry(ctx context.Context, databaseURL string, maxDuration time.Duration) (*pgxpool.Pool, error) {
	log.Printf("Attempting to connect to database...")
	
	var pool *pgxpool.Pool
	var lastErr error
	
	startTime := time.Now()
	attempt := 0
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second
	
	for time.Since(startTime) < maxDuration {
		attempt++
		
		// Create connection pool
		var err error
		pool, err = pgxpool.New(ctx, databaseURL)
		if err != nil {
			lastErr = fmt.Errorf("failed to create connection pool (attempt %d): %w", attempt, err)
			log.Printf("%v", lastErr)
		} else {
			// Test the connection
			if err = pool.Ping(ctx); err != nil {
				lastErr = fmt.Errorf("failed to ping database (attempt %d): %w", attempt, err)
				log.Printf("%v", lastErr)
				pool.Close()
			} else {
				// Success!
				log.Printf("Connected to database successfully after %d attempts in %v", 
					attempt, time.Since(startTime).Round(time.Millisecond))
				return pool, nil
			}
		}
		
		// Calculate delay with exponential backoff
		delay := time.Duration(attempt) * baseDelay
		if delay > maxDelay {
			delay = maxDelay
		}
		
		// Check if we have time for another attempt
		if time.Since(startTime)+delay >= maxDuration {
			break
		}
		
		log.Printf("Retrying in %v... (elapsed: %v)", delay, time.Since(startTime).Round(time.Millisecond))
		
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
	
	return nil, fmt.Errorf("failed to connect to database after %v: %w", maxDuration, lastErr)
}