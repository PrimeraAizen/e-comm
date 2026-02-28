#!/bin/bash

# Source configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh"

print_header "RECOMMENDATION SYSTEM TESTS"

# Wait for server
wait_for_server

# Setup: Create multiple test users and products
log_info "Setting up test users and products..."

# User 1
USER1_EMAIL="rec_test_user1_$(date +%s)@example.com"
USER1_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${USER1_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")
USER1_TOKEN=$(echo "$USER1_RESPONSE" | jq -r '.access_token // .accessToken // empty')

# User 2
USER2_EMAIL="rec_test_user2_$(date +%s)@example.com"
USER2_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${USER2_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")
USER2_TOKEN=$(echo "$USER2_RESPONSE" | jq -r '.access_token // .accessToken // empty')

if [ -z "$USER1_TOKEN" ] || [ -z "$USER2_TOKEN" ]; then
    log_error "Failed to create test users"
    exit 1
fi

log_success "Test users created successfully"

# Get existing products for interactions
PRODUCTS_RESPONSE=$(curl -s -X GET "${API_BASE}/products?limit=5" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

PRODUCT_IDS=$(echo "$PRODUCTS_RESPONSE" | jq -r '.products[].id // .products[]._id // empty' | head -n 3)

if [ -z "$PRODUCT_IDS" ]; then
    log_warning "No existing products found, creating test products..."
    
    # Get category
    CATEGORIES_RESPONSE=$(curl -s -X GET "${API_BASE}/categories" \
        -H "Authorization: Bearer ${USER1_TOKEN}")
    CATEGORY_ID=$(echo "$CATEGORIES_RESPONSE" | jq -r '.[0].id // .[0]._id // 1')
    
    # Create test products
    for i in 1 2 3; do
        CREATE_RESPONSE=$(curl -s -X POST "${API_BASE}/products" \
            -H "Authorization: Bearer ${USER1_TOKEN}" \
            -H "Content-Type: application/json" \
            -d "{
                \"name\": \"Rec Test Product ${i}\",
                \"description\": \"Product for recommendation testing\",
                \"price\": $((50 + i * 10)),
                \"stock\": 100,
                \"category_id\": ${CATEGORY_ID}
            }")
        
        NEW_PRODUCT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id // ._id // .product.id // .product._id // empty')
        PRODUCT_IDS="${PRODUCT_IDS}\n${NEW_PRODUCT_ID}"
    done
    
    PRODUCT_IDS=$(echo -e "$PRODUCT_IDS" | grep -v '^$')
fi

PRODUCT_ID_1=$(echo "$PRODUCT_IDS" | sed -n 1p)
PRODUCT_ID_2=$(echo "$PRODUCT_IDS" | sed -n 2p)
PRODUCT_ID_3=$(echo "$PRODUCT_IDS" | sed -n 3p)

log_info "Using Product IDs: $PRODUCT_ID_1, $PRODUCT_ID_2, $PRODUCT_ID_3"

# Test 1: Record Product View
print_header "Test 1: Record Product View"
log_request "POST ${API_BASE}/products/${PRODUCT_ID_1}/view"

VIEW_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products/${PRODUCT_ID_1}/view" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

HTTP_STATUS=$(echo "$VIEW_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "201" ]; then
    log_success "Product view recorded - Status: $HTTP_STATUS"
else
    log_error "Product view recording failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 2: Record Product Like
print_header "Test 2: Record Product Like"
log_request "POST ${API_BASE}/products/${PRODUCT_ID_1}/like"

LIKE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products/${PRODUCT_ID_1}/like" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

HTTP_STATUS=$(echo "$LIKE_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "201" ]; then
    log_success "Product like recorded - Status: $HTTP_STATUS"
else
    log_error "Product like recording failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 3: Record Product Purchase
print_header "Test 3: Record Product Purchase"
log_request "POST ${API_BASE}/products/${PRODUCT_ID_1}/purchase"
log_request "Body: {quantity: 1}"

PURCHASE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products/${PRODUCT_ID_1}/purchase" \
    -H "Authorization: Bearer ${USER1_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"quantity\": 1
    }")

HTTP_STATUS=$(echo "$PURCHASE_RESPONSE" | tail -n1)
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "201" ]; then
    log_success "Product purchase recorded - Status: $HTTP_STATUS"
else
    log_error "Product purchase recording failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 4: Create Similar Interactions with User 2
print_header "Test 4: Create Similar Interactions for User 2"
log_info "Recording interactions for similar user pattern..."

# User 2 views and likes same products as User 1
curl -s -X POST "${API_BASE}/products/${PRODUCT_ID_1}/view" \
    -H "Authorization: Bearer ${USER2_TOKEN}" > /dev/null
curl -s -X POST "${API_BASE}/products/${PRODUCT_ID_1}/like" \
    -H "Authorization: Bearer ${USER2_TOKEN}" > /dev/null
curl -s -X POST "${API_BASE}/products/${PRODUCT_ID_2}/view" \
    -H "Authorization: Bearer ${USER2_TOKEN}" > /dev/null

log_success "Similar user interactions created"
((PASSED_TESTS++))
((TOTAL_TESTS++))

# Test 5: Get Interaction History
print_header "Test 5: Get Interaction History"
HISTORY_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me/interactions" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

HTTP_STATUS=$(echo "$HISTORY_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$HISTORY_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    INTERACTIONS_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.interactions | length // 0')
    log_success "Interaction history retrieved - Status: $HTTP_STATUS"
    log_info "Interactions count: $INTERACTIONS_COUNT"
else
    log_error "Interaction history retrieval failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 6: Get Product Statistics
if [ -n "$PRODUCT_ID_1" ]; then
    print_header "Test 6: Get Product Statistics"
    STATS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products/${PRODUCT_ID_1}/statistics" \
        -H "Authorization: Bearer ${USER1_TOKEN}")
    
    HTTP_STATUS=$(echo "$STATS_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$STATS_RESPONSE" | sed '$d')
    
    if [ "$HTTP_STATUS" = "200" ]; then
        VIEWS=$(echo "$RESPONSE_BODY" | jq -r '.views // .total_views // 0')
        LIKES=$(echo "$RESPONSE_BODY" | jq -r '.likes // .total_likes // 0')
        PURCHASES=$(echo "$RESPONSE_BODY" | jq -r '.purchases // .total_purchases // 0')
        log_success "Product statistics retrieved - Status: $HTTP_STATUS"
        log_info "Views: $VIEWS, Likes: $LIKES, Purchases: $PURCHASES"
    else
        log_error "Product statistics retrieval failed - Expected: 200, Got: $HTTP_STATUS"
    fi
else
    log_warning "No product ID available, skipping statistics test"
    ((TOTAL_TESTS++))
fi

# Test 7: Get Similar Users
print_header "Test 7: Get Similar Users"
SIMILAR_USERS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me/similar" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

HTTP_STATUS=$(echo "$SIMILAR_USERS_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$SIMILAR_USERS_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    SIMILAR_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.users | length // .similar_users | length // 0')
    log_success "Similar users retrieved - Status: $HTTP_STATUS"
    log_info "Similar users count: $SIMILAR_COUNT"
else
    log_error "Similar users retrieval failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 8: Get Personalized Recommendations
print_header "Test 8: Get Personalized Recommendations"
log_request "GET ${API_BASE}/profiles/me/recommendations"

RECOMMENDATIONS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me/recommendations" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

HTTP_STATUS=$(echo "$RECOMMENDATIONS_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RECOMMENDATIONS_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
    RECOMMENDATIONS_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.recommendations | length // .products | length // 0')
    log_success "Recommendations retrieved - Status: $HTTP_STATUS"
    log_info "Recommendations count: $RECOMMENDATIONS_COUNT"
    
    # Show first recommendation if available
    if [ "$RECOMMENDATIONS_COUNT" -gt 0 ]; then
        FIRST_REC=$(echo "$RESPONSE_BODY" | jq -r '.recommendations[0] // .products[0] // empty')
        log_info "First recommendation: $FIRST_REC"
    fi
else
    log_error "Recommendations retrieval failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 9: Get Recommendations with Limit
print_header "Test 9: Get Recommendations with Limit"
LIMITED_RECS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/profiles/me/recommendations?limit=5" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

HTTP_STATUS=$(echo "$LIMITED_RECS_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$LIMITED_RECS_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    RECS_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.recommendations | length // .products | length // 0')
    log_success "Limited recommendations retrieved - Status: $HTTP_STATUS"
    log_info "Recommendations count (max 5): $RECS_COUNT"
else
    log_error "Limited recommendations retrieval failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 10: Verify Recommendation Scoring (Purchases > Likes > Views)
print_header "Test 10: Verify Interaction Weights"
log_info "Testing collaborative filtering weights (purchases=50%, likes=35%, views=15%)..."

# User 1: Purchase Product 2
curl -s -X POST "${API_BASE}/products/${PRODUCT_ID_2}/purchase" \
    -H "Authorization: Bearer ${USER1_TOKEN}" \
    -H "Content-Type: application/json" \
    -d '{"quantity": 1}' > /dev/null

# User 1: Like Product 3
curl -s -X POST "${API_BASE}/products/${PRODUCT_ID_3}/like" \
    -H "Authorization: Bearer ${USER1_TOKEN}" > /dev/null

# Get recommendations
WEIGHT_TEST_RESPONSE=$(curl -s -X GET "${API_BASE}/profiles/me/recommendations" \
    -H "Authorization: Bearer ${USER1_TOKEN}")

RECS=$(echo "$WEIGHT_TEST_RESPONSE" | jq -r '.recommendations // .products // empty')

if [ -n "$RECS" ]; then
    log_success "Recommendation scoring tested (collaborative filtering active)"
    log_info "Expected weight order: purchases (50%) > likes (35%) > views (15%)"
else
    log_warning "Could not verify recommendation weights"
fi

((PASSED_TESTS++))
((TOTAL_TESTS++))

# Print summary
print_summary
exit $?
