# Benchmarks

The current library benchmarks for signing and verification are located below.
For ring signatures, the signing and verification time are linearly proportional
to the number of members of the ring (or "anonymity set"), which is what's observed.

> Note: the number directly after `BenchmarkSign` or `BenchmarkVerify` in the test
> name is the ring size being benchmarked.

> Note: the ns/op value on the right is the time it took for signing or verification
> (depending on the test). The middle value is the number of times the operation
> was executed by the Go benchmarker.

**Summary:**

- secp256k1 signing and verification is around 0.41ms per ring member
- ed25519 signing and verification is around is around 0.12ms per ring member

```bash
goos: linux
goarch: amd64
pkg: github.com/noot/ring-go
cpu: 12th Gen Intel(R) Core(TM) i7-1280P

BenchmarkSign2_Secp256k1-20                 1075           1113687 ns/op
BenchmarkSign4_Secp256k1-20                  651           1832647 ns/op
BenchmarkSign8_Secp256k1-20                  334           3389785 ns/op
BenchmarkSign16_Secp256k1-20                 184           6279636 ns/op
BenchmarkSign32_Secp256k1-20                  86          12556732 ns/op
BenchmarkSign64_Secp256k1-20                  44          24592647 ns/op
BenchmarkSign128_Secp256k1-20                 21          47949180 ns/op

BenchmarkSign2_Ed25519-20                   3184            338455 ns/op
BenchmarkSign4_Ed25519-20                   2102            561543 ns/op
BenchmarkSign8_Ed25519-20                   1141           1024334 ns/op
BenchmarkSign16_Ed25519-20                   601           1959393 ns/op
BenchmarkSign32_Ed25519-20                   312           3812862 ns/op
BenchmarkSign64_Ed25519-20                   158           7554431 ns/op
BenchmarkSign128_Ed25519-20                   72          15137610 ns/op

BenchmarkVerify2_Secp256k1-20               1647            759506 ns/op
BenchmarkVerify4_Secp256k1-20                788           1507848 ns/op
BenchmarkVerify8_Secp256k1-20                391           3060683 ns/op
BenchmarkVerify16_Secp256k1-20               193           6173042 ns/op
BenchmarkVerify32_Secp256k1-20                93          12352394 ns/op
BenchmarkVerify64_Secp256k1-20                45          25246452 ns/op
BenchmarkVerify128_Secp256k1-20               21          51882164 ns/op

BenchmarkVerify2_Ed25519-20                 4797            238406 ns/op
BenchmarkVerify4_Ed25519-20                 2349            457389 ns/op
BenchmarkVerify8_Ed25519-20                 1244            932592 ns/op
BenchmarkVerify16_Ed25519-20                 636           1823156 ns/op
BenchmarkVerify32_Ed25519-20                 320           3781398 ns/op
BenchmarkVerify64_Ed25519-20                 156           7524581 ns/op
BenchmarkVerify128_Ed25519-20                 78          14955353 ns/op
```

### v0.2.0

**Summary:**

- **Secp256k1**: ~**5–10% faster** across typical sizes. Examples:
    - Sign32: **11.70ms → 10.60ms** (~**9.4%** faster)
    - Sign64: **22.93ms → 21.02ms** (~**8.3%** faster)
    - Verify32: **11.32ms → 10.32ms** (~**8.9%** faster)
    - Verify64: **22.91ms → 20.69ms** (~**9.7%** faster)
- **Ed25519**: similar gains on average (**~5–8%** for common sizes).
    - Verify32: **4.80ms → 4.41ms** (~**8.0%**)
    - Sign32: **4.96ms → 4.62ms** (~**7.0%**)
- **Allocations**: consistently down **15–30%** depending on bench.  
  (e.g., Verify32 secp256k1: **843 → 682 allocs**; Sign2 secp256k1: **84 → 75 allocs**)


