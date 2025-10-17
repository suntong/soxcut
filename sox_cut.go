package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ============================ CONFIGURATION ===================================

const (
	// The input audio file (e.g., wav, mp3, flac).
	inputFile = "input.mp3"

	// The desired output file.
	outputFile = "output.mp3"

	// The file containing the start and end times for each clip.
	// Format per line: HH:MM:SS.mmm HH:MM:SS.mmm (e.g., 00:01:10.500 00:01:15.500)
	timingsFile = "test/segments.txt"

	// The duration of the cross-fade overlap.
	excessDuration = 500 * time.Millisecond

	// The search window for finding the best splice point.
	leewayDuration = 200 * time.Millisecond
)

// ========================== END OF CONFIGURATION ==============================

// ClipTiming holds the start and end time for a single audio segment.
type ClipTiming struct {
	Start time.Duration
	End   time.Duration
}

func main() {
	log.Println("Go Audio Splicer started.")

	// 1. Sanity Checks
	if !commandExists("sox") {
		log.Fatal("SoX not found in PATH. Please install it to continue.")
	}
	log.Println("Found SoX executable.")

	// 2. Read and parse the clip timings file.
	timings, err := parseTimingsFile(timingsFile)
	if err != nil {
		log.Fatalf("Error reading timings file '%s': %v", timingsFile, err)
	}
	if len(timings) == 0 {
		log.Fatal("No clip timings found in the file. Exiting.")
	}
	log.Printf("Found %d clip(s) to process.", len(timings))

	// 3. Create a temporary directory for intermediate files.
	tempDir, err := os.MkdirTemp("", "go_audio_splice_*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	log.Printf("Temporary directory created at: %s", tempDir)

	// 4. Extract and prepare all clips for splicing.
	preparedClipPaths, err := prepareClips(timings, tempDir)
	if err != nil {
		log.Fatalf("Failed during clip preparation: %v", err)
	}
	log.Println("All clips extracted and prepared successfully.")

	// 5. Splice the prepared clips together.
	finalClipPath, err := spliceClips(preparedClipPaths, tempDir)
	if err != nil {
		log.Fatalf("Failed during splicing: %v", err)
	}
	log.Println("All clips spliced successfully.")

	// 6. Move the final result to the output file.
	if err := os.Rename(finalClipPath, outputFile); err != nil {
		log.Fatalf("Failed to move final file to '%s': %v", outputFile, err)
	}

	log.Println("-----------------------------------")
	log.Printf("Processing complete! Final audio saved to: %s", outputFile)
}

// prepareClips loops through the timings, trimming each clip from the source
// with the correct excess/leeway for perfect splicing.
func prepareClips(timings []ClipTiming, tempDir string) ([]string, error) {
	var preparedClipPaths []string
	clipCount := len(timings)

	for i, timing := range timings {
		if timing.Start >= timing.End {
			return nil, fmt.Errorf("invalid timing for clip %d: start time is after end time", i+1)
		}

		idealDuration := timing.End - timing.Start
		var trimStart, trimDuration time.Duration
		clipPath := filepath.Join(tempDir, fmt.Sprintf("clip_%d_prep.wav", i))

		// Determine trim parameters based on clip position (first, middle, last).
		isFirst := (i == 0)
		isLast := (i == clipCount-1)

		switch {
		case isFirst && isLast: // Only one clip
			trimStart = timing.Start
			trimDuration = idealDuration
		case isFirst: // First clip of many
			trimStart = timing.Start
			trimDuration = idealDuration + excessDuration
		case isLast: // Last clip of many
			trimStart = timing.Start - (excessDuration + leewayDuration)
			trimDuration = idealDuration + excessDuration + leewayDuration
		default: // A middle clip
			trimStart = timing.Start - (excessDuration + leewayDuration)
			trimDuration = idealDuration + excessDuration + leewayDuration + excessDuration
		}

		if trimStart < 0 {
			log.Printf("Warning: Clip %d start time is too early for full leeway. Trimming from 0.", i+1)
			trimDuration += trimStart // Adjust duration since we start later.
			trimStart = 0
		}

		log.Printf(" -> Preparing clip %d: trimming from %.3fs for %.3fs", i+1, trimStart.Seconds(), trimDuration.Seconds())

		cmd := exec.Command("sox", inputFile, clipPath, "trim",
			fmt.Sprintf("%f", trimStart.Seconds()),
			fmt.Sprintf("%f", trimDuration.Seconds()),
		)
		if output, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("failed to trim clip %d: %v\nOutput: %s", i+1, err, string(output))
		}
		preparedClipPaths = append(preparedClipPaths, clipPath)
	}
	return preparedClipPaths, nil
}

