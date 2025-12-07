# ADR-002: Replace euicc-go/bertlv with Minimalistic BER-TLV Encoder/Decoder for Performance

## Status
Accepted

## Context

The TLV encoder/decoder for ISO8583 field 55 was originally implemented using the [github.com/euicc-go/bertlv](https://github.com/euicc-go/bertlv) library. However, benchmarks revealed significant performance and memory allocation issues with this approach, making it unsuitable for high-throughput ISO8583 systems. To address this, a custom minimalistic BER-TLV encoder/decoder was implemented, optimized for flat TLV structures and short-form lengths typical of ISO8583 field 55. Benchmarks were run to validate the improvements.

### TLV Encoder Benchmarks: Evolution and Rationale

#### Previous Implementation (Library-Based)
| Encoder     | ns/op | B/op | allocs/op |
|-------------|-------|------|-----------|
| TLV Encode  | 430   | 528  | 24        |
| TLV Decode  | 430   | 528  | 24        |

- The TLV encoder/decoder was based on the [github.com/euicc-go/bertlv](https://github.com/euicc-go/bertlv) BER-TLV library implementation.
- It was much slower and allocated far more memory than other encoders.
- Most allocations were unnecessary, coming from repeated struct/slice creation and deep parsing logic in the library.
- This was not suitable for high-throughput ISO8583 systems, especially for field 55 (DE55) EMV data.

#### Minimalistic Custom Implementation (Current)
| Encoder         | ns/op   | B/op | allocs/op |
|----------------|---------|------|-----------|
| TLV Encode     | 56.37   | 24   | 2         |
| TLV Decode     | 56.26   | 24   | 2         |

```
BenchmarkTLVEncode-10           21225777                56.37 ns/op           24 B/op          2 allocs/op
BenchmarkTLVDecode-10           21018322                56.26 ns/op           24 B/op          2 allocs/op
ok      github.com/hkumarmk/iso8583-lite/pkg/encoding   3.615s
```

#### Reason for Change
- The minimalistic TLV encoder/decoder was designed to:
	- Eliminate unnecessary allocations
	- Optimize for flat TLV structures and short-form lengths (as used in ISO8583 field 55)
	- Achieve performance and memory usage comparable to other encoders
	- Simplify code for maintainability and future extension
- Benchmarks show a dramatic improvement: ~8x faster and ~20x less memory usage than the previous approach.
- Real-world DE55 samples validate correctness and robustness.

## Decision
Adopt the minimalistic TLV encoder/decoder implementation for production use, based on:
- Fast encode/decode times
- Low memory allocations
- Simplicity and maintainability
- Real-world DE55 sample validation

## Consequences
- TLV logic is robust and suitable for high-throughput ISO8583 systems
- Future optimizations should maintain or improve these metrics
- Benchmark results will be referenced in documentation and future design reviews

## TODOs, Limitations, and Next Steps

### Current Limitations
- Only supports flat TLV structures (no constructed/nested tags)
- Only supports short-form lengths (single-byte length fields)
- Limited error handling for malformed or deeply nested TLVs
- May not be fully compliant with all BER-TLV/EMV/ISO8583 extensions

### TODO / Next Steps
- Extend parser to support constructed tags (nested TLVs) if required
- Add support for long-form lengths (multi-byte length fields)
- Improve error handling and validation for malformed TLVs
- Benchmark and validate against additional real-world data and full spec requirements
- Document compliance and limitations in user-facing docs

### Standards & References
- ISO 8583: https://en.wikipedia.org/wiki/ISO_8583
- EMV Book 3: https://www.emvco.com/emv-technologies/specifications/
- ISO8583 Field 55: https://www.eftlab.com/knowledge-base/211-iso-8583-field-55-icc-system-related-data/
- BER-TLV: https://www.eftlab.com/knowledge-base/128-ber-tlv/

## References
- See `/pkg/encoding/tlv.go` and `/pkg/encoding/tlv_test.go` for implementation and tests
