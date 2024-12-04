package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	ffmpegProcess *exec.Cmd
	mu            sync.Mutex
)

// Supported resolutions and frame rates
var supportedResolutions = map[string]bool{
	"640x480":   true,
	"800x600":   true,
	"1024x768":  true,
	"1280x720":  true,
	"1280x800":  true,
	"1366x768":  true,
	"1600x900":  true,
	"1920x1080": true,
	"1920x1200": true,
	"2560x1440": true,
	"2560x1600": true,
	"3840x2160": true,
	"5120x2880": true,
	"7680x4320": true,
}

var supportedFrameRates = map[string]bool{
	"24":    true,
	"25":    true,
	"30":    true,
	"48":    true,
	"50":    true,
	"59":    true,
	"59.94": true,
	"60":    true,
	"120":   true,
	"240":   true,
	"15":    true,
	"29.97": true,
}

// Function to convert fraction to decimal
func convertFractionToDecimal(fraction string) (string, error) {
	parts := strings.Split(fraction, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("Internal Error : Invalid Fraction Format")
	}

	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", err
	}

	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}

	// Calculate the decimal value
	decimal := float64(numerator) / float64(denominator)
	return fmt.Sprintf("%.2f", decimal), nil
}

func formatFrameRate(frameRate string) string {
	// Parse the frame rate string to a float
	parsedFrameRate, err := strconv.ParseFloat(frameRate, 64)
	if err != nil {
		// Return the original frame rate string if parsing fails
		return frameRate
	}

	// If the frame rate is a whole number, return it as an integer string
	if parsedFrameRate == float64(int(parsedFrameRate)) {
		return fmt.Sprintf("%.0f", parsedFrameRate)
	}

	// Otherwise, return the original frame rate string (including decimals)
	return frameRate
}

// Function to get resolution and frame rate of a video using ffprobe
func getVideoProbe(videoFile string) (string, string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", "./videos/InputVideos/"+videoFile)
	output, err := cmd.Output()
	var properFrameRate string
	if err != nil {
		return "", "", err
	}

	// Split the output by lines
	outputLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(outputLines) < 3 {
		return "", "", fmt.Errorf("Invalid Video")
	}

	// Resolution (width and height)
	width := outputLines[0]
	height := outputLines[1]

	// Frame rate
	frameRate := outputLines[2]

	// Check if the resolution and frame rate are supported
	resolution := width + "x" + height
	if _, exists := supportedResolutions[resolution]; !exists {
		return "", "", fmt.Errorf("Unsupported Resolution: %s", resolution)
	}

	// Normalize frame rate to remove decimals (e.g., 25.00 becomes 25)
	normalizedFrameRate, err := convertFractionToDecimal(frameRate)
	if err != nil {
		return "", "", err
	}
	properFrameRate = formatFrameRate(normalizedFrameRate)
	if _, exists := supportedFrameRates[properFrameRate]; !exists {
		return "", "", fmt.Errorf("Unsupported Frame Rate: %s", normalizedFrameRate)
	}

	return resolution, properFrameRate, nil
}

// Function to get the number of audio tracks in a video file
func getAudioTrackCount(videoFile string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a", "-show_entries", "stream=codec_type", "-of", "csv=p=0", "./videos/InputVideos/"+videoFile)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Count the number of audio streams by counting the lines in the output
	outputLines := strings.Split(string(output), "\n")
	audioTrackCount := (len(outputLines) - 2) / 2 // Subtract 1 to ignore the empty last line

	return fmt.Sprintf("%d", audioTrackCount), nil
}

