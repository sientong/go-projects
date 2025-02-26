package analytics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type AccessLog struct {
	ShortURL   string
	RemoteAddr string
	UserAgent  string
	AccessTime time.Time
}

var accessChannel = make(chan AccessLog, 100)

func LogAccess(log AccessLog) {
	accessChannel <- log
}

func StartAccessLogger(ctx context.Context) {
	log.Println("Access logger started") // Debug log
	go func() {
		for {
			select {
			case access := <-accessChannel:
				log.Println("Received access log:", access) // Debug log
				// Store access log in Redis
				err := storeAccessLog(ctx, access)
				if err != nil {
					log.Printf("Error storing access log: %v", err)
				}
			case <-ctx.Done():
				log.Println("Access logger stopped") // Debug log
				return
			}
		}
	}()
}

func storeAccessLog(ctx context.Context, access AccessLog) error {
	// Ensure the log directory exists
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Printf("Error creating log directory: %v", err) // Debug log
		return err
	}

	file, err := os.OpenFile("log/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err) // Debug log
		return err
	}
	defer file.Close()

	logEntry := fmt.Sprintf("%s - [%s] \"%s\" \"%s\"\n", access.RemoteAddr, access.AccessTime.Format(time.RFC1123), access.ShortURL, access.UserAgent)
	_, err = file.WriteString(logEntry)
	if err != nil {
		log.Printf("Error writing to log file: %v", err) // Debug log
		return err
	}

	return nil
}

// LoggingMiddleware logs the incoming HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.RequestURI, time.Since(start))

		// Add access log to the channel
		accessLog := AccessLog{
			ShortURL:   r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			AccessTime: time.Now(),
		}
		accessChannel <- accessLog
	})
}
