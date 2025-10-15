# YAKVS Benchmarking Suite

This document describes the comprehensive benchmarking suite for YAKVS (Yet Another Key-Value Store) that measures the throughput and performance of SET, GET, and DEL commands.

## Overview

The benchmarking suite includes:
- **Basic Command Benchmarks**: Individual SET, GET, and DEL command performance
- **Cache Miss Benchmarks**: Performance with non-existent keys
- **Large Value Benchmarks**: Performance with large data (1KB+ values)
- **Mixed Workload Benchmarks**: Realistic workloads combining multiple operations
- **Concurrent Benchmarks**: Multi-threaded performance testing
- **Memory Usage Benchmarks**: Memory allocation and usage analysis

## Quick Start

### Running All Benchmarks

```bash
# Run all benchmarks with detailed output
./scripts/run_benchmarks.sh

# Run specific benchmark categories
./scripts/run_benchmarks.sh basic      # Basic commands only
./scripts/run_benchmarks.sh mixed      # Mixed workloads
./scripts/run_benchmarks.sh concurrent # Concurrent operations
./scripts/run_benchmarks.sh report     # Generate detailed report
```

### Manual Benchmark Execution

```bash
# Run all benchmarks
go test -bench=. -benchmem -count=3

# Run specific benchmarks
go test -bench=BenchmarkSetCommand -benchmem -count=3
go test -bench=BenchmarkGetCommand -benchmem -count=3
go test -bench=BenchmarkDelCommand -benchmem -count=3

# Run with profiling
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
```

## Benchmark Categories

### 1. Basic Command Benchmarks

These benchmarks test the core functionality of each command:

- **`BenchmarkSetCommand`**: Measures SET command throughput
- **`BenchmarkGetCommand`**: Measures GET command throughput  
- **`BenchmarkDelCommand`**: Measures DEL command throughput

### 2. Cache Miss Benchmarks

These benchmarks test performance when keys don't exist:

- **`BenchmarkGetCommand_Misses`**: GET operations on non-existent keys
- **`BenchmarkDelCommand_Misses`**: DEL operations on non-existent keys

### 3. Large Value Benchmarks

These benchmarks test performance with large data:

- **`BenchmarkSetCommand_LargeValues`**: SET operations with 1KB values

### 4. Mixed Workload Benchmarks

These benchmarks simulate realistic usage patterns:

- **`BenchmarkMixedWorkload_SetGet`**: 70% SET, 30% GET operations
- **`BenchmarkMixedWorkload_SetGetDel`**: 50% SET, 30% GET, 20% DEL operations

### 5. Concurrent Benchmarks

These benchmarks test multi-threaded performance:

- **`BenchmarkConcurrentSet`**: Concurrent SET operations
- **`BenchmarkConcurrentGet`**: Concurrent GET operations
- **`BenchmarkConcurrentDel`**: Concurrent DEL operations

### 6. Memory Usage Benchmarks

These benchmarks analyze memory consumption:

- **`BenchmarkMemoryUsage_Set`**: Memory usage for SET operations

## Understanding Benchmark Results

### Key Metrics

1. **Throughput (ops/sec)**: Operations per second - higher is better
2. **Latency (ns/op)**: Nanoseconds per operation - lower is better
3. **Memory Allocations (allocs/op)**: Memory allocations per operation - lower is better
4. **Memory Usage (B/op)**: Bytes allocated per operation - lower is better

### Sample Output

```
BenchmarkSetCommand-8             1000000    1200 ns/op    0 allocs/op    0 B/op
BenchmarkGetCommand-8             2000000     600 ns/op    0 allocs/op    0 B/op
BenchmarkDelCommand-8             1500000     800 ns/op    0 allocs/op    0 B/op
```

This means:
- SET: 1M operations, 1200ns per operation, no memory allocations
- GET: 2M operations, 600ns per operation, no memory allocations  
- DEL: 1.5M operations, 800ns per operation, no memory allocations

## Performance Analysis

### Expected Performance Characteristics

1. **GET operations** should be fastest (read-only, no modifications)
2. **SET operations** should be moderate (write operations)
3. **DEL operations** should be similar to SET (write operations)
4. **Concurrent operations** should scale with CPU cores
5. **Large values** will have higher memory usage

### Performance Optimization Tips

1. **High Memory Usage**: Look for operations with high `B/op` values
2. **Low Throughput**: Look for operations with low throughput
3. **High Allocation Count**: Look for operations with high `allocs/op` values

## Advanced Usage

### Profiling

Generate detailed performance profiles:

```bash
# CPU profiling
go test -bench=. -cpuprofile=cpu.prof

# Memory profiling  
go test -bench=. -memprofile=mem.prof

# Block profiling
go test -bench=. -blockprofile=block.prof
```

### Custom Benchmark Parameters

Modify benchmark parameters in `benchmark_test.go`:

```go
// Change number of operations
b.N = 1000000

// Change value sizes
largeValue := make([]byte, 4096) // 4KB values

// Change concurrency level
numGoroutines := 16
```

### Benchmark Analysis

Use the included analyzer:

```bash
# Generate benchmark output
go test -bench=. -benchmem > benchmark_output.txt

# Analyze results
go run benchmark_analyzer.go benchmark_output.txt
```

## Benchmark Configuration

### Environment Variables

```bash
# Set number of CPU cores for concurrent benchmarks
export GOMAXPROCS=8

# Enable memory profiling
export GOGC=100

# Set benchmark timeout
export BENCHMARK_TIMEOUT=10m
```

### Benchmark Parameters

Key parameters that affect benchmark results:

1. **`b.N`**: Number of operations (automatically set by Go)
2. **Value Size**: Size of values being stored
3. **Key Distribution**: Pattern of keys used
4. **Concurrency Level**: Number of goroutines for concurrent benchmarks

## Troubleshooting

### Common Issues

1. **Benchmark Timeout**: Increase timeout with `-timeout=30m`
2. **Memory Issues**: Reduce `b.N` or value sizes
3. **Inconsistent Results**: Run with `-count=5` for more stable results

### Performance Issues

1. **Low Throughput**: Check for locks, inefficient algorithms
2. **High Memory Usage**: Check for memory leaks, excessive allocations
3. **Inconsistent Performance**: Check for race conditions, improper synchronization

## Benchmark Results Interpretation

### Good Performance Indicators

- **GET operations**: > 1M ops/sec
- **SET operations**: > 500K ops/sec  
- **DEL operations**: > 500K ops/sec
- **Memory usage**: < 100 bytes/op
- **Allocations**: < 1 alloc/op

### Performance Bottlenecks

- **Lock contention**: Low concurrent performance
- **Memory pressure**: High memory usage
- **GC pressure**: Frequent garbage collection
- **Cache misses**: Poor key distribution

## Contributing

To add new benchmarks:

1. Create new benchmark functions following the naming convention
2. Add proper setup and teardown
3. Include memory allocation reporting
4. Add to the benchmark runner script
5. Update this documentation

## License

This benchmarking suite is part of the YAKVS project and follows the same license terms.
