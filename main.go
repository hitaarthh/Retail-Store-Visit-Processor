// main.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// JobStatus represents the current state of a job
type JobStatus string

const (
	StatusOngoing   JobStatus = "ongoing"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

// JobRequest represents the incoming job submission
type JobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

// Visit represents a store visit with images
type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

// Store represents store information from Store Master
type Store struct {
	ID       string `json:"store_id"`
	Name     string `json:"store_name"`
	AreaCode string `json:"area_code"`
}

// Job represents a processing job with updated response format
type Job struct {
	Status JobStatus `json:"status"`
	JobID  string    `json:"job_id"`
	Errors []Error   `json:"error,omitempty"`
}

// Error represents a processing error
type Error struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

// JobManager handles job storage and processing
type JobManager struct {
	jobs    map[int64]*Job
	mu      sync.RWMutex
	counter int64
	stores  map[string]Store // Store Master data
}

func loadStoresFromCSV(filename string) (map[string]Store, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	stores := make(map[string]Store)
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		stores[record[2]] = Store{ // record[2] is StoreID
			AreaCode: record[0],
			Name:     record[1],
			ID:       record[2],
		}
	}

	return stores, nil
}

func NewJobManager(csvPath string) (*JobManager, error) {
	stores, err := loadStoresFromCSV(csvPath)
	if err != nil {
		return nil, err
	}

	return &JobManager{
		jobs:   make(map[int64]*Job),
		stores: stores,
	}, nil
}

func (jm *JobManager) CreateJob() int64 {
	id := atomic.AddInt64(&jm.counter, 1)
	job := &Job{
		Status: StatusOngoing,
		JobID:  "",
	}

	jm.mu.Lock()
	jm.jobs[id] = job
	jm.mu.Unlock()

	return id
}

func (jm *JobManager) GetJob(id int64) (*Job, bool) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()
	job, exists := jm.jobs[id]
	return job, exists
}

func (jm *JobManager) ValidateStore(storeID string) bool {
	_, exists := jm.stores[storeID]
	return exists
}

func (jm *JobManager) ProcessJob(jobID int64, visits []Visit) {
	job, exists := jm.GetJob(jobID)
	if !exists {
		return
	}

	job.Status = StatusOngoing

	for _, visit := range visits {
		// Validate if store exists
		if !jm.ValidateStore(visit.StoreID) {
			job.Status = StatusFailed
			job.Errors = append(job.Errors, Error{
				StoreID: visit.StoreID,
				Error:   fmt.Sprintf("Store ID %s does not exist in the store master data.", visit.StoreID),
			})
			continue // Continue with other visits if a store fails validation
		}

		// Process each image
		for _, url := range visit.ImageURLs {
			if err := processImage(url); err != nil {
				job.Status = StatusFailed
				job.Errors = append(job.Errors, Error{
					StoreID: visit.StoreID,
					Error:   fmt.Sprintf("Failed to download or process image from URL %s: %v", url, err),
				})
				continue // Continue with other images if one fails
			}
		}
	}

	// If no errors were encountered, mark job as completed
	if len(job.Errors) == 0 {
		job.Status = StatusCompleted
	} else {
		job.Status = StatusFailed // If there were errors, mark job as failed
	}
}

func serveFrontend() http.Handler {
	frontendPath := "./frontend" // Path to your frontend files
	return http.FileServer(http.Dir(frontendPath))
}

// Update main.go
func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Initialize job manager with CSV
	manager, err := NewJobManager("StoreMasterAssignment.csv")
	if err != nil {
		log.Fatalf("Failed to initialize job manager: %v", err)
	}
	jobManager = manager

	http.Handle("/", serveFrontend())

	// Set up routes
	http.HandleFunc("/api/submit", handleSubmitJob)
	http.HandleFunc("/api/status", handleJobStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Println()
	fmt.Printf("1. Open your browser.\n")
	fmt.Printf("2. Enter http://localhost:%s/ in the address bar.\n", port)
	fmt.Println("   (You should see the frontend interface for Store Management and job status checking.)")
	fmt.Println("3. If the site is inactive or you prefer, use Postman/Insomnia to test the backend endpoints:")
	fmt.Printf("   - POST http://localhost:%s/api/submit\n", port)
	fmt.Printf("   - GET http://localhost:%s/api/status?jobid=<job_id>\n", port)
	fmt.Println("     Replace <job_id> with a valid job ID to check its status.")
	fmt.Println()

	// Check if the port is available
	fmt.Printf("Checking if port %s is available...\n", port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error: Port %s is already in use. Please free the port or use a different one.\n", port)
		os.Exit(1)
	}
	ln.Close() // Port is free; close the listener before starting the server
	fmt.Printf("Port %s is available. Starting server...\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func processImage(url string) error {
	// Download image
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	// Calculate perimeter
	bounds := img.Bounds()
	_ = float64(2 * (bounds.Max.X - bounds.Min.X + bounds.Max.Y - bounds.Min.Y))

	// Random sleep (0.1 to 0.4 seconds)
	sleepTime := time.Duration(100+rand.Intn(300)) * time.Millisecond
	time.Sleep(sleepTime)

	return nil
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Updated handleSubmitJob function
func handleSubmitJob(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	if r.Method == http.MethodOptions {
		return // Handle preflight request and return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jobReq JobRequest
	if err := json.NewDecoder(r.Body).Decode(&jobReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON format. Please check the input format.",
		})
		return
	}

	// Validate that "visits" field is provided
	if len(jobReq.Visits) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing 'visits' field. At least one visit is required.",
		})
		return
	}

	// Validate each visit field
	for i, visit := range jobReq.Visits {
		if visit.StoreID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Missing 'store_id' in visit at index %d.", i),
			})
			return
		}

		if len(visit.ImageURLs) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Missing 'image_url' array in visit for store %s.", visit.StoreID),
			})
			return
		}

		if visit.VisitTime == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Missing 'visit_time' in visit for store %s.", visit.StoreID),
			})
			return
		}
	}

	// Validate count matches the number of visits
	if jobReq.Count != len(jobReq.Visits) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("'count' mismatch: expected %d, got %d.", jobReq.Count, len(jobReq.Visits)),
		})
		return
	}

	// Create a new job and start processing
	jobID := jobManager.CreateJob()
	go jobManager.ProcessJob(jobID, jobReq.Visits)

	// Send the job ID as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": fmt.Sprintf("%d", jobID),
	})
}

// Updated handleJobStatus function
func handleJobStatus(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	if r.Method == http.MethodOptions {
		return // Handle preflight request and return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract job ID from query parameters
	jobIDStr := r.URL.Query().Get("jobid")
	if jobIDStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing 'jobid' parameter in query.",
		})
		return
	}

	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid job ID format: %s", jobIDStr),
		})
		return
	}

	// Retrieve job status from the JobManager
	job, exists := jobManager.GetJob(jobID)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Job ID %s not found.", jobIDStr),
		})
		return
	}

	// Structure the response based on job status
	response := struct {
		Status JobStatus `json:"status"`
		JobID  string    `json:"job_id"`
		Errors []Error   `json:"error,omitempty"`
	}{
		Status: job.Status,
		JobID:  jobIDStr,
	}

	// Include errors if the job failed
	if job.Status == StatusFailed {
		response.Errors = job.Errors
	}

	// Send the response in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var jobManager *JobManager
