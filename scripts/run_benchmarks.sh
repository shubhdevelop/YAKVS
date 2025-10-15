#!/bin/bash

# YAKVS Benchmark Runner Script
# This script runs comprehensive benchmarks for SET, GET, and DEL commands

echo "=========================================="
echo "YAKVS Command Throughput Benchmarks"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to run benchmark and capture results
run_benchmark() {
    local benchmark_name=$1
    local description=$2
    
    echo -e "${BLUE}Running: $description${NC}"
    echo "----------------------------------------"
    
    # Run benchmark 3 times and get the best result
    go test -bench="$benchmark_name" -benchmem -count=3 -timeout=10m | grep -E "(Benchmark|PASS|FAIL)" | tail -1
    
    echo ""
}

# Function to run all basic command benchmarks
run_basic_benchmarks() {
    echo -e "${YELLOW}=== Basic Command Benchmarks ===${NC}"
    echo ""
    
    run_benchmark "BenchmarkSetCommand" "SET Command Throughput"
    run_benchmark "BenchmarkGetCommand" "GET Command Throughput"
    run_benchmark "BenchmarkDelCommand" "DEL Command Throughput"
}

# Function to run cache miss benchmarks
run_cache_miss_benchmarks() {
    echo -e "${YELLOW}=== Cache Miss Benchmarks ===${NC}"
    echo ""
    
    run_benchmark "BenchmarkGetCommand_Misses" "GET Command with Cache Misses"
    run_benchmark "BenchmarkDelCommand_Misses" "DEL Command with Non-existent Keys"
}

# Function to run large value benchmarks
run_large_value_benchmarks() {
    echo -e "${YELLOW}=== Large Value Benchmarks ===${NC}"
    echo ""
    
    run_benchmark "BenchmarkSetCommand_LargeValues" "SET Command with Large Values (1KB)"
}

# Function to run mixed workload benchmarks
run_mixed_benchmarks() {
    echo -e "${YELLOW}=== Mixed Workload Benchmarks ===${NC}"
    echo ""
    
    run_benchmark "BenchmarkMixedWorkload_SetGet" "Mixed SET/GET Workload (70% SET, 30% GET)"
    run_benchmark "BenchmarkMixedWorkload_SetGetDel" "Mixed SET/GET/DEL Workload (50% SET, 30% GET, 20% DEL)"
}

# Function to run concurrent benchmarks
run_concurrent_benchmarks() {
    echo -e "${YELLOW}=== Concurrent Benchmarks ===${NC}"
    echo ""
    
    run_benchmark "BenchmarkConcurrentSet" "Concurrent SET Operations"
    run_benchmark "BenchmarkConcurrentGet" "Concurrent GET Operations"
    run_benchmark "BenchmarkConcurrentDel" "Concurrent DEL Operations"
}

# Function to run memory benchmarks
run_memory_benchmarks() {
    echo -e "${YELLOW}=== Memory Usage Benchmarks ===${NC}"
    echo ""
    
    run_benchmark "BenchmarkMemoryUsage_Set" "Memory Usage for SET Operations"
}

# Function to run all benchmarks
run_all_benchmarks() {
    echo -e "${GREEN}Running All Benchmarks...${NC}"
    echo ""
    
    run_basic_benchmarks
    run_cache_miss_benchmarks
    run_large_value_benchmarks
    run_mixed_benchmarks
    run_concurrent_benchmarks
    run_memory_benchmarks
}

# Function to run quick benchmarks (basic commands only)
run_quick_benchmarks() {
    echo -e "${GREEN}Running Quick Benchmarks (Basic Commands Only)...${NC}"
    echo ""
    
    run_basic_benchmarks
}

# Function to generate detailed report
generate_report() {
    local output_file="benchmark_report_$(date +%Y%m%d_%H%M%S).txt"
    
    echo -e "${BLUE}Generating detailed benchmark report...${NC}"
    echo "Report will be saved to: $output_file"
    
    {
        echo "YAKVS Benchmark Report"
        echo "Generated: $(date)"
        echo "=========================================="
        echo ""
        
        echo "System Information:"
        echo "OS: $(uname -s)"
        echo "Architecture: $(uname -m)"
        echo "CPU Cores: $(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo "Unknown")"
        echo ""
        
        echo "Go Version:"
        go version
        echo ""
        
        echo "Benchmark Results:"
        echo "=================="
        
    } > "$output_file"
    
    # Run all benchmarks and append to report
    go test -bench=. -benchmem -count=3 -timeout=10m >> "$output_file" 2>&1
    
    echo -e "${GREEN}Report generated: $output_file${NC}"
}

# Function to show help
show_help() {
    echo "YAKVS Benchmark Runner"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  all       Run all benchmarks (default)"
    echo "  basic     Run basic command benchmarks only"
    echo "  quick     Run quick benchmarks (basic commands)"
    echo "  mixed     Run mixed workload benchmarks"
    echo "  concurrent Run concurrent benchmarks"
    echo "  memory    Run memory usage benchmarks"
    echo "  report    Generate detailed benchmark report"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all benchmarks"
    echo "  $0 basic              # Run basic benchmarks only"
    echo "  $0 report             # Generate detailed report"
}

# Main script logic
case "${1:-all}" in
    "all")
        run_all_benchmarks
        ;;
    "basic")
        run_basic_benchmarks
        ;;
    "quick")
        run_quick_benchmarks
        ;;
    "mixed")
        run_mixed_benchmarks
        ;;
    "concurrent")
        run_concurrent_benchmarks
        ;;
    "memory")
        run_memory_benchmarks
        ;;
    "report")
        generate_report
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        echo -e "${RED}Unknown option: $1${NC}"
        show_help
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}Benchmark run completed!${NC}"
echo ""
echo "To run specific benchmarks manually:"
echo "  go test -bench=BenchmarkSetCommand -benchmem -count=3"
echo "  go test -bench=BenchmarkGetCommand -benchmem -count=3"
echo "  go test -bench=BenchmarkDelCommand -benchmem -count=3"
echo ""
echo "For more detailed analysis, use:"
echo "  go test -bench=. -benchmem -count=3 -cpuprofile=cpu.prof -memprofile=mem.prof"
