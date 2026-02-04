# Testing Guide

## Prerequisites

Before testing, ensure you have:
1. A PostgreSQL database with GnuCash data
2. Backend running on http://localhost:8080
3. Frontend running on http://localhost:3000

## Backend Testing

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response: `{"status":"ok"}`

### 2. Test Authentication

**Register a user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Save the `access_token` and `refresh_token` from the response.

**Refresh token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### 3. Test Accounts API

**Get all accounts:**
```bash
curl http://localhost:8080/api/v1/accounts
```

**Get account hierarchy:**
```bash
curl http://localhost:8080/api/v1/accounts/hierarchy
```

**Get specific account:**
```bash
curl http://localhost:8080/api/v1/accounts/{ACCOUNT_GUID}
```

**Get account balance:**
```bash
curl http://localhost:8080/api/v1/accounts/{ACCOUNT_GUID}/balance
```

### 4. Test Transactions API

**Get transactions:**
```bash
curl "http://localhost:8080/api/v1/transactions?limit=10&offset=0"
```

**Filter by date:**
```bash
curl "http://localhost:8080/api/v1/transactions?start_date=2024-01-01&end_date=2024-12-31"
```

**Filter by account:**
```bash
curl "http://localhost:8080/api/v1/transactions?account_guid={ACCOUNT_GUID}"
```

**Search by description:**
```bash
curl "http://localhost:8080/api/v1/transactions?description=groceries"
```

### 5. Test Analytics API

**Income vs Expense:**
```bash
curl "http://localhost:8080/api/v1/analytics/income-expense"
```

**Category Breakdown:**
```bash
curl "http://localhost:8080/api/v1/analytics/category-breakdown"
```

**Net Worth:**
```bash
curl http://localhost:8080/api/v1/analytics/net-worth
```

**With date range:**
```bash
curl "http://localhost:8080/api/v1/analytics/income-expense?start_date=2024-01-01&end_date=2024-12-31"
```

## Frontend Testing

### 1. Manual Testing

Open http://localhost:3000 in your browser.

**Test Authentication:**
1. Click "Login" in the navbar
2. Register a new account
3. Login with the account
4. Verify you see your email in the navbar
5. Click "Logout"

**Test Dashboard:**
1. Navigate to the Dashboard (home page)
2. Verify Net Worth summary displays
3. Verify Income vs Expense chart displays
4. Check that numbers are formatted correctly

**Test Accounts:**
1. Navigate to "Accounts" page
2. Verify account hierarchy displays
3. Click on account folders to expand/collapse
4. Check that balances display correctly
5. Verify account types are shown

**Test Transactions:**
1. Navigate to "Transactions" page
2. Verify transaction table displays
3. Test date range filter
4. Test description search
5. Verify pagination info displays
6. Check that splits are shown correctly

**Test Analytics:**
1. Navigate to "Analytics" page
2. Verify all charts render:
   - Net Worth summary
   - Income vs Expense chart
   - Income pie chart
   - Expense pie chart
3. Test date range filters
4. Verify chart interactions (hover, tooltips)

### 2. Browser Console Testing

Open browser DevTools (F12) and check:
- No console errors
- Network requests succeed (200 status)
- API responses are correctly formatted

### 3. Responsive Design Testing

Test on different screen sizes:
- Desktop (1920x1080)
- Tablet (768x1024)
- Mobile (375x667)

Verify:
- Navigation works on mobile
- Charts are responsive
- Tables scroll horizontally on small screens
- Forms are usable

## Common Issues

### Backend won't start
- Check database credentials in `configs/config.yaml`
- Ensure PostgreSQL is running
- Verify GnuCash database exists

### Frontend can't connect to backend
- Check `.env` file has correct `REACT_APP_API_BASE_URL`
- Verify backend is running on port 8080
- Check CORS is enabled in backend

### No data displays
- Ensure your GnuCash database has data
- Check that accounts table is populated
- Verify transactions table has entries

### Authentication fails
- Check JWT secret is configured
- Verify `app_users` table was created
- Check password meets minimum length (8 chars)

### Charts don't render
- Ensure Recharts is installed: `npm list recharts`
- Check browser console for errors
- Verify analytics API returns data

## Performance Testing

### Load Testing (Optional)

Install Apache Bench:
```bash
# Ubuntu/Debian
sudo apt-get install apache2-utils

# macOS
brew install httpd
```

Test account endpoint:
```bash
ab -n 1000 -c 10 http://localhost:8080/api/v1/accounts
```

Test transaction endpoint:
```bash
ab -n 1000 -c 10 "http://localhost:8080/api/v1/transactions?limit=50"
```

Expected results:
- Accounts: < 50ms average response time
- Transactions: < 100ms average response time
- Analytics: < 500ms average response time

## Automated Testing (Future)

Consider adding:
- Backend unit tests: `go test ./...`
- Frontend unit tests: `npm test`
- E2E tests with Cypress or Playwright
- API integration tests with Postman/Newman
