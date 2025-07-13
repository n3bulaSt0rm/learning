#!/bin/bash

# Demo script for microservices
set -e

BASE_URL="http://localhost:8080"
echo "🚀 Starting Microservices Demo"
echo "================================"

# Check if services are running
echo "🔍 Checking service health..."
curl -s "$BASE_URL/health" | jq . || {
    echo "❌ API Gateway is not running. Please start services first:"
    echo "   make docker-up"
    echo "   OR"
    echo "   make run-user && make run-product && make run-order && make run-gateway"
    exit 1
}

echo "✅ API Gateway is healthy"
echo ""

# Step 1: Create users
echo "👥 Step 1: Creating users..."
USER1_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890"
  }')

USER1_ID=$(echo "$USER1_RESPONSE" | jq -r '.user.id')
echo "✅ Created user: $USER1_ID"

USER2_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith", 
    "email": "jane@example.com",
    "phone": "+0987654321"
  }')

USER2_ID=$(echo "$USER2_RESPONSE" | jq -r '.user.id')
echo "✅ Created user: $USER2_ID"
echo ""

# Step 2: Create products
echo "📦 Step 2: Creating products..."
PRODUCT1_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "price": 999.99,
    "stock": 50,
    "category": "Electronics"
  }')

PRODUCT1_ID=$(echo "$PRODUCT1_RESPONSE" | jq -r '.product.id')
echo "✅ Created product: iPhone 15 ($PRODUCT1_ID)"

PRODUCT2_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "Professional laptop",
    "price": 1999.99,
    "stock": 25,
    "category": "Electronics"
  }')

PRODUCT2_ID=$(echo "$PRODUCT2_RESPONSE" | jq -r '.product.id')
echo "✅ Created product: MacBook Pro ($PRODUCT2_ID)"

PRODUCT3_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nike Air Max",
    "description": "Comfortable running shoes",
    "price": 129.99,
    "stock": 100,
    "category": "Fashion"
  }')

PRODUCT3_ID=$(echo "$PRODUCT3_RESPONSE" | jq -r '.product.id')
echo "✅ Created product: Nike Air Max ($PRODUCT3_ID)"
echo ""

# Step 3: List products
echo "📋 Step 3: Listing products..."
curl -s "$BASE_URL/api/v1/products?page=1&page_size=10" | jq '.products[] | {id: .id, name: .name, price: .price, stock: .stock}'
echo ""

# Step 4: Create orders
echo "🛒 Step 4: Creating orders..."
ORDER1_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER1_ID\",
    \"items\": [
      {
        \"product_id\": \"$PRODUCT1_ID\",
        \"quantity\": 2
      },
      {
        \"product_id\": \"$PRODUCT3_ID\",
        \"quantity\": 1
      }
    ]
  }")

ORDER1_ID=$(echo "$ORDER1_RESPONSE" | jq -r '.order.id')
echo "✅ Created order for John: $ORDER1_ID"

ORDER2_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER2_ID\",
    \"items\": [
      {
        \"product_id\": \"$PRODUCT2_ID\",
        \"quantity\": 1
      }
    ]
  }")

ORDER2_ID=$(echo "$ORDER2_RESPONSE" | jq -r '.order.id')
echo "✅ Created order for Jane: $ORDER2_ID"
echo ""

# Step 5: Check updated stock
echo "📊 Step 5: Checking updated stock..."
echo "Product stock after orders:"
curl -s "$BASE_URL/api/v1/products/$PRODUCT1_ID" | jq '{name: .product.name, stock: .product.stock}'
curl -s "$BASE_URL/api/v1/products/$PRODUCT2_ID" | jq '{name: .product.name, stock: .product.stock}'
curl -s "$BASE_URL/api/v1/products/$PRODUCT3_ID" | jq '{name: .product.name, stock: .product.stock}'
echo ""

# Step 6: Get order details
echo "📄 Step 6: Order details..."
echo "Order 1 details:"
curl -s "$BASE_URL/api/v1/orders/$ORDER1_ID" | jq '.order | {id: .id, user_id: .user_id, total_amount: .total_amount, status: .status, items: .items}'
echo ""

echo "Order 2 details:"
curl -s "$BASE_URL/api/v1/orders/$ORDER2_ID" | jq '.order | {id: .id, user_id: .user_id, total_amount: .total_amount, status: .status, items: .items}'
echo ""

# Step 7: List user orders
echo "👤 Step 7: User orders..."
echo "John's orders:"
curl -s "$BASE_URL/api/v1/users/$USER1_ID/orders?page=1&page_size=10" | jq '.orders[] | {id: .id, total_amount: .total_amount, status: .status}'
echo ""

echo "Jane's orders:"
curl -s "$BASE_URL/api/v1/users/$USER2_ID/orders?page=1&page_size=10" | jq '.orders[] | {id: .id, total_amount: .total_amount, status: .status}'
echo ""

# Step 8: Update order status
echo "🔄 Step 8: Updating order status..."
curl -s -X PUT "$BASE_URL/api/v1/orders/$ORDER1_ID/status" \
  -H "Content-Type: application/json" \
  -d '{"status": "ORDER_STATUS_CONFIRMED"}' | jq '.order | {id: .id, status: .status}'
echo ""

# Step 9: List all orders
echo "📈 Step 9: All orders in system..."
curl -s "$BASE_URL/api/v1/orders?page=1&page_size=10" | jq '.orders[] | {id: .id, user_id: .user_id, total_amount: .total_amount, status: .status}'
echo ""

# Step 10: Test error cases
echo "❌ Step 10: Testing error cases..."
echo "Trying to create order with invalid user:"
curl -s -X POST "$BASE_URL/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "invalid-user-id",
    "items": [
      {
        "product_id": "'$PRODUCT1_ID'",
        "quantity": 1
      }
    ]
  }' | jq '.message // .error // .'

echo ""
echo "Trying to create order with insufficient stock:"
curl -s -X POST "$BASE_URL/api/v1/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER1_ID\",
    \"items\": [
      {
        \"product_id\": \"$PRODUCT1_ID\",
        \"quantity\": 1000
      }
    ]
  }" | jq '.message // .error // .'

echo ""
echo "🎉 Demo completed successfully!"
echo "================================"
echo "Summary:"
echo "- Created 2 users"
echo "- Created 3 products" 
echo "- Created 2 orders"
echo "- Verified stock updates"
echo "- Demonstrated inter-service communication"
echo "- Tested error handling"
echo ""
echo "You can explore more endpoints using the API documentation in README.md" 