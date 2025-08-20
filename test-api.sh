#!/bin/bash

# API Testing Script for Inventory Management API
# Usage: ./scripts/test-api.sh [BASE_URL]

set -e

# Configuration
BASE_URL=${1:-"http://localhost:8080"}
TEST_EMAIL="test@example.com"
TEST_PASSWORD="test123456"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️ $1${NC}"
}

# Test health endpoint
test_health() {
    print_header "Testing Health Endpoint"
    
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Health check passed"
        echo "Response: $body"
    else
        print_error "Health check failed (HTTP $http_code)"
        exit 1
    fi
}

# Test user registration
test_registration() {
    print_header "Testing User Registration"
    
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "201" ]; then
        print_success "User registration successful"
        echo "Response: $body"
    elif [ "$http_code" = "409" ]; then
        print_warning "User already exists (expected if running multiple times)"
    else
        print_error "User registration failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

# Test user login and get token
test_login() {
    print_header "Testing User Login"
    
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "User login successful"
        # Extract token
        TOKEN=$(echo "$body" | jq -r '.token')
        if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
            print_success "Token extracted: ${TOKEN:0:20}..."
            echo "$TOKEN" > /tmp/api_token
        else
            print_error "Failed to extract token"
            exit 1
        fi
    else
        print_error "User login failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

# Test getting all products
test_get_products() {
    print_header "Testing Get All Products"
    
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/products")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Get products successful"
        total=$(echo "$body" | jq -r '.total')
        echo "Total products: $total"
    else
        print_error "Get products failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

# Test creating a product (requires auth)
test_create_product() {
    print_header "Testing Create Product (Auth Required)"
    
    if [ ! -f /tmp/api_token ]; then
        print_error "No auth token found. Please run login test first."
        exit 1
    fi
    
    TOKEN=$(cat /tmp/api_token)
    
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/products" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "name": "Test Product API",
            "description": "Product created via API test",
            "quantity": 10,
            "price": 99.99,
            "category": "Test"
        }')
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "201" ]; then
        print_success "Product creation successful"
        PRODUCT_ID=$(echo "$body" | jq -r '.product.id')
        echo "Product ID: $PRODUCT_ID"
        echo "$PRODUCT_ID" > /tmp/test_product_id
    else
        print_error "Product creation failed (HTTP $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

# Test getting a single product
test_get_single_product() {
    print_header "Testing Get Single Product"
    
    if [ ! -f /tmp/test_product_id ]; then
        print_warning "No test product ID found. Skipping single product test."
        return
    fi
    
    PRODUCT_ID=$(cat /tmp/test_product_id)
    
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/products/$PRODUCT_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Get single product successful"
        name=$(echo "$body" | jq -r '.product.name')
        echo "Product name: $name"
    else
        print_error "Get single product failed (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Test updating a product (requires auth)
test_update_product() {
    print_header "Testing Update Product (Auth Required)"
    
    if [ ! -f /tmp/test_product_id ] || [ ! -f /tmp/api_token ]; then
        print_warning "Missing test product ID or token. Skipping update test."
        return
    fi
    
    PRODUCT_ID=$(cat /tmp/test_product_id)
    TOKEN=$(cat /tmp/api_token)
    
    response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/products/$PRODUCT_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "name": "Updated Test Product API",
            "description": "Updated product via API test",
            "quantity": 15,
            "price": 149.99,
            "category": "Test Updated"
        }')
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Product update successful"
    else
        print_error "Product update failed (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Test low stock products
test_low_stock() {
    print_header "Testing Low Stock Products"
    
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/products/low-stock?threshold=10")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Low stock products retrieved"
        total=$(echo "$body" | jq -r '.total')
        echo "Low stock products: $total"
    else
        print_error "Low stock products failed (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Test alerts with concurrency
test_alerts() {
    print_header "Testing Concurrent Alerts (Auth Required)"
    
    if [ ! -f /tmp/api_token ]; then
        print_warning "No auth token found. Skipping alerts test."
        return
    fi
    
    TOKEN=$(cat /tmp/api_token)
    
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/products/alerts?threshold=10" \
        -H "Authorization: Bearer $TOKEN")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Alerts generation successful"
        total=$(echo "$body" | jq -r '.total')
        echo "Alerts generated: $total"
    else
        print_error "Alerts generation failed (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Test deleting the test product (requires auth)
test_delete_product() {
    print_header "Testing Delete Product (Auth Required)"
    
    if [ ! -f /tmp/test_product_id ] || [ ! -f /tmp/api_token ]; then
        print_warning "Missing test product ID or token. Skipping delete test."
        return
    fi
    
    PRODUCT_ID=$(cat /tmp/test_product_id)
    TOKEN=$(cat /tmp/api_token)
    
    response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/products/$PRODUCT_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        print_success "Product deletion successful"
    else
        print_error "Product deletion failed (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Cleanup function
cleanup() {
    print_header "Cleaning Up"
    rm -f /tmp/api_token /tmp/test_product_id
    print_success "Cleanup completed"
}

# Main test execution
main() {
    echo -e "${BLUE}🚀 Starting API Tests for: $BASE_URL${NC}\n"
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        print_error "jq is required but not installed. Please install jq to run this script."
        exit 1
    fi
    
    # Run tests
    test_health
    test_registration
    test_login
    test_get_products
    test_create_product
    test_get_single_product
    test_update_product
    test_low_stock
    test_alerts
    test_delete_product
    
    cleanup
    
    echo -e "\n${GREEN}🎉 All tests completed successfully!${NC}"
}

# Run main function
main "$@"