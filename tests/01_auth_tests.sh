#!/bin/bash

# Source configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh"

print_header "AUTHENTICATION & PROFILE TESTS"

# Wait for server
wait_for_server

# Test 1: User Registration
print_header "Test 1: User Registration (Valid Data)"
log_request "POST ${API_BASE}/auth/register"
log_request "Body: {email: ${TEST_EMAIL}, password: ***, password_confirm: ***}"

REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${TEST_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

HTTP_STATUS=$(echo "$REGISTER_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "201" ] || [ "$HTTP_STATUS" = "200" ]; then
    ACCESS_TOKEN=$(echo "$RESPONSE_BODY" | jq -r '.access_token // .accessToken // empty')
    REFRESH_TOKEN=$(echo "$RESPONSE_BODY" | jq -r '.refresh_token // .refreshToken // empty')
    log_success "User registered successfully - Status: $HTTP_STATUS"
    log_info "Access Token: ${ACCESS_TOKEN:0:20}..."
else
    log_error "User registration failed - Expected: 201, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 2: Duplicate Email Registration
print_header "Test 2: Duplicate Email Registration (Expect 409)"
log_request "POST ${API_BASE}/auth/register (duplicate email)"

DUPLICATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${TEST_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

HTTP_STATUS=$(echo "$DUPLICATE_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "409" ] || [ "$HTTP_STATUS" = "400" ]; then
    log_success "Duplicate email correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Duplicate email validation failed - Expected: 409, Got: $HTTP_STATUS"
fi

# Test 3: Login with Correct Credentials
print_header "Test 3: Login with Correct Credentials"
log_request "POST ${API_BASE}/auth/login"
log_request "Body: {email: ${TEST_EMAIL}, password: ***}"

LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${TEST_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\"
    }")

HTTP_STATUS=$(echo "$LOGIN_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
    ACCESS_TOKEN=$(echo "$RESPONSE_BODY" | jq -r '.access_token // .accessToken // empty')
    REFRESH_TOKEN=$(echo "$RESPONSE_BODY" | jq -r '.refresh_token // .refreshToken // empty')
    log_success "Login successful - Status: $HTTP_STATUS"
    log_info "New Access Token: ${ACCESS_TOKEN:0:20}..."
else
    log_error "Login failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 4: Login with Wrong Password
print_header "Test 4: Login with Wrong Password (Expect 401)"
WRONG_LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${TEST_EMAIL}\",
        \"password\": \"WrongPassword123!\"
    }")

HTTP_STATUS=$(echo "$WRONG_LOGIN_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "401" ]; then
    log_success "Wrong password correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Wrong password validation failed - Expected: 401, Got: $HTTP_STATUS"
fi

# Test 5: Get Profile with Authentication
print_header "Test 5: Get Profile with Authentication"
log_request "GET ${API_BASE}/profiles/me"
log_request "Headers: Authorization: Bearer ${ACCESS_TOKEN:0:20}..."

PROFILE_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$PROFILE_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$PROFILE_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
    PROFILE_EMAIL=$(echo "$RESPONSE_BODY" | jq -r '.email // empty')
    log_success "Profile retrieved successfully - Status: $HTTP_STATUS"
    log_info "Profile Email: $PROFILE_EMAIL"
else
    log_error "Profile retrieval failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 6: Update Profile
print_header "Test 6: Update Profile (Partial Update)"
log_request "PUT ${API_BASE}/profiles/me"
log_request "Body: {first_name: 'Test', last_name: 'User Updated'}"

UPDATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "${API_BASE}/profiles/me" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"first_name\": \"Test\",
        \"last_name\": \"User Updated\"
    }")

HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$UPDATE_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"
if [ "$HTTP_STATUS" = "200" ]; then
    log_response "Body: $(echo "$RESPONSE_BODY" | jq -c '{first_name, last_name, email}')"
fi

if [ "$HTTP_STATUS" = "200" ]; then
    UPDATED_NAME=$(echo "$RESPONSE_BODY" | jq -r '.first_name // empty')
    log_success "Profile updated successfully - Status: $HTTP_STATUS"
    log_info "Updated Name: $UPDATED_NAME"
else
    log_error "Profile update failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 7: Access Protected Endpoint without Token
print_header "Test 7: Access Protected Endpoint without Token (Expect 401)"
NO_AUTH_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me")

HTTP_STATUS=$(echo "$NO_AUTH_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "401" ]; then
    log_success "Unauthorized access correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Unauthorized access validation failed - Expected: 401, Got: $HTTP_STATUS"
fi

# Test 8: Token Refresh
print_header "Test 8: Token Refresh with Valid Refresh Token"
if [ -n "$REFRESH_TOKEN" ]; then
    REFRESH_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/auth/refresh" \
        -H "Content-Type: application/json" \
        -d "{
            \"refresh_token\": \"${REFRESH_TOKEN}\"
        }")
    
    HTTP_STATUS=$(echo "$REFRESH_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$REFRESH_RESPONSE" | sed '$d')
    
    if [ "$HTTP_STATUS" = "200" ]; then
        NEW_ACCESS_TOKEN=$(echo "$RESPONSE_BODY" | jq -r '.access_token // .accessToken // empty')
        log_success "Token refreshed successfully - Status: $HTTP_STATUS"
        log_info "New Access Token: ${NEW_ACCESS_TOKEN:0:20}..."
        ACCESS_TOKEN=$NEW_ACCESS_TOKEN
    else
        log_error "Token refresh failed - Expected: 200, Got: $HTTP_STATUS"
        echo "Response: $RESPONSE_BODY"
    fi
else
    log_warning "No refresh token available, skipping test"
    ((TOTAL_TESTS++))
fi

# Test 9: Invalid Token Format
print_header "Test 9: Access with Invalid Token Format (Expect 401)"
INVALID_TOKEN_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me" \
    -H "Authorization: Bearer invalid.token.here")

HTTP_STATUS=$(echo "$INVALID_TOKEN_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "401" ]; then
    log_success "Invalid token correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Invalid token validation failed - Expected: 401, Got: $HTTP_STATUS"
fi

# Print summary
print_summary
exit $?
