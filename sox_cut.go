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
	// The duration of the cross-fade overlap.
	excessDuration = 500 * time.Millisecond

	// The search window for finding the best splice point.
	leewayDuration = 200 * time.Millisecond
)

var (
	// The file containing the start and end times for each clip.
	// Format per line: HH:MM:SS.mmm HH:MM:SS.mmm (e.g., 00:01:10 00:01:15.6)
	timingsFile = "test/segments.txt"

	inputFile  string
	outputFile string
)

// ========================== END OF CONFIGURATION ==============================

// ClipTiming holds the start and end time for a single audio segment.
type ClipTiming struct {
	Start time.Duration
	End   time.Duration
}

func main() {
	log.Println("Audio Splicer started.")

	// Parse Command-Line Arguments
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go <inputFile> <outputFile> [timingsFile] [sox_options...]")
		fmt.Println("\nExample (WAV to MP3):")
		fmt.Println("  go run main.go input.wav output.mp3 timings.txt -C 192")
		fmt.Println("\nExample (WAV to Opus):")
		fmt.Println("  go run main.go audio.flac final.opus timings.txt -C 128k")
		os.Exit(0)
	}

	inputFile = args[0]
	outputFile = args[1]
	//timingsFilePath := defaultTimingsFile
	var soxOptions []string

	if len(args) > 2 {
		// Check if the third argument is a potential timings file or a sox option
		if !strings.HasPrefix(args[2], "-") {
			timingsFile = args[2]
			if len(args) > 3 {
				soxOptions = args[3:]
			}
		} else {
			soxOptions = args[2:]
		}
	}

	// Dependency Check: Ensure sox is installed.
	if !commandExists("sox") {
		log.Fatal("SoX not found in PATH. Please install it to continue.")
	}
	//log.Println("Found SoX executable.")

	// Read and parse the clip timings file.
	timings, err := parseTimingsFile(timingsFile)
	if err != nil {
		log.Fatalf("Error reading timings file '%s': %v", timingsFile, err)
	}
	if len(timings) == 0 {
		log.Fatal("No clip timings found in the file. Exiting.")
	}
	log.Printf("Found %d clip(s) to process from '%s'.", len(timings), timingsFile)

	// Create a temporary directory for intermediate files.
	tempDir, err := os.MkdirTemp("", "sc_*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	log.Printf("Temporary directory created at: %s", tempDir)

	// Extract and prepare all clips for splicing.
	preparedClipPaths, err := prepareClips(timings, tempDir)
	if err != nil {
		log.Fatalf("Failed during clip preparation: %v", err)
	}
	log.Println("All clips extracted and prepared successfully.")

	// Splice the prepared clips together.
	finalClipPath, err := spliceClips(preparedClipPaths, timings, tempDir)
	if err != nil {
		log.Fatalf("Failed during splicing: %v", err)
	}
	log.Println("All clips spliced successfully.")

	// Perform final encode to the output file, with user options.
	//   sox <input> <output> <options>
	finalCmdArgs := []string{finalClipPath, outputFile}
	finalCmdArgs = append(finalCmdArgs, soxOptions...)
	log.Printf("Encoding final file with\n\t\t '%v'...", finalCmdArgs)

	cmd := exec.Command("sox", finalCmdArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to execute final sox command: %v\nOutput: %s", err, string(output))
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

		log.Printf(" -> Preparing clip %d: trimming from %v for %.3fs (%.3fs)",
			i+1, trimStart, idealDuration.Seconds(), trimDuration.Seconds())

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
func spliceClips(clipPaths []string, timings []ClipTiming, tempDir string) (string, error) {
	if len(clipPaths) <= 1 {
		return clipPaths[0], nil // Only one clip, no splicing needed.
	}

	currentCombinedFile := clipPaths[0]
	accumulatedIdealDuration := timings[0].End - timings[0].Start

	for i := 1; i < len(clipPaths); i++ {
		log.Printf(" -> Splicing clip %d at joint point: %.3fs", i+1, accumulatedIdealDuration.Seconds())
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
		accumulatedIdealDuration += (timings[i].End - timings[i].Start)
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

		start, err := parseISOTime(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid start time format '%s': %w", lineNumber, parts[0], err)
		}
		end, err := parseISOTime(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid end time format '%s': %w", lineNumber, parts[1], err)
		}
		timings = append(timings, ClipTiming{Start: start, End: end})
	}

	return timings, scanner.Err()
}

var durationFormat = []string{"", "05", "04:05", "15:04:05"}
var d0, _ = time.Parse("15:04:05", "00:00:00")

// parseISOTime converts a [[HH:]MM:]SS[.mmm] string to a time.Duration.
func parseISOTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	numParts := len(parts)

	if numParts == 0 || numParts > 3 {
		return 0, fmt.Errorf("invalid time format, expected [[HH:]MM:]SS[.mmm]")
	}

	dur, err := time.Parse(durationFormat[numParts], s)
	if err != nil {
		return 0, err
	}

	return dur.Sub(d0), nil
}

// commandExists checks if a command is available in the system's PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
