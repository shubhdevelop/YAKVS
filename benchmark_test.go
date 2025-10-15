package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/shubhdevelop/YAKVS/command"
	"github.com/shubhdevelop/YAKVS/executor"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// BenchmarkTestResult results structure for analysis
type BenchmarkTestResult struct {
	Operation    string
	Duration     time.Duration
	Operations   int
	Throughput   float64 // operations per second
	MemoryAllocs uint64
	MemoryBytes  uint64
}

// Global store for benchmarks
var benchStore *store.Store

func init() {
	benchStore = store.NewStore()
}

// Helper function to create a command
func createCommand(name string, args []string) *parser.Command {
	return &parser.Command{
		Name: name,
		Args: args,
	}
}

// Helper function to generate random keys and values
func generateRandomData(count int) ([]string, []string) {
	keys := make([]string, count)
	values := make([]string, count)
	
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key_%d_%d", i, rand.Intn(1000000))
		values[i] = fmt.Sprintf("value_%d_%d", i, rand.Intn(1000000))
	}
	
	return keys, values
}

// Benchmark SET command throughput via async executor
func BenchmarkSetCommand(b *testing.B) {
	// Clean store before each benchmark
	benchStore = store.NewStore()
	
	keys, values := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], values[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
}

// Benchmark SET command with different value sizes via async executor
func BenchmarkSetCommand_LargeValues(b *testing.B) {
	benchStore = store.NewStore()
	
	// Create large values (1KB each)
	largeValue := make([]byte, 1024)
	for i := range largeValue {
		largeValue[i] = byte(rand.Intn(256))
	}
	
	keys, _ := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], string(largeValue)})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
}

// Benchmark GET command throughput via async executor
func BenchmarkGetCommand(b *testing.B) {
	// Pre-populate store with data
	keys, values := generateRandomData(b.N)
	
	// Set up data
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], values[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("GET", []string{keys[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
}

// Benchmark GET command with cache misses (non-existent keys) via async executor
func BenchmarkGetCommand_Misses(b *testing.B) {
	benchStore = store.NewStore()
	
	// Use keys that don't exist in the store
	keys, _ := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("GET", []string{keys[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
}

// Benchmark DEL command throughput via async executor
func BenchmarkDelCommand(b *testing.B) {
	// Pre-populate store with data
	keys, values := generateRandomData(b.N)
	
	// Set up data
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], values[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("DEL", []string{keys[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
}

// Benchmark DEL command with non-existent keys via async executor
func BenchmarkDelCommand_Misses(b *testing.B) {
	benchStore = store.NewStore()
	
	// Use keys that don't exist in the store
	keys, _ := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("DEL", []string{keys[i]})
		ch := make(chan executor.ResultWithError, 1)
		executor.ExecuteCommandAysnc(cmd, benchStore, ch)
		<-ch // Wait for completion
	}
}

// Mixed workload benchmarks via async executor
func BenchmarkMixedWorkload_SetGet(b *testing.B) {
	benchStore = store.NewStore()
	keys, values := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// 70% SET operations, 30% GET operations
		if i%10 < 7 {
			// SET operation
			cmd := createCommand("SET", []string{keys[i], values[i]})
			ch := make(chan executor.ResultWithError, 1)
			executor.ExecuteCommandAysnc(cmd, benchStore, ch)
			<-ch // Wait for completion
		} else {
			// GET operation
			cmd := createCommand("GET", []string{keys[i]})
			ch := make(chan executor.ResultWithError, 1)
			executor.ExecuteCommandAysnc(cmd, benchStore, ch)
			<-ch // Wait for completion
		}
	}
}

func BenchmarkMixedWorkload_SetGetDel(b *testing.B) {
	benchStore = store.NewStore()
	keys, values := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// 50% SET, 30% GET, 20% DEL operations
		op := i % 10
		if op < 5 {
			// SET operation
			cmd := createCommand("SET", []string{keys[i], values[i]})
			ch := make(chan executor.ResultWithError, 1)
			executor.ExecuteCommandAysnc(cmd, benchStore, ch)
			<-ch // Wait for completion
		} else if op < 8 {
			// GET operation
			cmd := createCommand("GET", []string{keys[i]})
			ch := make(chan executor.ResultWithError, 1)
			executor.ExecuteCommandAysnc(cmd, benchStore, ch)
			<-ch // Wait for completion
		} else {
			// DEL operation
			cmd := createCommand("DEL", []string{keys[i]})
			ch := make(chan executor.ResultWithError, 1)
			executor.ExecuteCommandAysnc(cmd, benchStore, ch)
			<-ch // Wait for completion
		}
	}
}

// Concurrent benchmarks
func BenchmarkConcurrentSet(b *testing.B) {
	benchStore = store.NewStore()
	keys, values := generateRandomData(b.N)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	var wg sync.WaitGroup
	numGoroutines := runtime.NumCPU()
	operationsPerGoroutine := b.N / numGoroutines
	
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				cmd := createCommand("SET", []string{keys[i], values[i]})
				setCmd := command.NewSetCommand(cmd, benchStore)
				setCmd.Execute()
			}
		}(g*operationsPerGoroutine, (g+1)*operationsPerGoroutine)
	}
	
	wg.Wait()
}

func BenchmarkConcurrentGet(b *testing.B) {
	// Pre-populate store
	keys, values := generateRandomData(b.N)
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], values[i]})
		setCmd := command.NewSetCommand(cmd, benchStore)
		setCmd.Execute()
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	var wg sync.WaitGroup
	numGoroutines := runtime.NumCPU()
	operationsPerGoroutine := b.N / numGoroutines
	
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				cmd := createCommand("GET", []string{keys[i]})
				getCmd := command.NewGetCommand(cmd, benchStore)
				getCmd.Execute()
			}
		}(g*operationsPerGoroutine, (g+1)*operationsPerGoroutine)
	}
	
	wg.Wait()
}

