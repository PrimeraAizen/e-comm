#!/bin/bash

# Test Configuration
export BASE_URL="http://localhost:8080"
export API_BASE="${BASE_URL}/api/v1"

# Colors for output
export RED='\033[0;31m'
export GREEN='\033[0;32m'
export YELLOW='\033[1;33m'
export BLUE='\033[0;34m'
export NC='\033[0m' # No Color

# Test counters
export TOTAL_TESTS=0
export PASSED_TESTS=0
export FAILED_TESTS=0

# Test user credentials
export TEST_EMAIL="test_user_$(date +%s)@example.com"
export TEST_PASSWORD="TestPassword123!"
export TEST_EMAIL_2="test_user_2_$(date +%s)@example.com"

# Store tokens
export ACCESS_TOKEN=""
export REFRESH_TOKEN=""

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED_TESTS++))
    ((TOTAL_TESTS++))
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_request() {
    echo -e "${BLUE}[REQUEST]${NC} $1"
}

log_response() {
    echo -e "${BLUE}[RESPONSE]${NC} $1"
}

print_separator() {
    echo -e "\n${BLUE}================================================${NC}"
}

print_header() {
    print_separator
    echo -e "${BLUE}  $1${NC}"
    print_separator
}

check_response() {
    local response=$1
    local expected_status=$2
    local test_name=$3
    
    local status=$(echo "$response" | jq -r '.status // empty')
    
    if [ -z "$status" ]; then
        # Try to get HTTP status from response
        status=$(echo "$response" | grep -o "HTTP/[0-9.]* [0-9]*" | awk '{print $2}')
    fi
    
    if [ "$status" = "$expected_status" ]; then
        log_success "$test_name - Status: $status"
        return 0
    else
        log_error "$test_name - Expected: $expected_status, Got: $status"
        echo "Response: $response" | head -n 5
        return 1
    fi
}

extract_json_value() {
    local json=$1
    local key=$2
    echo "$json" | jq -r ".$key // empty"
}

wait_for_server() {
    log_info "Waiting for server to be ready..."
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "${BASE_URL}/ping" > /dev/null 2>&1; then
            log_success "Server is ready!"
            return 0
        fi
        echo -n "."
        sleep 1
        ((attempt++))
    done
    
    log_error "Server is not responding after ${max_attempts} seconds"
    exit 1
}

print_summary() {
    print_separator
    echo -e "${BLUE}TEST SUMMARY${NC}"
    print_separator
    echo -e "Total Tests:  ${BLUE}${TOTAL_TESTS}${NC}"
    echo -e "Passed:       ${GREEN}${PASSED_TESTS}${NC}"
    echo -e "Failed:       ${RED}${FAILED_TESTS}${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
        print_separator
        return 0
    else
        echo -e "\n${RED}✗ Some tests failed${NC}"
        print_separator
        return 1
    fi
}