// spliceClips iteratively joins the prepared clips using the splice effect.
func spliceClips(clipPaths []string, tempDir string) (string, error) {
	if len(clipPaths) <= 1 {
		return clipPaths[0], nil // Only one clip, no splicing needed.
	}

	currentCombinedFile := clipPaths[0]

	for i := 1; i < len(clipPaths); i++ {
		log.Printf(" -> Splicing clip %d onto the result...", i+1)
		nextClip := clipPaths[i]
		tempOutputFile := filepath.Join(tempDir, fmt.Sprintf("combined_%d.wav", i))

		// Get the duration of the current combined file to determine the splice position.
		// Per the man page, this is the duration of the first input file to the splice command.
		splicePos, err := getAudioDuration(currentCombinedFile)
		if err != nil {
			return "", fmt.Errorf("could not get duration of '%s': %v", currentCombinedFile, err)
		}

		spliceArgs := fmt.Sprintf("%f,%f,%f", splicePos.Seconds(), excessDuration.Seconds(), leewayDuration.Seconds())

		cmd := exec.Command("sox", currentCombinedFile, nextClip, tempOutputFile, "splice", "-q", spliceArgs)
		if output, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("failed to splice clip %d: %v\nOutput: %s", i+1, err, string(output))
		}
		currentCombinedFile = tempOutputFile
	}
	return currentCombinedFile, nil
}

// --- Helper Functions ---

// getAudioDuration uses `soxi` to get the precise duration of an audio file.
func getAudioDuration(filePath string) (time.Duration, error) {
	cmd := exec.Command("soxi", "-D", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("soxi command failed: %w: %s", err, string(output))
	}

	durationSec, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse soxi duration '%s': %w", string(output), err)
	}
	return time.Duration(durationSec * float64(time.Second)), nil
}

// parseTimingsFile reads the HH:MM:SS.mmm formatted file.
func parseTimingsFile(filePath string) ([]ClipTiming, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var timings []ClipTiming
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { // Skip empty lines and comments
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields (start and end time), got %d", lineNumber, len(parts))
		}

		start, err := parseSoxTime(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid start time format '%s': %w", lineNumber, parts[0], err)
		}
		end, err := parseSoxTime(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid end time format '%s': %w", lineNumber, parts[1], err)
		}
		timings = append(timings, ClipTiming{Start: start, End: end})
	}

	return timings, scanner.Err()
}

// parseSoxTime converts a HH:MM:SS.mmm string to a time.Duration.
func parseSoxTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format, expected HH:MM:SS.mmm")
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours: %w", err)
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	secParts := strings.Split(parts[2], ".")
	if len(secParts) > 2 {
		return 0, fmt.Errorf("invalid seconds format")
	}
	sec, err := strconv.Atoi(secParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	ms := 0
	if len(secParts) == 2 {
		// Pad with zeros to handle cases like .5 -> 500ms
		msStr := secParts[1]
		for len(msStr) < 3 {
			msStr += "0"
		}
		ms, err = strconv.Atoi(msStr[:3]) // Take only first 3 digits
		if err != nil {
			return 0, fmt.Errorf("invalid milliseconds: %w", err)
		}
	}

	return time.Hour*time.Duration(h) +
		time.Minute*time.Duration(m) +
		time.Second*time.Duration(sec) +
		time.Millisecond*time.Duration(ms), nil
}

// commandExists checks if a command is available in the system's PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