func BenchmarkConcurrentDel(b *testing.B) {
	// Pre-populate store
	keys, values := generateRandomData(b.N)
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], values[i]})
		setCmd := command.NewSetCommand(cmd, benchStore)
		setCmd.Execute()
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	var wg sync.WaitGroup
	numGoroutines := runtime.NumCPU()
	operationsPerGoroutine := b.N / numGoroutines
	
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				cmd := createCommand("DEL", []string{keys[i]})
				delCmd := command.NewDelCommand(cmd, benchStore)
				delCmd.Execute()
			}
		}(g*operationsPerGoroutine, (g+1)*operationsPerGoroutine)
	}
	
	wg.Wait()
}

// Memory usage benchmarks
func BenchmarkMemoryUsage_Set(b *testing.B) {
	benchStore = store.NewStore()
	keys, values := generateRandomData(b.N)
	
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		cmd := createCommand("SET", []string{keys[i], values[i]})
		setCmd := command.NewSetCommand(cmd, benchStore)
		setCmd.Execute()
	}
	
	b.StopTimer()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(m2.Alloc-m1.Alloc)/float64(b.N), "bytes/op")
}

// Throughput analysis function
func analyzeThroughput(b *testing.B, operation string) {
	duration := b.Elapsed()
	throughput := float64(b.N) / duration.Seconds()
	
	fmt.Printf("\n=== %s Throughput Analysis ===\n", operation)
	fmt.Printf("Operations: %d\n", b.N)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Throughput: %.2f ops/sec\n", throughput)
	fmt.Printf("Average latency: %.6f ms\n", float64(duration.Nanoseconds())/float64(b.N)/1e6)
}

// Custom benchmark runner with detailed analysis
func runDetailedBenchmark(b *testing.B, operation string, benchmarkFunc func(*testing.B)) {
	start := time.Now()
	benchmarkFunc(b)
	duration := time.Since(start)
	
	throughput := float64(b.N) / duration.Seconds()
	avgLatency := float64(duration.Nanoseconds()) / float64(b.N) / 1e6
	
	fmt.Printf("\n=== %s Detailed Results ===\n", operation)
	fmt.Printf("Operations: %d\n", b.N)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Throughput: %.2f ops/sec\n", throughput)
	fmt.Printf("Average latency: %.6f ms\n", avgLatency)
	// Note: Memory allocation stats are automatically reported by Go's benchmark framework
	// when using -benchmem flag
}

// Example usage function for running benchmarks programmatically
func ShowBenchmarkUsage() {
	// This function shows how to run benchmarks programmatically
	// and analyze results
	
	fmt.Println("YAKVS Command Throughput Benchmarks")
	fmt.Println("====================================")
	
	// You can run specific benchmarks like this:
	// go test -bench=BenchmarkSetCommand -benchmem -count=3
	// go test -bench=BenchmarkGetCommand -benchmem -count=3
	// go test -bench=BenchmarkDelCommand -benchmem -count=3
	
	fmt.Println("\nTo run benchmarks, use:")
	fmt.Println("go test -bench=. -benchmem -count=3")
	fmt.Println("\nFor specific benchmarks:")
	fmt.Println("go test -bench=BenchmarkSetCommand -benchmem")
	fmt.Println("go test -bench=BenchmarkGetCommand -benchmem")
	fmt.Println("go test -bench=BenchmarkDelCommand -benchmem")
	fmt.Println("\nFor concurrent benchmarks:")
	fmt.Println("go test -bench=BenchmarkConcurrent -benchmem")
	fmt.Println("\nFor mixed workload benchmarks:")
	fmt.Println("go test -bench=BenchmarkMixed -benchmem")
}