// Function to get video files along with their resolution and frame rate
func getVideoFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./videos/InputVideos") // Change to your directory where videos are stored
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var videoFiles []map[string]string
	for _, file := range files {
		if !file.IsDir() {
			if file.Name() == "README.md" {
				continue
			}

			// Get video resolution and frame rate
			resolution, frameRate, err := getVideoProbe(file.Name())
			if err != nil {
				// Log error and continue to the next file if invalid resolution/frame rate
				fmt.Printf("VideoFile '%s' Issue\nError probing video: %s\n\n", file.Name(), err)
				continue
			}

			// Get the number of audio tracks
			audioCount, err := getAudioTrackCount(file.Name())
			if err != nil {
				// Log error and continue to the next file if there is an issue probing audio tracks
				fmt.Printf("VideoFile '%s' Issue\nError probing audio: %s\n\n", file.Name(), err)
				audioCount = "0" // Set to 0 in case of error
			}

			// Add video info to the response slice if resolution and frame rate are valid
			videoFiles = append(videoFiles, map[string]string{
				"fileName":    file.Name(),
				"resolution":  resolution,
				"frameRate":   frameRate,
				"audioTracks": audioCount,
			})
		}
	}

	// Set Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Send the response as JSON
	if err := json.NewEncoder(w).Encode(videoFiles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func openVLC(w http.ResponseWriter, r *http.Request) {
	// Run the VLC command
	cmd := exec.Command("vlc", "udp://@127.0.0.1:2430") // Adjust the command depending on your OS (e.g., "vlc", "vlc.exe" for Windows)
	err := cmd.Start()
	if err != nil {
		// If there's an error starting the command, send a failure response
		http.Error(w, fmt.Sprintf("Failed to execute command: %v", err), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("TekMedia's MPEG-TS Generator")
	fmt.Println("Copyright (c) 2024 TekMedia Software")
	fmt.Println("\nOpen http://127.0.0.1:8080 in your browser...")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/videos", getVideoFiles)
	http.HandleFunc("/open-vlc", openVLC)
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ffmpegArgs := []string{"-y"}
			videoinput := r.FormValue("videoinput")
			signal := r.FormValue("signal")
			resolution := r.FormValue("resolution")
			framerate := r.FormValue("framerate")
			pixelFormat := r.FormValue("pixelFormat")
			scanType := r.FormValue("scanType")
			videoCodec := r.FormValue("videoCodec")
			videoBitrate := r.FormValue("videoBitrate")
			streaming := r.FormValue("stream")
			destAddress := r.FormValue("destIp")
			port := r.FormValue("destPort")
			fileName := r.FormValue("fileName")
			audioEnabled := r.FormValue("enableAudio")
			audioTrack := r.FormValue("audioTracks")
			externalAudioInput := r.FormValue("audioinput")
			audioLayout := r.Form["audioLayout[]"]
			audioCodec := r.Form["audioCodec[]"]
			audioBitrate := r.Form["audioBitrate[]"]
			var audioValue, params string
			audioTracks, err := strconv.Atoi(audioTrack)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				// return
			}

			videoBitrate += "000000"
			bufSize, err := strconv.Atoi(videoBitrate)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				// return
			}
			bufSize *= 2
			bufsize := strconv.Itoa(bufSize)

			// Print the values for debugging
			// fmt.Println("videoinput:", videoinput)
			// fmt.Println("signal:", signal)
			// fmt.Println("resolution:", resolution)
			// fmt.Println("framerate:", framerate)
			// fmt.Println("pixelFormat:", pixelFormat)
			// fmt.Println("scanType:", scanType)
			// fmt.Println("videoCodec:", videoCodec)
			// fmt.Println("videoBitrate:", videoBitrate)
			// fmt.Println("streaming:", streaming)
			// fmt.Println("destAddress:", destAddress)
			// fmt.Println("port:", port)
			// fmt.Println("fileName:", fileName)
			// fmt.Println("audio:", audioEnabled)
			// fmt.Println("audiotr:", audioTracks)
			// fmt.Println("audioinput:", externalAudioInput)

			if videoinput == "true" {
				ffmpegArgs = append(ffmpegArgs, "-stream_loop", "-1")
				ffmpegArgs = append(ffmpegArgs, "-i", fmt.Sprintf("./videos/InputVideos/%s", signal))
			} else {
				ffmpegArgs = append(ffmpegArgs, "-re")
				ffmpegArgs = append(ffmpegArgs, "-fflags", "+genpts")
				ffmpegArgs = append(ffmpegArgs, "-f", "lavfi")
				ffmpegArgs = append(ffmpegArgs, "-i", fmt.Sprintf("%s=size=%s:rate=%s", signal, resolution, framerate))
			}

			if audioEnabled == "true" && externalAudioInput != "true" {
				for i := 0; i < audioTracks; i++ {
					audioValue += fmt.Sprintf("[1]aformat=channel_layouts=%s[a%d];", audioLayout[i], i+1)
				}
				ffmpegArgs = append(ffmpegArgs, "-f", "lavfi")
				ffmpegArgs = append(ffmpegArgs, "-i", "sine=frequency=1000")
				ffmpegArgs = append(ffmpegArgs, "-filter_complex", audioValue[:len(audioValue)-1])
			}

			ffmpegArgs = append(ffmpegArgs, "-map", "0:v")
			for i := 0; i < audioTracks; i++ {
				if externalAudioInput != "true" {
					ffmpegArgs = append(ffmpegArgs, "-map", fmt.Sprintf("[a%d]", i+1))
				} else {
					ffmpegArgs = append(ffmpegArgs, "-map", fmt.Sprintf("0:a:%d", i))
				}
			}
			ffmpegArgs = append(ffmpegArgs, "-pix_fmt", pixelFormat)

			if videoCodec == "avc" {
				ffmpegArgs = append(ffmpegArgs, "-c:v", "libx264")
				params = "-x264-params"
			} else if videoCodec == "hevc" {
				ffmpegArgs = append(ffmpegArgs, "-c:v", "libx265")
				params = "-x265-params"
			}

			ffmpegArgs = append(ffmpegArgs, "-b:v", videoBitrate)
			ffmpegArgs = append(ffmpegArgs, "-minrate", videoBitrate)
			ffmpegArgs = append(ffmpegArgs, "-maxrate", videoBitrate)
			ffmpegArgs = append(ffmpegArgs, "-bufsize", bufsize)

			if scanType == "ilc" {
				ffmpegArgs = append(ffmpegArgs, "-vf", fmt.Sprintf("format=%s,interlace", pixelFormat))
				ffmpegArgs = append(ffmpegArgs, "-flags", "+ilme+ildct")
				ffmpegArgs = append(ffmpegArgs, "-top", "1")
			}

			ffmpegArgs = append(ffmpegArgs, params, "nal-hrd=cbr")
			ffmpegArgs = append(ffmpegArgs, "-g", "150")
			ffmpegArgs = append(ffmpegArgs, "-bf", "3")

			for i := 0; i < audioTracks; i++ {
				ffmpegArgs = append(ffmpegArgs, fmt.Sprintf("-c:a:%d", i), audioCodec[i])
				ffmpegArgs = append(ffmpegArgs, fmt.Sprintf("-b:a:%d", i), audioBitrate[i])
			}

			ffmpegArgs = append(ffmpegArgs, "-f", "tee")

			output := "[f=mpegts]udp://127.0.0.1:2430?pkt_size=1316|[f=mpegts]"

			if streaming == "true" {
				output += fmt.Sprintf("udp://%s:%s?pkt_size=1316", destAddress, port)
			} else {
				output += fmt.Sprintf("./videos/GeneratedVideos/%s.ts", fileName)
			}

			ffmpegArgs = append(ffmpegArgs, output)

			fmt.Println("\nCommand Executed:")
			fmt.Println("ffmpeg", strings.Join(ffmpegArgs, " "))

			// Prepare the FFmpeg command
			ffmpegProcess = exec.Command("ffmpeg", ffmpegArgs...)

			// Redirect stdout and stderr to the log file
			logFile, err := os.OpenFile("ffmpeg.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				http.Error(w, "Could not open log file", http.StatusInternalServerError)
				return
			}
			defer logFile.Close()
			ffmpegProcess.Stdout = logFile
			ffmpegProcess.Stderr = logFile

			// Run the FFmpeg command
			err = ffmpegProcess.Start() // Use Start() instead of Run()
			if err != nil {
				ffmpegProcess = nil // Reset on failure
				http.Error(w, "Failed to start streaming", http.StatusInternalServerError)
				return
			}

			fmt.Fprintln(w, "Streaming started successfully.")
		}
	})

	// Handle stop streaming
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock() // Lock before accessing ffmpegProcess
		defer mu.Unlock()

		if ffmpegProcess != nil {
			err := ffmpegProcess.Process.Kill()
			if err != nil {
				http.Error(w, "Failed to stop streaming: "+err.Error(), http.StatusInternalServerError)
				return
			}
			ffmpegProcess = nil
			fmt.Fprintln(w, "Streaming stopped successfully.")
		} else {
			http.Error(w, "No streaming process found", http.StatusBadRequest)
		}
	})

	// Start the HTTP server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// Check if the error is related to the port already being in use
		if opErr, ok := err.(*os.SyscallError); ok {
			// Check if it's an address already in use error (port busy)
			if opErr.Err == syscall.EADDRINUSE {
				fmt.Println("\nPort 8080 is already in use. Please try a different port.")
			} else {
				// Handle other system errors
				//log.Fatalf("Error starting server: %v", opErr)
				fmt.Println("\nError starting server")
			}
		} else {
			// Handle other errors
			//log.Fatalf("Error starting server: %v", err)
			fmt.Println("\nError starting server. Ensure no other instances are running.")
		}
	} else {
		fmt.Println("Server is running on port 8080")
	}
}
