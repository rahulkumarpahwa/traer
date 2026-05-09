#!/bin/bash

BASE_URL="http://localhost:8080"

EMAIL="testuser@example.com"
USERNAME="testuser"
PASSWORD="TestPass123!"

echo "======================================"
echo "1. Health Check"
echo "======================================"

curl -X GET "$BASE_URL/"
echo -e "\n"

echo "======================================"
echo "2. Create User"
echo "======================================"

curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "'"$USERNAME"'",
    "email": "'"$EMAIL"'",
    "password": "'"$PASSWORD"'"
  }'

echo -e "\n"

echo "======================================"
echo "3. Login"
echo "======================================"

LOGIN_RESPONSE=$(curl -i -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'"$EMAIL"'",
    "password": "'"$PASSWORD"'"
  }')

echo "$LOGIN_RESPONSE"

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -i "Set-Cookie" | sed -n 's/.*token=\([^;]*\).*/\1/p')

echo -e "\nTOKEN:"
echo "$TOKEN"

echo -e "\n"

echo "======================================"
echo "4. Get User By Email"
echo "======================================"

curl -X GET "$BASE_URL/users/email?email=$EMAIL" \
  -H "Cookie: token=$TOKEN"

echo -e "\n"

echo "======================================"
echo "5. Get User By Username"
echo "======================================"

curl -X GET "$BASE_URL/users/username?username=$USERNAME" \
  -H "Cookie: token=$TOKEN"

echo -e "\n"

echo "======================================"
echo "6. Update User"
echo "======================================"

curl -X PUT "$BASE_URL/users?id=1" \
  -H "Content-Type: application/json" \
  -H "Cookie: token=$TOKEN" \
  -d '{
    "username": "updateduser"
  }'

echo -e "\n"

echo "======================================"
echo "7. Create Job"
echo "======================================"

JOB_RESPONSE=$(curl -s -X POST "$BASE_URL/jobs/create" \
  -H "Content-Type: application/json" \
  -H "Cookie: token=$TOKEN" \
  -d '{
    "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "type": "audio"
  }')

echo "$JOB_RESPONSE"

JOB_ID=$(echo "$JOB_RESPONSE" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')

echo -e "\nJOB ID:"
echo "$JOB_ID"

echo -e "\n"

echo "======================================"
echo "8. Get Active Jobs"
echo "======================================"

curl -X GET "$BASE_URL/jobs/active" \
  -H "Cookie: token=$TOKEN"

echo -e "\n"

echo "======================================"
echo "9. Get Job Status"
echo "======================================"

curl -X GET "$BASE_URL/jobs/status?id=$JOB_ID" \
  -H "Cookie: token=$TOKEN"

echo -e "\n"

echo "======================================"
echo "10. Get Instances"
echo "======================================"

curl -X GET "$BASE_URL/instances/get" \
  -H "Cookie: token=$TOKEN"

echo -e "\n"

echo "======================================"
echo "ALL TESTS COMPLETED"
echo "======================================"