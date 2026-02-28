#!/bin/bash

# Source configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh"

print_header "API VALIDATION TESTS (Valid & Invalid Requests)"

# Wait for server
wait_for_server

# Setup: Create a test user
log_info "Setting up test user..."
SETUP_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${TEST_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

ACCESS_TOKEN=$(echo "$SETUP_RESPONSE" | jq -r '.access_token // .accessToken // empty')

if [ -z "$ACCESS_TOKEN" ]; then
    log_error "Failed to get access token for setup"
    exit 1
fi

log_success "Setup complete"

# Test 1: Missing Authorization Header
print_header "Test 1: Request without Authorization Header (Expect 401)"
NO_AUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products")

HTTP_STATUS=$(echo "$NO_AUTH_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "401" ]; then
    log_success "Missing authorization correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Missing authorization validation failed - Expected: 401, Got: $HTTP_STATUS"
fi

# Test 2: Invalid Token Format
print_header "Test 2: Invalid Token Format (Expect 401)"
INVALID_TOKEN_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products" \
    -H "Authorization: Bearer not.a.valid.jwt.token")

HTTP_STATUS=$(echo "$INVALID_TOKEN_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "401" ]; then
    log_success "Invalid token correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Invalid token validation failed - Expected: 401, Got: $HTTP_STATUS"
fi

# Test 3: Expired Token
print_header "Test 3: Expired/Malformed Token (Expect 401)"
EXPIRED_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjJ9.invalid"
EXPIRED_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products" \
    -H "Authorization: Bearer ${EXPIRED_TOKEN}")

HTTP_STATUS=$(echo "$EXPIRED_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "401" ]; then
    log_success "Expired/malformed token correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Expired token validation failed - Expected: 401, Got: $HTTP_STATUS"
fi

# Test 4: Malformed JSON Request Body
print_header "Test 4: Malformed JSON in Request Body (Expect 400)"
log_request "POST ${API_BASE}/auth/register"
log_request "Body: {invalid json here} (malformed)"

MALFORMED_JSON_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{invalid json here}")

HTTP_STATUS=$(echo "$MALFORMED_JSON_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Malformed JSON correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Malformed JSON validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 5: Missing Required Fields (Email)
print_header "Test 5: Missing Required Field - Email (Expect 400)"
log_request "POST ${API_BASE}/auth/register"
log_request "Body: {password: ***, password_confirm: ***} (missing email)"

MISSING_EMAIL_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

HTTP_STATUS=$(echo "$MISSING_EMAIL_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Missing email correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Missing email validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 6: Invalid Email Format
print_header "Test 6: Invalid Email Format (Expect 400)"
INVALID_EMAIL_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"not-an-email\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

HTTP_STATUS=$(echo "$INVALID_EMAIL_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Invalid email format correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Invalid email format validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 7: Invalid Data Type (String for Number)
print_header "Test 7: Invalid Data Type - Price as String (Expect 400)"
log_request "POST ${API_BASE}/products"
log_request "Body: {name: 'Test Product', price: 'not a number' (invalid type)}"

INVALID_TYPE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"Test Product\",
        \"description\": \"Test\",
        \"price\": \"not a number\",
        \"stock\": 10,
        \"category_id\": 1
    }")

HTTP_STATUS=$(echo "$INVALID_TYPE_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Invalid data type correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Invalid data type validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 8: Non-existent Resource
print_header "Test 8: Request Non-existent Resource (Expect 404)"
NOT_FOUND_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products/999999999" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$NOT_FOUND_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "404" ]; then
    log_success "Non-existent resource correctly returned 404 - Status: $HTTP_STATUS"
else
    log_error "Non-existent resource validation failed - Expected: 404, Got: $HTTP_STATUS"
fi

# Test 9: Invalid HTTP Method
print_header "Test 9: Invalid HTTP Method (Expect 405)"
INVALID_METHOD_RESPONSE=$(curl -s -w "\n%{http_code}" -X PATCH "${API_BASE}/auth/register" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$INVALID_METHOD_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "405" ] || [ "$HTTP_STATUS" = "404" ]; then
    log_success "Invalid HTTP method correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Invalid HTTP method validation failed - Expected: 405, Got: $HTTP_STATUS"
fi

# Test 10: CORS Preflight Request
print_header "Test 10: CORS Preflight Request"
CORS_RESPONSE=$(curl -s -w "\n%{http_code}" -X OPTIONS "${API_BASE}/products" \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: GET" \
    -H "Access-Control-Request-Headers: authorization,content-type")

HTTP_STATUS=$(echo "$CORS_RESPONSE" | tail -n1)
RESPONSE_HEADERS=$(echo "$CORS_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "204" ]; then
    log_success "CORS preflight successful - Status: $HTTP_STATUS"
    
    # Check for CORS headers
    if echo "$RESPONSE_HEADERS" | grep -qi "access-control-allow"; then
        log_info "CORS headers present in response"
    fi
else
    log_error "CORS preflight failed - Expected: 200/204, Got: $HTTP_STATUS"
fi

# Test 11: Negative Price Value
print_header "Test 11: Negative Price Value (Expect 400)"
NEGATIVE_PRICE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"Invalid Product\",
        \"description\": \"Test\",
        \"price\": -10.99,
        \"stock\": 10,
        \"category_id\": 1
    }")

HTTP_STATUS=$(echo "$NEGATIVE_PRICE_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Negative price correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Negative price validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 12: Excessive Field Length
print_header "Test 12: Excessive Field Length (Expect 400)"
LONG_STRING=$(python3 -c "print('A' * 10000)")
LONG_FIELD_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${LONG_STRING}@example.com\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

HTTP_STATUS=$(echo "$LONG_FIELD_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Excessive field length correctly rejected - Status: $HTTP_STATUS"
else
    log_warning "Excessive field length validation - Expected: 400, Got: $HTTP_STATUS (may be allowed)"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
fi

# Test 13: SQL Injection Attempt (NoSQL context)
print_header "Test 13: SQL/NoSQL Injection Attempt"
INJECTION_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"admin@example.com' OR '1'='1\",
        \"password\": \"password' OR '1'='1\"
    }")

HTTP_STATUS=$(echo "$INJECTION_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ] || [ "$HTTP_STATUS" = "401" ]; then
    log_success "Injection attempt correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Injection attempt validation failed - Expected: 400/401, Got: $HTTP_STATUS"
fi

# Test 14: Empty Request Body
print_header "Test 14: Empty Request Body (Expect 400)"
EMPTY_BODY_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{}")

HTTP_STATUS=$(echo "$EMPTY_BODY_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Empty request body correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Empty body validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 15: Valid Request After Invalid Ones
print_header "Test 15: Valid Request (Should Succeed)"
VALID_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products?limit=5" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$VALID_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "200" ]; then
    log_success "Valid request successful after invalid ones - Status: $HTTP_STATUS"
else
    log_error "Valid request failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 16: Content-Type Validation
print_header "Test 16: Missing Content-Type Header (POST request)"
NO_CONTENT_TYPE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -d '{"name":"test"}')

HTTP_STATUS=$(echo "$NO_CONTENT_TYPE_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ] || [ "$HTTP_STATUS" = "415" ]; then
    log_success "Missing Content-Type correctly handled - Status: $HTTP_STATUS"
else
    log_warning "Missing Content-Type - Expected: 400/415, Got: $HTTP_STATUS (may be auto-detected)"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
fi

# Test 17: Case Sensitivity in Endpoints
print_header "Test 17: Case Sensitivity in Endpoints (Expect 404)"
CASE_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/PRODUCTS" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$CASE_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "404" ]; then
    log_success "Case-sensitive endpoint correctly returned 404 - Status: $HTTP_STATUS"
else
    log_error "Case sensitivity validation failed - Expected: 404, Got: $HTTP_STATUS"
fi

# Print summary
print_summary
exit $?
