# Backend TODO

## Critical ✅ COMPLETED
- [x] Fix CORS configuration — remove wildcard origin with credentials
- [x] Move JWT secret and DB credentials to environment variables only
- [x] Add rate limiting on auth endpoints (5 req/min per IP)
- [x] Fix auth service type assertions — add methods to repository interface
- [x] Fix balance calculation denominator bug in account_repository.go
- [x] Fix go.mod version (changed from 1.25 to 1.23)

## Security ✅ COMPLETED
- [x] Add security headers middleware (X-Frame-Options, X-Content-Type-Options, X-XSS-Protection)
- [x] Add request size limits (10MB)
- [x] Add refresh token cleanup mechanism (runs every 6 hours)

## Performance ✅ COMPLETED
- [x] Fix N+1 queries in analytics service (replaced with single aggregated SQL query)
- [x] Add structured logging with slog

## Remaining Tasks
- [ ] Add test coverage
- [ ] Add DB transactions for multi-step operations (user creation + token)
- [ ] Pre-allocate slices where size is known
- [ ] Extract magic strings to constants
- [ ] Strengthen password policy (add complexity requirements)
- [ ] Sanitize error messages — don't leak DB details
- [ ] Consistent error handling with custom error types
