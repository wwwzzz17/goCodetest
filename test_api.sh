#!/bin/bash

# Simple API Test Script
BASE_URL="http://localhost:8080"

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Simple API Testing${NC}\n"

# 1. Create 4 products
echo -e "${YELLOW}=== Creating 4 Products ===${NC}"

echo "1. Creating Coffee"
curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name": "Coffee", "price": 12.99, "quantity": 100}'
echo ""

echo "2. Creating Tea"
curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name": "Tea", "price": 8.50, "quantity": 50}'
echo ""

echo "3. Creating Green Tea"
curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name": "Green Tea", "price": 9.99, "quantity": 75}'
echo ""

echo "4. Creating Coffee Beans"
curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name": "Coffee Beans", "price": 15.99, "quantity": 30}'
echo -e "\n"

# 2. Test search functionality
echo -e "${YELLOW}=== Testing Search Functionality ===${NC}"

echo "Search for 'coffee':"
curl -s "$BASE_URL/products?name=coffee"
echo -e "\n"

echo "Search for 'tea':"
curl -s "$BASE_URL/products?name=tea"
echo -e "\n"

echo "Price range search (8-10):"
curl -s "$BASE_URL/products?price_min=8&price_max=10"
echo -e "\n"

echo "Quantity range search (quantity > 50):"
curl -s "$BASE_URL/products?quantity_min=50"
echo -e "\n"

# 3. Update product
echo -e "${YELLOW}=== Updating Product ===${NC}"

echo "Updating product ID=1 price:"
curl -s -X PUT "$BASE_URL/products/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Premium Coffee", "price": 16.99, "quantity": 80}'
echo -e "\n"

# 4. Delete product
echo -e "${YELLOW}=== Deleting Product ===${NC}"

echo "Deleting product ID=4:"
curl -s -X DELETE "$BASE_URL/products/4"
echo -e "\n"

# 5. View final results
echo -e "${YELLOW}=== Final Product List ===${NC}"
curl -s "$BASE_URL/products"
echo -e "\n"

echo -e "${GREEN}✅"