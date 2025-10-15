package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BenchmarkResult represents a parsed benchmark result
type BenchmarkResult struct {
	Name           string
	Operations     int
	Duration       time.Duration
	Throughput     float64 // operations per second
	MemoryAllocs   int64
	MemoryBytes    int64
	AllocsPerOp    float64
	BytesPerOp     float64
}

// BenchmarkAnalyzer handles parsing and analysis of benchmark results
type BenchmarkAnalyzer struct {
	Results []BenchmarkResult
}

// NewBenchmarkAnalyzer creates a new analyzer
func NewBenchmarkAnalyzer() *BenchmarkAnalyzer {
	return &BenchmarkAnalyzer{
		Results: make([]BenchmarkResult, 0),
	}
}

// ParseBenchmarkOutput parses Go benchmark output
func (ba *BenchmarkAnalyzer) ParseBenchmarkOutput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	// Regex patterns for parsing benchmark output
	benchmarkPattern := regexp.MustCompile(`Benchmark(\w+)\s+(\d+)\s+(\d+(?:\.\d+)?)ns/op\s+(\d+)\s+(\d+(?:\.\d+)?)allocs/op\s+(\d+)\s+(\d+(?:\.\d+)?)B/op`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if matches := benchmarkPattern.FindStringSubmatch(line); matches != nil {
			result := BenchmarkResult{
				Name: matches[1],
			}
			
			// Parse operations
			if ops, err := strconv.Atoi(matches[2]); err == nil {
				result.Operations = ops
			}
			
			// Parse duration (nanoseconds)
			if ns, err := strconv.ParseFloat(matches[3], 64); err == nil {
				result.Duration = time.Duration(ns) * time.Nanosecond
				result.Throughput = float64(result.Operations) / result.Duration.Seconds()
			}
			
			// Parse memory allocations
			if allocs, err := strconv.ParseInt(matches[4], 10, 64); err == nil {
				result.MemoryAllocs = allocs
			}
			
			// Parse allocs per operation
			if allocsPerOp, err := strconv.ParseFloat(matches[5], 64); err == nil {
				result.AllocsPerOp = allocsPerOp
			}
			
			// Parse bytes per operation
			if bytesPerOp, err := strconv.ParseFloat(matches[7], 64); err == nil {
				result.BytesPerOp = bytesPerOp
			}
			
			// Parse total memory bytes
			if bytes, err := strconv.ParseInt(matches[6], 10, 64); err == nil {
				result.MemoryBytes = bytes
			}
			
			ba.Results = append(ba.Results, result)
		}
	}
	
	return scanner.Err()
}

// AnalyzeResults performs analysis on the benchmark results
func (ba *BenchmarkAnalyzer) AnalyzeResults() {
	fmt.Println("=== YAKVS Benchmark Analysis ===")
	fmt.Println()
	
	// Group results by command type
	commandGroups := make(map[string][]BenchmarkResult)
	for _, result := range ba.Results {
		commandType := ba.extractCommandType(result.Name)
		commandGroups[commandType] = append(commandGroups[commandType], result)
	}
	
	// Analyze each command type
	for commandType, results := range commandGroups {
		fmt.Printf("=== %s Command Analysis ===\n", strings.ToUpper(commandType))
		ba.analyzeCommandGroup(commandType, results)
		fmt.Println()
	}
	
	// Overall performance summary
	ba.printPerformanceSummary()
}

// extractCommandType extracts the command type from benchmark name
func (ba *BenchmarkAnalyzer) extractCommandType(name string) string {
	if strings.Contains(name, "Set") {
		return "SET"
	} else if strings.Contains(name, "Get") {
		return "GET"
	} else if strings.Contains(name, "Del") {
		return "DEL"
	} else if strings.Contains(name, "Mixed") {
		return "MIXED"
	} else if strings.Contains(name, "Concurrent") {
		return "CONCURRENT"
	} else if strings.Contains(name, "Memory") {
		return "MEMORY"
	}
	return "OTHER"
}

// analyzeCommandGroup analyzes a group of related benchmarks
func (ba *BenchmarkAnalyzer) analyzeCommandGroup(commandType string, results []BenchmarkResult) {
	if len(results) == 0 {
		return
	}
	
	// Sort by throughput (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Throughput > results[j].Throughput
	})
	
	fmt.Printf("Found %d benchmark(s) for %s command:\n", len(results), commandType)
	fmt.Println()
	
	for _, result := range results {
		fmt.Printf("  %-30s: %8.2f ops/sec (%.2f ms/op)\n", 
			result.Name, 
			result.Throughput, 
			float64(result.Duration.Nanoseconds())/float64(result.Operations)/1e6)
	}
	
	// Calculate statistics
	ba.printCommandStatistics(commandType, results)
}

