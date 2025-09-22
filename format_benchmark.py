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
        if not line.startswith('Benchmark'):
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
        # BenchmarkVerify32_Fast-10 -> Verify, 32, Fast
        match = re.match(r'Benchmark(\w+?)(\d+)_(\w+)-\d+', bench_name)
        if not match:
            continue

        operation, ring_size, backend = match.groups()
        ring_size = int(ring_size)

        data[operation][ring_size][backend] = {
            'ns': float(ns_per_op),
            'bytes': float(bytes_per_op),
            'allocs': float(allocs_per_op),
            'iterations': int(iterations)
        }

    return data

def print_formatted_results(data):
    """Print formatted benchmark results"""
    operations = ['Sign', 'Verify']
    ring_sizes = [2, 4, 8, 16, 32]
    medals = ['🥇', '🥈', '🥉']

    for operation in operations:
        if operation not in data:
            continue

        print(f"\n🔍 {operation.upper()} PERFORMANCE (Ring Signatures):")
        print(f"{'Ring Size':<10} {'Backend':<15} {'Time/op':<12} {'Memory/op':<12} {'Allocs/op':<12} {'Performance':<12}")
        print(f"{'--------':<10} {'-------':<15} {'--------':<12} {'---------':<12} {'---------':<12} {'-----------':<12}")

        for ring_size in ring_sizes:
            if ring_size not in data[operation]:
                continue

            backends_data = data[operation][ring_size]
            # Sort backends by performance (time)
            backends_sorted = sorted(backends_data.items(), key=lambda x: x[1]['ns'])

            for i, (backend, metrics) in enumerate(backends_sorted):
                time_str = format_time(metrics['ns'])
                memory_str = format_memory(metrics['bytes'])
                allocs_str = format_number(metrics['allocs'])
                medal = medals[i] if i < len(medals) else ''

                print(f"{ring_size:<10} {backend:<15} {time_str:<12} {memory_str:<12} {allocs_str:<12} {medal}")

            # Add separator between ring sizes
            if ring_size != ring_sizes[-1]:
                print()

def print_performance_summary(data):
    """Print performance improvement summary"""
    print(f"\n📊 PERFORMANCE IMPROVEMENTS:")
    print(f"{'Operation':<12} {'Ring Size':<10} {'Decred':<12} {'Fast/Ethereum':<15} {'Improvement':<12}")
    print(f"{'--------':<12} {'--------':<10} {'------':<12} {'-------------':<15} {'-----------':<12}")

    operations = ['Sign', 'Verify']
    ring_sizes = [2, 8, 32]

    for operation in operations:
        if operation not in data:
            continue

        for ring_size in ring_sizes:
            if ring_size not in data[operation]:
                continue

            backends = data[operation][ring_size]
            if 'Decred' in backends and ('Fast' in backends or 'Ethereum' in backends):
                decred_time = backends['Decred']['ns']
                fast_backend = 'Fast' if 'Fast' in backends else 'Ethereum'
                fast_time = backends[fast_backend]['ns']

                improvement = ((decred_time - fast_time) / decred_time) * 100

                decred_str = format_time(decred_time)
                fast_str = format_time(fast_time)
                improvement_str = f"{improvement:.0f}% faster" if improvement > 0 else f"{abs(improvement):.0f}% slower"

                print(f"{operation:<12} {ring_size:<10} {decred_str:<12} {fast_str:<15} {improvement_str}")

if __name__ == "__main__":
    data = parse_benchmark_output()
    if not data:
        print("❌ No benchmark data found. Check that:")
        print("   • Go benchmarks are running correctly")
        print("   • Benchmark names match the expected pattern")
        print("   • CGO dependencies are available if needed")
        sys.exit(1)

    print_formatted_results(data)
    print_performance_summary(data)

    print(f"\n💡 KEY INSIGHTS:")
    print(f"   🥇 = Fastest    🥈 = Second fastest    🥉 = Third fastest")
    print(f"")
    print(f"   • Ethereum/Fast backend provides significant performance improvements")
    print(f"   • Larger ring sizes benefit more from optimized backends")
    print(f"   • Decred backend offers excellent portability (no CGO required)")
    print(f"   • Choose backend based on performance vs portability requirements")