#!/bin/bash

# Source configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh"

print_header "PRODUCT OPERATIONS TESTS"

# Wait for server
wait_for_server

# Setup: Register and login a test user
log_info "Setting up test user..."
REGISTER_RESPONSE=$(curl -s -X POST "${API_BASE}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${TEST_EMAIL}\",
        \"password\": \"${TEST_PASSWORD}\",
        \"password_confirm\": \"${TEST_PASSWORD}\"
    }")

ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token // .accessToken // empty')

if [ -z "$ACCESS_TOKEN" ]; then
    log_error "Failed to get access token for setup"
    exit 1
fi

log_success "Setup complete - Access Token obtained"

# Test 1: List Categories
print_header "Test 1: List Categories"
CATEGORIES_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/categories" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$CATEGORIES_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$CATEGORIES_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    CATEGORY_ID=$(echo "$RESPONSE_BODY" | jq -r '.[0].id // .[0]._id // empty')
    CATEGORY_NAME=$(echo "$RESPONSE_BODY" | jq -r '.[0].name // empty')
    log_success "Categories listed successfully - Status: $HTTP_STATUS"
    log_info "First Category ID: $CATEGORY_ID, Name: $CATEGORY_NAME"
else
    log_error "Categories listing failed - Expected: 200, Got: $HTTP_STATUS"
    # Try to get first category anyway from DB
    CATEGORY_ID=1
fi

# Test 2: Create Product
print_header "Test 2: Create Product"
PRODUCT_NAME="Test Product $(date +%s)"

if [ -z "$CATEGORY_ID" ] || [ "$CATEGORY_ID" = "null" ]; then
    CATEGORY_ID=1
fi

log_request "POST ${API_BASE}/products"
log_request "Body: {name: '${PRODUCT_NAME}', price: 99.99, stock: 100, category_id: ${CATEGORY_ID}}"

CREATE_PRODUCT_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"${PRODUCT_NAME}\",
        \"description\": \"Test product description\",
        \"price\": 99.99,
        \"stock\": 100,
        \"category_id\": ${CATEGORY_ID}
    }")

HTTP_STATUS=$(echo "$CREATE_PRODUCT_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$CREATE_PRODUCT_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "201" ] || [ "$HTTP_STATUS" = "200" ]; then
    PRODUCT_ID=$(echo "$RESPONSE_BODY" | jq -r '.id // ._id // .product.id // .product._id // empty')
    log_success "Product created successfully - Status: $HTTP_STATUS"
    log_info "Product ID: $PRODUCT_ID"
else
    log_error "Product creation failed - Expected: 201, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 3: List Products with Pagination
print_header "Test 3: List Products with Pagination"
log_request "GET ${API_BASE}/products?page=1&limit=10"

LIST_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products?page=1&limit=10" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$LIST_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$LIST_RESPONSE" | sed '$d')
log_response "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
    PRODUCTS_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.products | length // 0')
    log_success "Products listed successfully - Status: $HTTP_STATUS"
    log_info "Products count: $PRODUCTS_COUNT"
else
    log_error "Products listing failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 4: Search Products by Name
print_header "Test 4: Search Products by Name"
SEARCH_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products?search=Test" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$SEARCH_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$SEARCH_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    SEARCH_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.products | length // 0')
    log_success "Product search successful - Status: $HTTP_STATUS"
    log_info "Search results count: $SEARCH_COUNT"
else
    log_error "Product search failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 5: Filter Products by Category
print_header "Test 5: Filter Products by Category"
FILTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products?category_id=${CATEGORY_ID}" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$FILTER_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$FILTER_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    FILTERED_COUNT=$(echo "$RESPONSE_BODY" | jq -r '.products | length // 0')
    log_success "Product filtering successful - Status: $HTTP_STATUS"
    log_info "Filtered products count: $FILTERED_COUNT"
else
    log_error "Product filtering failed - Expected: 200, Got: $HTTP_STATUS"
    echo "Response: $RESPONSE_BODY"
fi

