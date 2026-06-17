# JWT Troubleshooting Notes

## Problem

JWT authentication fails when using mock-jwks and gen-token together against the update endpoint.

## Test steps and commands tried

In project root directory in two separate shells:
shell1: go run ./testdata/mock-jwks/
shell2: JWKS_URL=http://localhost:8080/.well-known/jwks.json JWT_AUDIENCE=min-app-test ALLOWED_NAMESPACES=lagring STORE=memory go run main.go     

Then from a third shell I execute:
shell3: curl -si -X POST localhost:8072/tasks -H "Content-Type: application/json" -H "Authorization: Bearer $(./gen-token)" -d '{"title": "test"}'

## Error output

HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Wed, 17 Jun 2026 14:51:06 GMT
Content-Length: 61

could not verify message using any of the signatures or keys

## Expected behavior

Should allow the creation of a new task.