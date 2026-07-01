# QA Review: original_vs_port/compare

Date: 2026-07-01

Scope: `test/original_vs_port/compare`

This review focuses on the reliability of the comparison tool from a QA engineering perspective: false negatives, false positives, diagnostics, and test coverage for the comparer itself.

## Findings

### P1: Non-`none` pending Go actions can be silently dropped

File: `test/original_vs_port/compare/comparer.go`

Relevant code:

- `processEvent` stores a Go port action in `fc.pending`.
- `flushPendingBeforeNonSelf` clears `fc.pending`.
- It reports a mismatch for non-`none` pending actions only when the next original line is `comparable`.

Risk:

If the Go port incorrectly returns `dahai`, `reach`, `hora`, or another real action, and the next log line is not a comparable action (`dora`, `reach_accepted`, `end_game`, etc.), the pending action is discarded without a mismatch. This creates a false negative, which is critical for a tool whose purpose is behavioral difference detection.

Recommendation:

When a pending Go action is non-`none`, reaching any subsequent line without a matching original self action should be treated as a mismatch. If there are protocol timing exceptions, encode them explicitly and cover them with tests.

## P2: Implicit `none` is always counted as a match

File: `test/original_vs_port/compare/comparer.go`

Relevant code:

- `flushPendingBeforeNonSelf`
- `flushPendingAtEOF`

Risk:

Pending `none` actions are counted as matches when no original self action appears. This may be intended to compensate for original logs not recording explicit pass decisions, but it can also hide cases where the Go port returns `none` at a point where no action opportunity exists.

Recommendation:

Validate that the previous event was a legal explicit-pass opportunity before counting this as a match, or track implicit passes in a separate counter instead of merging them into `matches`. The summary should make it clear how many matches are direct action matches versus inferred pass matches.

## P2: `consumed` order may produce false mismatches

File: `test/original_vs_port/compare/action.go`

Relevant code:

- `normalizeRawAction`
- `actionsEqual`

Risk:

`consumed` is compared with `slices.Equal`, so the order must match exactly. For actions where consumed tile order does not affect action meaning, this can produce false positives. This is especially relevant for `pon`, `ankan`, `kakan`, and `daiminkan`. For `chi`, confirm the protocol semantics before changing behavior.

Recommendation:

Normalize consumed tiles for order-insensitive action types before comparison. Keep `chi` order-sensitive only if the protocol requires that order to carry meaning.

## P2: Core comparer behavior has insufficient direct tests

File: `test/original_vs_port/compare/comparer_test.go`

Risk:

The current tests cover `findPlayer` and representative action normalization, but not the stateful comparison logic. The highest-risk behavior is in pending action handling, which is currently not directly characterized.

Recommendation:

Add tests for:

- Go action followed by matching original self action.
- Go non-`none` action followed by a non-comparable original line.
- Go non-`none` action followed by EOF or `end_game`.
- Go `none` inferred as an implicit pass.
- Original self action with no pending Go action.
- Mismatch limit behavior.
- Summary counters for direct matches, implicit pass matches, mismatches, and errors.

## Verification

Command:

```sh
GOEXPERIMENT=jsonv2 go test ./test/original_vs_port/compare
```

Result:

The package tests passed when run outside the sandbox. The first sandboxed run failed because the Go build cache under `AppData\Local\go-build` was not accessible.
