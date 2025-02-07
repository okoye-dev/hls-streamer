package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	videoPath = "./test/tanjiro.mp4"
	outputPath = "./output/output.m3u8"
)

func convertVideo() error {
	// Run FFmpeg script
	cmd := exec.Command( "/bin/sh", "ffmpeg.sh", videoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func streamVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if the output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		fmt.Println("Processing video conversion...")

		// Convert the video
		if err := convertVideo(); err != nil {
			http.Error(w, "Error processing video", http.StatusInternalServerError)
			return
		}
	}

	http.ServeFile(w, r, outputPath)
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "server is up"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Serve video files
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/output/", http.StripPrefix("/output/", fs))

	// Handle processing request
	http.HandleFunc("/test/convert", streamVideo)
	http.HandleFunc("/test", healthCheck)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}