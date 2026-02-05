# Backend Security & Architecture Improvements

## Summary
Completed major security hardening and performance optimizations for the GnuCash backend based on comprehensive code review findings.

## Critical Security Fixes ✅

### 1. Fixed CORS Misconfiguration
**Issue**: Wildcard origin (`*`) with credentials enabled — major security vulnerability
**Fix**:
- Removed wildcard CORS
- Added configurable allowed origins via `CORS_ALLOWED_ORIGINS` env var
- Default: `http://localhost:3000`
- Updated CORS middleware to validate origins against whitelist
- Changed OPTIONS response from 200 to 204 (correct status)

**Files changed**:
- `internal/interfaces/http/middleware/cors.go`
- `internal/config/config.go`
- `internal/interfaces/http/router.go`
- `cmd/server/main.go`

### 2. Removed Hardcoded Secrets
**Issue**: JWT secret and DB credentials stored in version control
**Fix**:
- Removed all secrets from `configs/config.yaml`
- Added validation requiring JWT_SECRET to be at least 32 characters
- Updated `.env.example` with proper documentation
- Fails fast if JWT_SECRET is not set or too weak

**Files changed**:
- `configs/config.yaml`
- `internal/config/config.go`
- `.env.example`

### 3. Added Rate Limiting
**Issue**: No protection against brute force attacks on auth endpoints
**Fix**:
- Implemented token bucket rate limiter
- Applied to all auth endpoints: 5 requests per minute per IP
- Includes automatic cleanup of old buckets
- Returns 429 Too Many Requests when limit exceeded

**Files added**:
- `internal/interfaces/http/middleware/ratelimit.go`

**Files changed**:
- `internal/interfaces/http/router.go`

### 4. Added Security Headers
**Issue**: Missing standard security headers
**Fix**:
- Added X-Content-Type-Options: nosniff
- Added X-Frame-Options: DENY
- Added X-XSS-Protection: 1; mode=block
- Added Referrer-Policy: strict-origin-when-cross-origin

**Files added**:
- `internal/interfaces/http/middleware/security.go`

### 5. Added Request Size Limits
**Issue**: No protection against memory exhaustion attacks
**Fix**:
- Added 10MB request body size limit
- Returns 413 Request Entity Too Large when exceeded

**Files added**:
- `internal/interfaces/http/middleware/request_size.go`

### 6. Added Refresh Token Cleanup
**Issue**: Expired tokens accumulated indefinitely in database
**Fix**:
- Created background service that runs every 6 hours
- Automatically deletes expired refresh tokens
- Logs cleanup operations with structured logging

**Files added**:
- `internal/infrastructure/persistence/postgres/token_cleanup.go`

**Files changed**:
- `cmd/server/main.go`

## Architecture Fixes ✅

### 7. Fixed Type Assertion Anti-Pattern
**Issue**: Auth service cast repository interface to concrete type, breaking Clean Architecture
**Fix**:
- Added refresh token methods to `UserRepository` interface
- Removed all type assertions from `auth_service.go`
- Removed unused import of `postgres` package from service layer

**Files changed**:
- `internal/domain/repository/user_repository.go`
- `internal/application/service/auth_service.go`

## Performance Optimizations ✅

### 8. Eliminated N+1 Queries in Analytics
**Issue**: Analytics service made separate database query for each account (N+1 problem)
**Fix**:
- Added `AggregateByAccountType` method to repository
- Replaced loop-based queries with single aggregated SQL query
- Uses PostgreSQL aggregation for better performance
- Normalized amounts to common denominator for correct calculations

**Files changed**:
- `internal/domain/repository/transaction_repository.go`
- `internal/infrastructure/persistence/postgres/transaction_repository.go`
- `internal/application/service/analytics_service.go`

### 9. Fixed Balance Calculation Bug
**Issue**: Used `MAX(quantity_denom)` which produces incorrect results with mixed-precision splits
**Fix**:
- Normalize all splits to common high-precision denominator (100,000)
- Use PostgreSQL `numeric` type to avoid integer division precision loss
- Applied to both `GetBalance()` and `GetBalanceWithChildren()`

**Files changed**:
- `internal/infrastructure/persistence/postgres/account_repository.go`

## Code Quality Improvements ✅

### 10. Added Structured Logging
**Issue**: Mixed use of `log.Printf` and `fmt.Printf` with no structure
**Fix**:
- Created logger package wrapping Go's `log/slog`
- Replaced all logging calls with structured logging
- Added request ID middleware for request tracing
- Added HTTP request logging middleware
- Supports JSON output for production (via `GO_ENV=production`)
- Supports text output for development

**Files added**:
- `pkg/logger/logger.go`
- `internal/interfaces/http/middleware/request_id.go`
- `internal/interfaces/http/middleware/logging.go`

**Files changed**:
- `cmd/server/main.go`
- `internal/interfaces/http/router.go`
- `go.mod` (added `github.com/google/uuid`)

### 11. Fixed go.mod Version
**Issue**: Referenced non-existent Go version 1.25
**Fix**: Changed to Go 1.23

**Files changed**:
- `go.mod`

## Environment Variables

Required environment variables:
```bash
# Required
JWT_SECRET=<min-32-chars>
DATABASE_USER=<username>
DATABASE_PASSWORD=<password>

# Optional (with defaults)
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=gnucash
DATABASE_SSLMODE=disable
SERVER_PORT=8080
GO_ENV=production  # for JSON logging
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

## Testing Recommendations

1. **Security Testing**:
   - Verify CORS only allows configured origins
   - Test rate limiting (should block after 5 requests/min)
   - Verify JWT secret validation fails with weak secrets
   - Test request size limits with large payloads

2. **Performance Testing**:
   - Compare analytics endpoint response times (should be much faster)
   - Verify balance calculations are correct with mixed denominators
   - Monitor database query counts (should be reduced)

3. **Monitoring**:
   - Check structured logs for request tracing
   - Monitor token cleanup logs every 6 hours
   - Track rate limit violations

## Remaining Work

See `backend/todo.md` for remaining items:
- Add comprehensive test coverage
- Add DB transactions for atomic operations
- Strengthen password policy
- Custom error types
- Additional code quality improvements