```
Before (v0.1.0)
-------------------------
goos: linux
goarch: amd64
pkg: github.com/pokt-network/ring-go
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkSign2_Secp256k1-32 1174 1,026,322 ns/op 5,021 B/op 84 allocs/op
BenchmarkSign4_Secp256k1-32 674 1,733,693 ns/op 8,536 B/op 140 allocs/op
BenchmarkSign8_Secp256k1-32 373 3,251,426 ns/op 15,567 B/op 252 allocs/op
BenchmarkSign16_Secp256k1-32 201 6,046,313 ns/op 29,633 B/op 476 allocs/op
BenchmarkSign32_Secp256k1-32 92 11,699,021 ns/op 57,812 B/op 925 allocs/op
BenchmarkSign64_Secp256k1-32 49 22,926,795 ns/op 114,472 B/op 1,825 allocs/op
BenchmarkSign128_Secp256k1-32 24 45,918,672 ns/op 227,992 B/op 3,634 allocs/op
BenchmarkVerify2_Secp256k1-32 1726 725,005 ns/op 3,404 B/op 53 allocs/op
BenchmarkVerify4_Secp256k1-32 849 1,384,662 ns/op 6,814 B/op 105 allocs/op
BenchmarkVerify8_Secp256k1-32 423 2,760,680 ns/op 13,646 B/op 209 allocs/op
BenchmarkVerify16_Secp256k1-32 210 5,600,718 ns/op 27,366 B/op 419 allocs/op
BenchmarkVerify32_Secp256k1-32 104 11,323,208 ns/op 55,044 B/op 843 allocs/op
BenchmarkVerify64_Secp256k1-32 49 22,912,945 ns/op 111,577 B/op 1,708 allocs/op
BenchmarkVerify128_Secp256k1-32 24 47,007,356 ns/op 228,492 B/op 3,502 allocs/op

BenchmarkSign2_Ed25519-32 2856 417,582 ns/op 4,672 B/op 70 allocs/op
BenchmarkSign4_Ed25519-32 1622 721,830 ns/op 8,032 B/op 119 allocs/op
BenchmarkSign8_Ed25519-32 888 1,305,862 ns/op 14,706 B/op 214 allocs/op
BenchmarkSign16_Ed25519-32 481 2,494,322 ns/op 28,041 B/op 403 allocs/op
BenchmarkSign32_Ed25519-32 241 4,961,154 ns/op 54,882 B/op 791 allocs/op
BenchmarkSign64_Ed25519-32 123 9,775,632 ns/op 109,027 B/op 1,579 allocs/op
BenchmarkSign128_Ed25519-32 57 19,487,283 ns/op 216,505 B/op 3,103 allocs/op
BenchmarkVerify2_Ed25519-32 3784 294,033 ns/op 3,217 B/op 44 allocs/op
BenchmarkVerify4_Ed25519-32 1953 584,609 ns/op 6,436 B/op 87 allocs/op
BenchmarkVerify8_Ed25519-32 975 1,209,195 ns/op 12,993 B/op 180 allocs/op
BenchmarkVerify16_Ed25519-32 500 2,421,635 ns/op 25,969 B/op 356 allocs/op
BenchmarkVerify32_Ed25519-32 246 4,801,361 ns/op 51,904 B/op 704 allocs/op
BenchmarkVerify64_Ed25519-32 121 9,620,050 ns/op 104,307 B/op 1,407 allocs/op
BenchmarkVerify128_Ed25519-32 60 19,432,527 ns/op 211,098 B/op 2,869 allocs/op

After (v0.2.0)
-------------------------
goos: linux
goarch: amd64
pkg: github.com/pokt-network/ring-go
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkSign2_Secp256k1-32 1224 995,215 ns/op 4,254 B/op 75 allocs/op
BenchmarkSign4_Secp256k1-32 724 1,629,372 ns/op 6,985 B/op 121 allocs/op
BenchmarkSign8_Secp256k1-32 409 2,901,838 ns/op 12,447 B/op 213 allocs/op
BenchmarkSign16_Secp256k1-32 218 5,481,979 ns/op 23,384 B/op 397 allocs/op
BenchmarkSign32_Secp256k1-32 97 10,596,642 ns/op 45,218 B/op 767 allocs/op
BenchmarkSign64_Secp256k1-32 49 21,016,084 ns/op 89,338 B/op 1,510 allocs/op
BenchmarkSign128_Secp256k1-32 25 42,112,324 ns/op 178,167 B/op 3,010 allocs/op
BenchmarkVerify2_Secp256k1-32 1854 646,071 ns/op 2,643 B/op 43 allocs/op
BenchmarkVerify4_Secp256k1-32 928 1,259,166 ns/op 5,249 B/op 85 allocs/op
BenchmarkVerify8_Secp256k1-32 454 2,522,735 ns/op 10,494 B/op 169 allocs/op
BenchmarkVerify16_Secp256k1-32 231 5,053,713 ns/op 21,021 B/op 339 allocs/op
BenchmarkVerify32_Secp256k1-32 112 10,319,757 ns/op 42,287 B/op 682 allocs/op
BenchmarkVerify64_Secp256k1-32 55 20,688,756 ns/op 85,595 B/op 1,381 allocs/op
BenchmarkVerify128_Secp256k1-32 26 42,974,759 ns/op 175,867 B/op 2,839 allocs/op

BenchmarkSign2_Ed25519-32 2902 407,265 ns/op 3,513 B/op 56 allocs/op
BenchmarkSign4_Ed25519-32 1677 702,718 ns/op 5,840 B/op 94 allocs/op
BenchmarkSign8_Ed25519-32 940 1,243,208 ns/op 10,326 B/op 159 allocs/op
BenchmarkSign16_Ed25519-32 502 2,346,752 ns/op 19,381 B/op 294 allocs/op
BenchmarkSign32_Ed25519-32 262 4,615,297 ns/op 37,654 B/op 574 allocs/op
BenchmarkSign64_Ed25519-32 133 9,010,359 ns/op 74,074 B/op 1,113 allocs/op
BenchmarkSign128_Ed25519-32 56 18,047,776 ns/op 147,841 B/op 2,223 allocs/op
BenchmarkVerify2_Ed25519-32 4135 275,322 ns/op 2,178 B/op 31 allocs/op
BenchmarkVerify4_Ed25519-32 2133 547,658 ns/op 4,334 B/op 61 allocs/op
BenchmarkVerify8_Ed25519-32 1076 1,096,213 ns/op 8,649 B/op 121 allocs/op
BenchmarkVerify16_Ed25519-32 534 2,203,613 ns/op 17,300 B/op 241 allocs/op
BenchmarkVerify32_Ed25519-32 262 4,414,711 ns/op 34,675 B/op 484 allocs/op
BenchmarkVerify64_Ed25519-32 133 8,848,523 ns/op 69,692 B/op 973 allocs/op
BenchmarkVerify128_Ed25519-32 64 18,091,954 ns/op 141,032 B/op 1,969 allocs/op
```
