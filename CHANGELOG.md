# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-02-03

### Added
- Initial release of Agentic Cash
- Clean architecture Go backend with RESTful API
- React TypeScript frontend with Tailwind CSS
- JWT authentication system
- Account hierarchy visualization
- Transaction viewing with filters
- Financial analytics dashboard
- Income vs Expense charts
- Category breakdown pie charts
- Net worth calculation
- Docker support with docker-compose
- Comprehensive documentation (README, DOCKER, DEPLOYMENT, TESTING, QUICKSTART)
- Makefile for common operations
- Health check endpoints
- CORS support
- Pagination support for transactions
- Date range filtering
- PostgreSQL integration with GnuCash database
- Recharts visualizations
- TanStack Query for server state
- Protected routes (optional)
- Token refresh mechanism
- Responsive design for mobile/tablet/desktop

### Backend Features
- `/health` - Health check endpoint
- `/api/v1/auth/*` - Authentication endpoints (register, login, refresh, logout)
- `/api/v1/accounts/*` - Account management endpoints
- `/api/v1/transactions/*` - Transaction query endpoints
- `/api/v1/analytics/*` - Analytics endpoints

### Frontend Features
- Dashboard page with net worth and income/expense overview
- Accounts page with interactive hierarchy tree
- Transactions page with filterable table
- Analytics page with comprehensive charts
- Login and registration pages
- User profile display in navbar
- Loading states and error handling
- Responsive navigation

### Infrastructure
- Docker Compose configuration
- Multi-stage Docker builds for optimization
- Nginx configuration for frontend
- Environment variable support
- Database initialization scripts
- Health checks for all services
- Development and production configurations
- Volume management for data persistence

### Documentation
- Comprehensive README with quick start
- Docker deployment guide
- Production deployment instructions
- Testing and validation guide
- Quick start guide
- API documentation
- Makefile with common commands

[1.0.0]: https://github.com/udai-kiran/agentic-cash/releases/tag/v1.0.0
