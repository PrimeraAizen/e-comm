#!/bin/bash

# Master test runner for NoSQL E-commerce System

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_header() {
    echo -e "\n${BLUE}================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}================================================${NC}\n"
}

# Test summary counters
TOTAL_SUITES=0
PASSED_SUITES=0
FAILED_SUITES=0

print_header "NoSQL E-commerce System - Test Suite"
echo -e "${BLUE}Testing Assignment: NoSQL-Based System${NC}"
echo -e "Date: $(date)"
echo -e "Server: http://localhost:8080\n"

# Check if server is running
echo -n "Checking server status... "
if curl -s http://localhost:8080/ping > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Server is running${NC}"
else
    echo -e "${RED}✗ Server is not running${NC}"
    echo -e "${YELLOW}Please start the server with: make run${NC}"
    exit 1
fi

# Make all test scripts executable
chmod +x "${SCRIPT_DIR}"/*.sh 2>/dev/null

# Run test suites
echo ""

# Test Suite 1: Authentication & Profile Tests
print_header "Test Suite 1: Authentication & Profile Management"
echo "Coverage: Registration, authentication, and profile editing"
if bash "${SCRIPT_DIR}/01_auth_tests.sh"; then
    ((PASSED_SUITES++))
    echo -e "${GREEN}✓ Authentication tests passed${NC}"
else
    ((FAILED_SUITES++))
    echo -e "${RED}✗ Authentication tests failed${NC}"
fi
((TOTAL_SUITES++))

sleep 2

# Test Suite 2: Product Operations Tests
print_header "Test Suite 2: Product Operations"
echo "Coverage: Adding and searching for products"
if bash "${SCRIPT_DIR}/02_product_tests.sh"; then
    ((PASSED_SUITES++))
    echo -e "${GREEN}✓ Product tests passed${NC}"
else
    ((FAILED_SUITES++))
    echo -e "${RED}✗ Product tests failed${NC}"
fi
((TOTAL_SUITES++))

sleep 2

# Test Suite 3: Recommendation System Tests
print_header "Test Suite 3: Recommendation System"
echo "Coverage: Generating and displaying recommendations"
if bash "${SCRIPT_DIR}/03_recommendation_tests.sh"; then
    ((PASSED_SUITES++))
    echo -e "${GREEN}✓ Recommendation tests passed${NC}"
else
    ((FAILED_SUITES++))
    echo -e "${RED}✗ Recommendation tests failed${NC}"
fi
((TOTAL_SUITES++))

sleep 2

# Test Suite 4: API Validation Tests
print_header "Test Suite 4: API Validation"
echo "Coverage: API functionality (valid and invalid requests)"
if bash "${SCRIPT_DIR}/04_api_validation_tests.sh"; then
    ((PASSED_SUITES++))
    echo -e "${GREEN}✓ API validation tests passed${NC}"
else
    ((FAILED_SUITES++))
    echo -e "${RED}✗ API validation tests failed${NC}"
fi
((TOTAL_SUITES++))

sleep 2

# Test Suite 5: Recommendation System Evaluation
print_header "Test Suite 5: Recommendation System Evaluation"
echo "Coverage: Precision, Recall, F1-Score metrics for different user types"
if bash "${SCRIPT_DIR}/05_recommendation_evaluation.sh"; then
    ((PASSED_SUITES++))
    echo -e "${GREEN}✓ Recommendation evaluation completed${NC}"
else
    ((FAILED_SUITES++))
    echo -e "${RED}✗ Recommendation evaluation failed${NC}"
fi
((TOTAL_SUITES++))

# Final Summary
print_header "FINAL TEST SUMMARY"
echo -e "Total Test Suites: ${BLUE}${TOTAL_SUITES}${NC}"
echo -e "Passed Suites:     ${GREEN}${PASSED_SUITES}${NC}"
echo -e "Failed Suites:     ${RED}${FAILED_SUITES}${NC}"
echo ""

if [ $FAILED_SUITES -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  ✓ ALL TEST SUITES PASSED!            ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Test Coverage:${NC}"
    echo "  ✓ Registration, authentication, and profile editing"
    echo "  ✓ Adding and searching for products"
    echo "  ✓ Generating and displaying recommendations"
    echo "  ✓ API functionality (valid and invalid requests)"
    echo "  ✓ Recommendation metrics (Precision, Recall, F1-Score)"
    echo ""
    exit 0
else
    echo -e "${RED}╔════════════════════════════════════════╗${NC}"
    echo -e "${RED}║  ✗ SOME TEST SUITES FAILED             ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}Please review the failed tests above.${NC}"
    echo ""
    exit 1
fi