// printCommandStatistics prints detailed statistics for a command group
func (ba *BenchmarkAnalyzer) printCommandStatistics(commandType string, results []BenchmarkResult) {
	if len(results) == 0 {
		return
	}
	
	// Calculate averages
	var totalThroughput, totalLatency, totalAllocs, totalBytes float64
	
	for _, result := range results {
		totalThroughput += result.Throughput
		totalLatency += float64(result.Duration.Nanoseconds()) / float64(result.Operations) / 1e6
		totalAllocs += result.AllocsPerOp
		totalBytes += result.BytesPerOp
	}
	
	avgThroughput := totalThroughput / float64(len(results))
	avgLatency := totalLatency / float64(len(results))
	avgAllocs := totalAllocs / float64(len(results))
	avgBytes := totalBytes / float64(len(results))
	
	fmt.Printf("\n  %s Command Statistics:\n", commandType)
	fmt.Printf("    Average Throughput: %.2f ops/sec\n", avgThroughput)
	fmt.Printf("    Average Latency:    %.2f ms/op\n", avgLatency)
	fmt.Printf("    Average Allocs:     %.2f allocs/op\n", avgAllocs)
	fmt.Printf("    Average Memory:     %.2f bytes/op\n", avgBytes)
	
	// Find best and worst performers
	best := results[0] // Already sorted by throughput
	worst := results[len(results)-1]
	
	fmt.Printf("\n  Performance Range:\n")
	fmt.Printf("    Best:  %s (%.2f ops/sec)\n", best.Name, best.Throughput)
	fmt.Printf("    Worst: %s (%.2f ops/sec)\n", worst.Name, worst.Throughput)
}

// printPerformanceSummary prints an overall performance summary
func (ba *BenchmarkAnalyzer) printPerformanceSummary() {
	fmt.Println("=== Overall Performance Summary ===")
	
	// Find the fastest operation overall
	var fastest BenchmarkResult
	for _, result := range ba.Results {
		if result.Throughput > fastest.Throughput {
			fastest = result
		}
	}
	
	if fastest.Name != "" {
		fmt.Printf("Fastest Operation: %s (%.2f ops/sec)\n", fastest.Name, fastest.Throughput)
	}
	
	// Calculate total operations across all benchmarks
	var totalOps int
	for _, result := range ba.Results {
		totalOps += result.Operations
	}
	
	fmt.Printf("Total Operations Tested: %d\n", totalOps)
	fmt.Printf("Total Benchmarks: %d\n", len(ba.Results))
	
	// Memory usage analysis
	ba.printMemoryAnalysis()
}

// printMemoryAnalysis prints memory usage analysis
func (ba *BenchmarkAnalyzer) printMemoryAnalysis() {
	fmt.Println("\n=== Memory Usage Analysis ===")
	
	var totalAllocs, totalBytes float64
	var memoryIntensive []BenchmarkResult
	
	for _, result := range ba.Results {
		totalAllocs += result.AllocsPerOp
		totalBytes += result.BytesPerOp
		
		// Identify memory-intensive operations
		if result.BytesPerOp > 1000 { // More than 1KB per operation
			memoryIntensive = append(memoryIntensive, result)
		}
	}
	
	avgAllocs := totalAllocs / float64(len(ba.Results))
	avgBytes := totalBytes / float64(len(ba.Results))
	
	fmt.Printf("Average Memory Allocations: %.2f allocs/op\n", avgAllocs)
	fmt.Printf("Average Memory Usage: %.2f bytes/op\n", avgBytes)
	
	if len(memoryIntensive) > 0 {
		fmt.Println("\nMemory-Intensive Operations (>1KB/op):")
		for _, result := range memoryIntensive {
			fmt.Printf("  %s: %.2f bytes/op\n", result.Name, result.BytesPerOp)
		}
	}
}

// GenerateReport generates a comprehensive benchmark report
func (ba *BenchmarkAnalyzer) GenerateReport() {
	fmt.Println("=== YAKVS Benchmark Report ===")
	fmt.Printf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()
	
	// System information
	fmt.Println("System Information:")
	fmt.Printf("  Go Version: %s\n", getGoVersion())
	fmt.Printf("  OS: %s\n", getOSInfo())
	fmt.Println()
	
	// Benchmark results
	ba.AnalyzeResults()
	
	// Recommendations
	ba.printRecommendations()
}

// printRecommendations prints performance recommendations
func (ba *BenchmarkAnalyzer) printRecommendations() {
	fmt.Println("=== Performance Recommendations ===")
	
	// Analyze results and provide recommendations
	var recommendations []string
	
	// Check for high memory usage
	for _, result := range ba.Results {
		if result.BytesPerOp > 10000 { // More than 10KB per operation
			recommendations = append(recommendations, 
				fmt.Sprintf("Consider optimizing %s - high memory usage (%.2f bytes/op)", 
					result.Name, result.BytesPerOp))
		}
	}
	
	// Check for low throughput
	for _, result := range ba.Results {
		if result.Throughput < 1000 { // Less than 1000 ops/sec
			recommendations = append(recommendations, 
				fmt.Sprintf("Consider optimizing %s - low throughput (%.2f ops/sec)", 
					result.Name, result.Throughput))
		}
	}
	
	if len(recommendations) > 0 {
		fmt.Println("Recommendations:")
		for i, rec := range recommendations {
			fmt.Printf("  %d. %s\n", i+1, rec)
		}
	} else {
		fmt.Println("All benchmarks show good performance!")
	}
}

// Helper functions
func getGoVersion() string {
	// This would typically run `go version` command
	return "Go 1.21+ (estimated)"
}

func getOSInfo() string {
	// This would typically get OS information
	return "macOS/Linux (estimated)"
}

// Main function for standalone usage
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run benchmark_analyzer.go <benchmark_output_file>")
		fmt.Println("Example: go run benchmark_analyzer.go benchmark_output.txt")
		os.Exit(1)
	}
	
	analyzer := NewBenchmarkAnalyzer()
	
	if err := analyzer.ParseBenchmarkOutput(os.Args[1]); err != nil {
		fmt.Printf("Error parsing benchmark output: %v\n", err)
		os.Exit(1)
	}
	
	analyzer.GenerateReport()
}
