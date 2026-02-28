# NoSQL E-commerce System - Test Suite

This directory contains comprehensive test scripts for the NoSQL-based e-commerce system, covering all aspects of the final assignment.

## Overview

The test suite validates the following key areas:
1. **Registration, authentication, and profile editing**
2. **Adding and searching for products**
3. **Generating and displaying recommendations**
4. **API functionality (valid and invalid requests)**

## Test Files

### `config.sh`
Central configuration file containing:
- API base URLs and endpoints
- Test user credentials
- Color codes for output formatting
- Helper functions for logging and assertions
- Test counter tracking

### `01_auth_tests.sh`
**Coverage**: Registration, authentication, and profile editing

Tests:
- User registration with valid data
- Duplicate email validation
- Login with correct credentials
- Login with wrong password
- Profile retrieval with authentication
- Profile updates (partial)
- Unauthorized access attempts
- Token refresh functionality
- Invalid token format handling

### `02_product_tests.sh`
**Coverage**: Adding and searching for products

Tests:
- Category listing
- Product creation
- Product listing with pagination
- Product search by name
- Filtering by category
- Filtering by price range
- Sorting by price
- Get specific product by ID
- Product updates (partial)
- Non-existent product handling
- Invalid product data validation
- Product deletion

### `03_recommendation_tests.sh`
**Coverage**: Generating and displaying recommendations

Tests:
- Recording product views
- Recording product likes
- Recording product purchases
- Creating similar user interaction patterns
- Interaction history retrieval
- Product statistics (views, likes, purchases)
- Similar users discovery
- Personalized recommendations
- Recommendations with limit
- Collaborative filtering weight verification (purchases=50%, likes=35%, views=15%)

### `04_api_validation_tests.sh`
**Coverage**: API functionality (valid and invalid requests)

Tests:
- Missing authorization header
- Invalid token format
- Expired/malformed tokens
- Malformed JSON requests
- Missing required fields
- Invalid email format
- Invalid data types
- Non-existent resources
- Invalid HTTP methods
- CORS preflight requests
- Negative price values
- Excessive field lengths
- SQL/NoSQL injection attempts
- Empty request bodies
- Valid requests after errors
- Content-Type validation
- Case sensitivity in endpoints

### `05_recommendation_evaluation.sh`
**Coverage**: Recommendation system effectiveness evaluation

Evaluates:
- **Precision**: Accuracy of recommendations (relevant items / total recommended)
- **Recall**: Coverage of user interests (relevant items recommended / total relevant)
- **F1-Score**: Harmonic mean of precision and recall

User Types Analyzed:
- **Heavy Buyer**: High purchase activity (5 purchases, 3 likes, 2 views)
- **Window Shopper**: High browse, low conversion (7 views, 2 likes, 1 purchase)
- **Engaged User**: Balanced interactions (3 purchases, 4 likes, 3 views)
- **New User**: Cold start scenario (1 view, 1 like)

Outputs:
- Individual metrics for each user type
- Average metrics across all user types
- Algorithm analysis and insights
- Recommendations for improvement

### `run_all_tests.sh`
Master test runner that:
- Executes all test suites in sequence
- Provides detailed progress reporting
- Generates comprehensive summary
- Returns appropriate exit codes
- Includes recommendation system evaluation

## Running Tests

### Prerequisites
- Server must be running on `localhost:8080`
- MongoDB must be accessible
- `curl` and `jq` must be installed

### Run All Tests
```bash
cd tests
chmod +x *.sh
./run_all_tests.sh
```

### Run Individual Test Suites
```bash
# Authentication tests
./01_auth_tests.sh

# Product tests
./02_product_tests.sh

# Recommendation tests
./03_recommendation_tests.sh

# API validation tests
./04_api_validation_tests.sh

# Recommendation evaluation
./05_recommendation_evaluation.sh
```

### Start the Server
```bash
# From project root
make run

# Or manually
go run cmd/web/main.go
```

## Test Output

Each test script provides:
- **Color-coded output**:
  - 🔵 Blue: Information and headers
  - 🟢 Green: Passed tests
  - 🔴 Red: Failed tests
  - 🟡 Yellow: Warnings

- **Test counters**: Total, Passed, Failed
- **Detailed responses** for failures
- **Summary report** at the end

### Example Output
```
================================================
  AUTHENTICATION & PROFILE TESTS
================================================

[INFO] Waiting for server to be ready...
[PASS] Server is ready!

================================================
  Test 1: User Registration (Valid Data)
================================================

[PASS] User registered successfully - Status: 201
[INFO] Access Token: eyJhbGciOiJIUzI1NiIs...

...

================================================
TEST SUMMARY
================================================
Total Tests:  9
Passed:       9
Failed:       0

✓ All tests passed!
================================================
```

## Implementation Details

### Collaborative Filtering
The recommendation system uses weighted scoring:
- **Purchases**: 50% weight (highest priority)
- **Likes**: 35% weight (medium priority)
- **Views**: 15% weight (lowest priority)

Tests verify that recommendations correctly reflect these weights by creating different interaction patterns and validating the results.

### Authentication
Tests validate:
- JWT token generation and validation
- Bearer token format (with and without "Bearer" prefix)
- Token refresh mechanism
- Unauthorized access protection

### Data Validation
Tests cover:
- Input validation (email format, required fields)
- Type checking (string vs number)
- Range validation (negative prices, excessive lengths)
- Security (injection attempts, malformed tokens)

### CORS
Validates CORS preflight requests for origins:
- `http://localhost:3000`
- `http://localhost:5173`
- `http://localhost:8080`

## Test Data Cleanup

Tests create temporary data with unique identifiers (timestamps) to avoid conflicts. The test users and products are created with:
- Dynamic emails: `test_user_<timestamp>@example.com`
- Unique product names: `Test Product <timestamp>`

## Troubleshooting

### Server Not Running
```bash
# Error: Server is not running
# Solution: Start the server
make run
```

### Authentication Failures
```bash
# Error: Failed to get access token
# Possible causes:
# - MongoDB not running
# - Invalid credentials in config
# Solution: Check MongoDB and config.yaml
docker-compose up -d
```

### Missing Dependencies
```bash
# Install jq (JSON processor)
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq

# Fedora
sudo dnf install jq
```

## Assignment Compliance

This test suite fulfills the final assignment requirements:

✅ **Test Case 1**: Registration, authentication, and profile editing
- Covered by `01_auth_tests.sh` (9 tests)

✅ **Test Case 2**: Adding and searching for products
- Covered by `02_product_tests.sh` (12 tests)

✅ **Test Case 3**: Generating and displaying recommendations
- Covered by `03_recommendation_tests.sh` (10 tests)

✅ **Test Case 4**: API functionality (valid and invalid requests)
- Covered by `04_api_validation_tests.sh` (17 tests)

✅ **Evaluation**: Recommendation system effectiveness
- Covered by `05_recommendation_evaluation.sh` (Precision, Recall, F1-Score)
- 4 user types analyzed with detailed metrics

**Total**: 48 comprehensive tests + recommendation evaluation across 5 test suites

## Exit Codes

- `0`: All tests passed
- `1`: One or more tests failed

## Notes

- Tests are designed to be idempotent when possible
- Each test suite is independent and can run separately
- The master runner (`run_all_tests.sh`) provides the most comprehensive validation
- Test execution includes 2-second delays between suites to prevent race conditions
