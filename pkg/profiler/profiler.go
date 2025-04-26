package profiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CPUProfile profiles CPU usage of a Go binary.
func CPUProfile(target string, outputFile string, duration int) error {
	fmt.Printf("Profiling CPU usage of %s for %d seconds...\n", target, duration)

	// Ensure target binary exists
	_, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("target binary not found: %w", err)
	}

	// Create absolute path for output file
	absOutput, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Run the binary with CPU profiling enabled
	cmd := exec.Command(target, "-cpuprofile", absOutput)

	// Start the process
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start target binary: %w", err)
	}

	// Kill the process after the specified duration
	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		cmd.Process.Kill()
	}()

	// Wait for the process to complete
	err = cmd.Wait()
	if err != nil && err.Error() != "signal: killed" {
		return fmt.Errorf("error running target binary: %w", err)
	}

	fmt.Printf("CPU profile saved to %s\n", absOutput)
	fmt.Println("Use 'goforge profile visualize " + absOutput + "' to analyze the profile")

	return nil
}

// MemoryProfile profiles memory usage of a Go binary.
func MemoryProfile(target string, outputFile string) error {
	fmt.Printf("Profiling memory usage of %s...\n", target)

	// Ensure target binary exists
	_, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("target binary not found: %w", err)
	}

	// Create absolute path for output file
	absOutput, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	// Run the binary with memory profiling enabled
	cmd := exec.Command(target, "-memprofile", absOutput)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run memory profile: %w\nOutput: %s", err, output)
	}

	fmt.Printf("Memory profile saved to %s\n", absOutput)
	fmt.Println("Use 'goforge profile visualize " + absOutput + "' to analyze the profile")

	return nil
}

// Visualize displays a profile in a human-readable format.
func Visualize(profileFile string) error {
	fmt.Printf("Visualizing profile %s...\n", profileFile)

	// Ensure profile file exists
	_, err := os.Stat(profileFile)
	if err != nil {
		return fmt.Errorf("profile file not found: %w", err)
	}

	// Use 'go tool pprof' to generate a visualization
	// Here we'll use the text output, but in a real implementation we could
	// generate graphical visualizations (SVG, PDF, etc.)
	cmd := exec.Command("go", "tool", "pprof", "-text", profileFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to visualize profile: %w", err)
	}

	// Display the profile information
	fmt.Println("\nProfile Analysis:")
	fmt.Println(string(output))

	// In a real implementation, we could also offer to open a web browser with
	// the interactive pprof interface
	fmt.Println("\nTip: For more detailed analysis, run:")
	fmt.Printf("go tool pprof -http=:8080 %s\n", profileFile)

	return nil
}
