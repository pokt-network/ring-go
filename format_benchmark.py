#!/usr/bin/env python3
"""
Benchmark comparison tool for ring-go crypto backends.
Compares Decred (pure Go) vs Ethereum (libsecp256k1) implementations.
"""

import sys
import re
import subprocess
import os
from collections import defaultdict

def run_benchmarks():
    """Run benchmarks for both backends and return results."""
    results = defaultdict(dict)

    print("🔬 Benchmarking ring signature crypto backends...")
    print("=" * 70)
    print()

    # Run Decred backend benchmarks
    print("Running benchmarks with Decred backend (Pure Go)...")
    cmd = ["go", "test", ".", "-bench=BenchmarkSign.*_Secp256k1|BenchmarkVerify.*_Secp256k1",
           "-benchmem", "-run=^$", "-benchtime=2s"]
    try:
        output = subprocess.run(cmd, capture_output=True, text=True, timeout=45).stdout
        parse_results(output, results, "Decred")
    except subprocess.TimeoutExpired:
        print("⚠️  Decred benchmarks timed out")

    # Run Ethereum backend benchmarks
    print("Running benchmarks with Ethereum backend (CGO + libsecp256k1)...")
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"
    cmd = ["go", "test", "-tags=ethereum_secp256k1", ".",
           "-bench=BenchmarkSign.*_Secp256k1|BenchmarkVerify.*_Secp256k1",
           "-benchmem", "-run=^$", "-benchtime=2s"]
    try:
        output = subprocess.run(cmd, capture_output=True, text=True, timeout=45, env=env).stdout
        parse_results(output, results, "Ethereum")
    except subprocess.TimeoutExpired:
        print("⚠️  Ethereum benchmarks timed out")

    return results

def parse_results(output, results, backend):
    """Parse benchmark output and store in results dict."""
    for line in output.split('\n'):
        # Pattern: BenchmarkSign2_Secp256k1-8   1000   1234567 ns/op
        match = re.match(r'Benchmark(Sign|Verify)(\d+)_Secp256k1.*?\s+\d+\s+(\d+(?:\.\d+)?)\s+(ns|µs|ms)/op', line)
        if match:
            operation = match.group(1)
            ring_size = int(match.group(2))
            time_val = float(match.group(3))
            time_unit = match.group(4)

            # Normalize to microseconds
            if time_unit == "ns":
                time_us = time_val / 1000
            elif time_unit == "µs":
                time_us = time_val
            elif time_unit == "ms":
                time_us = time_val * 1000
            else:
                time_us = time_val

            results[(operation, ring_size)][backend] = time_us

def parse_stdin():
    """Parse benchmark output from stdin (for piped input)."""
    results = defaultdict(dict)
    current_backend = None

    for line in sys.stdin:
        # Detect backend type
        if "# Decred Backend Results:" in line:
            current_backend = "Decred"
            continue
        elif "# Ethereum Backend Results:" in line:
            current_backend = "Ethereum"
            continue

        # Parse benchmark lines
        match = re.match(r'Benchmark(Sign|Verify)(\d+)_(Decred|Ethereum|Secp256k1).*?\s+\d+\s+(\d+(?:\.\d+)?)\s+(ns|µs|ms)/op', line)
        if match:
            operation = match.group(1)
            ring_size = int(match.group(2))
            backend_name = match.group(3)
            time_val = float(match.group(4))
            time_unit = match.group(5)

            # If backend was renamed via sed, use current_backend
            if backend_name in ["Decred", "Ethereum"]:
                backend = backend_name
            elif current_backend:
                backend = current_backend
            else:
                continue

            # Normalize to microseconds
            if time_unit == "ns":
                time_us = time_val / 1000
            elif time_unit == "µs":
                time_us = time_val
            elif time_unit == "ms":
                time_us = time_val * 1000
            else:
                time_us = time_val

            results[(operation, ring_size)][backend] = time_us

    return results

def format_time(microseconds):
    """Format time in appropriate unit."""
    if microseconds < 1:
        return f"{microseconds*1000:.1f} ns"
    elif microseconds < 1000:
        return f"{microseconds:.1f} µs"
    else:
        return f"{microseconds/1000:.1f} ms"

def print_comparison(results):
    """Print formatted comparison table."""
    if not results:
        print("No benchmark results found.")
        return

    print()
    # Print Sign performance
    print("🔍 SIGN PERFORMANCE (Ring Signatures):")
    print("Ring  Decred          Ethereum        Improvement")
    print("Size  (Pure Go)       (libsecp256k1)  (% faster)")
    print("-" * 4 + "  " + "-" * 15 + " " + "-" * 15 + " " + "-" * 11)

    sizes = sorted(set(size for op, size in results.keys() if op == "Sign"))
    for size in sizes[:5]:  # Limit to first 5 sizes
        key = ("Sign", size)
        if key in results and "Decred" in results[key] and "Ethereum" in results[key]:
            decred = results[key]["Decred"]
            ethereum = results[key]["Ethereum"]
            improvement = ((decred - ethereum) / decred) * 100
            print(f"{size:<4}  {format_time(decred):<15} {format_time(ethereum):<15} {improvement:.0f}%")

    print()
    # Print Verify performance
    print("🔍 VERIFY PERFORMANCE (Ring Signatures):")
    print("Ring  Decred          Ethereum        Improvement")
    print("Size  (Pure Go)       (libsecp256k1)  (% faster)")
    print("-" * 4 + "  " + "-" * 15 + " " + "-" * 15 + " " + "-" * 11)

    sizes = sorted(set(size for op, size in results.keys() if op == "Verify"))
    for size in sizes[:5]:  # Limit to first 5 sizes
        key = ("Verify", size)
        if key in results and "Decred" in results[key] and "Ethereum" in results[key]:
            decred = results[key]["Decred"]
            ethereum = results[key]["Ethereum"]
            improvement = ((decred - ethereum) / decred) * 100
            print(f"{size:<4}  {format_time(decred):<15} {format_time(ethereum):<15} {improvement:.0f}%")

    print()
    print("=" * 70)

def main():
    """Main entry point."""
    if "--run-benchmarks" in sys.argv:
        # Run benchmarks directly
        results = run_benchmarks()
        print_comparison(results)
    else:
        # Parse from stdin (for piped input)
        results = parse_stdin()
        print_comparison(results)

if __name__ == "__main__":
    main()