# Test 6: Filter by Price Range
print_header "Test 6: Filter Products by Price Range"
PRICE_FILTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products?min_price=50&max_price=150" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$PRICE_FILTER_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$PRICE_FILTER_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    log_success "Price range filtering successful - Status: $HTTP_STATUS"
else
    log_error "Price range filtering failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 7: Sort Products by Price
print_header "Test 7: Sort Products by Price (Ascending)"
SORT_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products?sort_by=price&order=asc" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$SORT_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$SORT_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    log_success "Product sorting successful - Status: $HTTP_STATUS"
else
    log_error "Product sorting failed - Expected: 200, Got: $HTTP_STATUS"
fi

# Test 8: Get Specific Product
if [ -n "$PRODUCT_ID" ]; then
    print_header "Test 8: Get Specific Product by ID"
    GET_PRODUCT_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products/${PRODUCT_ID}" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}")
    
    HTTP_STATUS=$(echo "$GET_PRODUCT_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$GET_PRODUCT_RESPONSE" | sed '$d')
    
    if [ "$HTTP_STATUS" = "200" ]; then
        PRODUCT_NAME_RESP=$(echo "$RESPONSE_BODY" | jq -r '.name // .product.name // empty')
        log_success "Product retrieved successfully - Status: $HTTP_STATUS"
        log_info "Product Name: $PRODUCT_NAME_RESP"
    else
        log_error "Product retrieval failed - Expected: 200, Got: $HTTP_STATUS"
        echo "Response: $RESPONSE_BODY"
    fi
else
    log_warning "No product ID available, skipping get test"
    ((TOTAL_TESTS++))
fi

# Test 9: Update Product
if [ -n "$PRODUCT_ID" ]; then
    print_header "Test 9: Update Product (Partial Update)"
    UPDATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "${API_BASE}/products/${PRODUCT_ID}" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}" \
        -H "Content-Type: application/json" \
        -d "{
            \"price\": 149.99,
            \"stock\": 75
        }")
    
    HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$UPDATE_RESPONSE" | sed '$d')
    
    if [ "$HTTP_STATUS" = "200" ]; then
        NEW_PRICE=$(echo "$RESPONSE_BODY" | jq -r '.price // .product.price // empty')
        log_success "Product updated successfully - Status: $HTTP_STATUS"
        log_info "New Price: $NEW_PRICE"
    else
        log_error "Product update failed - Expected: 200, Got: $HTTP_STATUS"
        echo "Response: $RESPONSE_BODY"
    fi
else
    log_warning "No product ID available, skipping update test"
    ((TOTAL_TESTS++))
fi

# Test 10: Get Non-existent Product
print_header "Test 10: Get Non-existent Product (Expect 404)"
NOT_FOUND_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE}/products/999999" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")

HTTP_STATUS=$(echo "$NOT_FOUND_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "404" ]; then
    log_success "Non-existent product correctly returned 404 - Status: $HTTP_STATUS"
else
    log_error "Non-existent product validation failed - Expected: 404, Got: $HTTP_STATUS"
fi

# Test 11: Create Product with Missing Fields
print_header "Test 11: Create Product with Missing Fields (Expect 400)"
INVALID_PRODUCT_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE}/products" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"Incomplete Product\"
    }")

HTTP_STATUS=$(echo "$INVALID_PRODUCT_RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" = "400" ]; then
    log_success "Invalid product data correctly rejected - Status: $HTTP_STATUS"
else
    log_error "Invalid product validation failed - Expected: 400, Got: $HTTP_STATUS"
fi

# Test 12: Delete Product
if [ -n "$PRODUCT_ID" ]; then
    print_header "Test 12: Delete Product"
    DELETE_RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "${API_BASE}/products/${PRODUCT_ID}" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}")
    
    HTTP_STATUS=$(echo "$DELETE_RESPONSE" | tail -n1)
    
    if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "204" ]; then
        log_success "Product deleted successfully - Status: $HTTP_STATUS"
    else
        log_error "Product deletion failed - Expected: 200/204, Got: $HTTP_STATUS"
    fi
else
    log_warning "No product ID available, skipping delete test"
    ((TOTAL_TESTS++))
fi

# Print summary
print_summary
exit $?
