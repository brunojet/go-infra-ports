Summary of changes: rest_v2 registry & tests

- Behavior change: `Merge` on the `rest_registry` now panics when called with a nil or a foreign registry type (requires concrete `*restRegistry`). This was an intentional breaking change per review guidance.
- Tests updated and added:
  - Updated existing tests to use V2 types (`RestRequestV2`, `RestResponseV2`).
  - Added mapping tests: `internal/repositories/rest_v2/rest_repository_local_map_test.go`.
  - Added execution/error-path tests: `internal/repositories/rest_v2/rest_repository_local_execute_test.go`.
- Linter fixes applied (revive comments and removed unused symbols).
- Coverage: `internal/repositories/rest_v2` now reports `93.7%` statements.

How to verify locally:

1. Run tests for the package:

   go test ./internal/repositories/rest_v2 -coverprofile=coverage/rest_v2.out -covermode=atomic

2. View coverage summary:

   go tool cover -func=coverage/rest_v2.out

3. Generate HTML:

   go tool cover -html=coverage/rest_v2.out -o coverage/rest_v2.html

Notes:
- I ran `golangci-lint run ./...` (no issues) and added targeted //nolint comments only in tests where the analyzer was a false-positive for test flows.
- If you want, I can open a PR with these changes and include the coverage HTML as an artifact.
