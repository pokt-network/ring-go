#!/usr/bin/env python3

import sys
import re
from collections import defaultdict

def format_time(ns):
    """Format nanoseconds into human-readable time units"""
    ns = float(ns)
    if ns >= 1_000_000:
        return f"{ns/1_000_000:.1f} ms"
    elif ns >= 1_000:
        return f"{ns/1_000:.1f} μs"
    else:
        return f"{ns:.0f} ns"

def format_memory(bytes_val):
    """Format bytes into human-readable memory units"""
    bytes_val = float(bytes_val)
    if bytes_val >= 1_048_576:
        return f"{bytes_val/1_048_576:.1f} MB"
    elif bytes_val >= 1_024:
        return f"{bytes_val/1_024:.1f} KB"
    else:
        return f"{bytes_val:.0f} B"

def format_number(num):
    """Format large numbers with K/M suffixes"""
    num = float(num)
    if num >= 1_000_000:
        return f"{num/1_000_000:.1f}M"
    elif num >= 1_000:
        return f"{num/1_000:.1f}K"
    else:
        return f"{num:.0f}"

def parse_benchmark_output():
    """Parse Go benchmark output from stdin"""
    data = defaultdict(lambda: defaultdict(lambda: defaultdict(dict)))

    for line in sys.stdin:
        line = line.strip()
        if not line.startswith('Benchmark') or 'ns/op' not in line:
            continue

        # Parse benchmark line: BenchmarkSign2_Decred-10  795  1550968 ns/op  5013 B/op  84 allocs/op
        parts = line.split()
        if len(parts) < 8:
            continue

        bench_name = parts[0]
        iterations = parts[1]
        ns_per_op = parts[2]
        bytes_per_op = parts[4]
        allocs_per_op = parts[6]

        # Extract operation, ring size, and backend from benchmark name
        # BenchmarkSign2_Decred-10 -> Sign, 2, Decred
        # BenchmarkVerify32_Ethereum-10 -> Verify, 32, Ethereum
        match = re.match(r'Benchmark(\w+?)(\d+)_(\w+)-\d+', bench_name)
        if not match:
            continue

        operation, ring_size, backend = match.groups()
        ring_size = int(ring_size)

        # Map backend names for clarity
        if backend == 'Decred':
            backend = 'Decred'
        elif backend == 'Ethereum':
            backend = 'Ethereum'
        elif backend == 'Secp256k1':
            backend = 'Secp256k1'
        elif backend == 'Ed25519':
            backend = 'Ed25519'

        data[operation][ring_size][backend] = {
            'ns': float(ns_per_op),
            'bytes': float(bytes_per_op),
            'allocs': float(allocs_per_op),
            'iterations': int(iterations)
        }

    return data

def print_formatted_results(data):
    """Print formatted benchmark results with side-by-side comparison"""
    operations = ['Sign', 'Verify']
    ring_sizes = [2, 4, 8, 16, 32]

    # ANSI color codes
    BOLD = '\033[1m'
    CYAN = '\033[36m'
    GREEN = '\033[32m'
    YELLOW = '\033[33m'
    RED = '\033[31m'
    BLUE = '\033[34m'
    RESET = '\033[0m'

    for operation in operations:
        if operation not in data:
            continue

        print(f"\n{BOLD}{CYAN}🔍 {operation.upper()} PERFORMANCE (Ring Signatures):{RESET}")
        print(f"{BOLD}{'Ring':<5} {'Decred':<15} {'Ethereum':<15} {'Improvement':<12}{RESET}")
        print(f"{BOLD}{'Size':<5} {'(Pure Go)':<15} {'(libsecp256k1)':<15} {'(% faster)':<12}{RESET}")
        print(f"{'----':<5} {'-----------':<15} {'---------------':<15} {'-----------':<12}")

        for ring_size in ring_sizes:
            if ring_size not in data[operation]:
                continue

            backends_data = data[operation][ring_size]

            if 'Decred' in backends_data and 'Ethereum' in backends_data:
                decred_time = backends_data['Decred']['ns']
                ethereum_time = backends_data['Ethereum']['ns']

                decred_str = format_time(decred_time)
                ethereum_str = format_time(ethereum_time)

                improvement = ((decred_time - ethereum_time) / decred_time) * 100

                # Color code the improvement based on performance
                if improvement >= 50:
                    improvement_color = GREEN
                elif improvement >= 30:
                    improvement_color = YELLOW
                else:
                    improvement_color = RED

                improvement_str = f"{improvement:.0f}%" if improvement > 0 else f"{abs(improvement):.0f}% slower"

                print(f"{BLUE}{ring_size:<5}{RESET} {decred_str:<15} {GREEN}{ethereum_str:<15}{RESET} {improvement_color}{improvement_str:<12}{RESET}")
            elif 'Decred' in backends_data:
                decred_time = backends_data['Decred']['ns']
                decred_str = format_time(decred_time)
                print(f"{BLUE}{ring_size:<5}{RESET} {decred_str:<15} {'N/A':<15} {'N/A':<12}")
            elif 'Ethereum' in backends_data:
                ethereum_time = backends_data['Ethereum']['ns']
                ethereum_str = format_time(ethereum_time)
                print(f"{BLUE}{ring_size:<5}{RESET} {'N/A':<15} {GREEN}{ethereum_str:<15}{RESET} {'N/A':<12}")

        print()  # Add separator between operations


if __name__ == "__main__":
    data = parse_benchmark_output()
    if not data:
        print("❌ No benchmark data found. Check that:")
        print("   • Go benchmarks are running correctly")
        print("   • Benchmark names match the expected pattern")
        print("   • CGO dependencies are available if needed")
        sys.exit(1)

    print_formatted_results(data)