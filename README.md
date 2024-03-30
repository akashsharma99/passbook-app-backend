` 23-03-2024 `

# Passbook App API

Plan to create api using golang for the passbook app [frontend](https://github.com/akashsharma99/passbook-app-web)

click [here](NOTES.md) for entity and functional requirement list.

base URL: `http://api.domain.app/v1/`

## Authentication

```
Authorization: Bearer <token>
```

## Generic Responses

### Success Response

```json
{
    "status": "success",
    "message": "success message",
    "data": {},
    "meta": {
        "total_pages": 100,
        "page": 1,
        "limit": 10
    }
}
```

### Error Response

```json
{
    "status": "error",
    "message": "error message",
    "data": {},
    "meta": {}
}
```
## Endpoints


#### `POST /auth/register` - Register User

**Request**

```json
{
    "username": "john_doe",
    "email": "john.doe@email.com",
    "password": "password"
}
```
**Responses**
- 201: User created successfully
```json
{
    "status": "success",
    "message": "User created successfully",
    "data": {
        "username": "john_doe",
        "email": "john_doe@email.com"
    }
}
```
- 400: Validation error
- 409: User already exists
- 500: Internal failures

#### `POST /auth/login` - Login User

**Request**

```json
{
    "username": "john.doe",
    "password": "password"
}
```
**Responses**
- 401: Unauthorized
- 200: User logged in successfully
```json
{
    "status": "success",
    "message": "User logged in successfully",
    "data":{
        "access_token": "jwt_token",
        "user": {
            "username": "john_doe",
            "email": "john.doe@email.com"
        }
    }
}
```
> Use `access_token` in Authorization header for further requests. The `refresh_token` is part of a secure httpOnly cookie sent by the server. Both tokens are jwt tokens. The `access_token` is short-lived (15 minutes) and the `refresh_token` is long-lived (24 hours).

**jwt token payload**
```json
{
  "userId": "a7437e4e-a898-4c17-b96e-acf32754ae6e", // uuid
  "exp": 1711901181,// expiry time
  "iat": 1711814781 // issued at
}
```
> All 🔒 endpoints require a Authorization header

#### `GET /users/me` 🔒 - Get User Profile
```bash
curl -X GET http://api.domain.app/v1/users/me -H "Authorization : Bearer <token>"
```


#### `POST /passbooks` 🔒 - Create Passbook

**Request**

```json
{
    "bank_name": "HDFC Bank",
    "account_number": "1234567890"
}
```
**Responses**
- 201: Passbook created successfully
- 400: Validation error

#### `GET /passbooks` 🔒 - Get All Passbooks created by logged in user

**Responses**
- 200: Passbooks fetched successfully
```json
{
    "status": "success",
    "code": 200,
    "message": "Passbooks fetched successfully",
    "data": [
        {
            "nickname": "My Savings",
            "bank_name": "HDFC Bank",
            "account_number": "123456789012345",
            "total_balance": 1500.00
        },
        {
            "nickname": "My Salary",
            "bank_name": "ICICI Bank",
            "account_number": "123456789012345",
            "total_balance": 2000.00
        }
    ],
    "meta": {}
}
```
#### `GET /passbooks/:passbook_id` 🔒 - Get Passbook
#### `DELETE /passbooks/:passbook_id` 🔒 - Delete Passbook
#### `PATCH /passbooks/:passbook_id` 🔒 - Update Passbook
#### `GET /passbooks/:passbook_id/transactions` 🔒 - Get All Transactions paginated
- Query params
    - page: 1
    - limit: 10
    - party_name: "Aditya Gupta"
    - tags: "fun,dividend"
    - type: "CREDIT"

**Responses**
- 200: Transactions fetched successfully
```json
{
    "status": "success",
    "code": 200,
    "message": "Transactions fetched successfully",
    "data": [
        {
            "amount": 1500.00,
            "transaction_date": "2023-12-31T14:48:00.000Z",
            "type": "CREDIT",
            "party_name": "Aditya Gupta",
            "description": "ice cream contribution",
            "tags": "vacation,food,fun"
        },
        {
            "amount": 2000.00,
            "transaction_date": "2023-12-31T15:48:00.000Z",
            "type": "DEBIT",
            "party_name": "TCS",
            "description": "dividend from stocks",
            "tags": "stocks,dividend"
        }
    ],
    "meta": {
        "total_pages": 100,
        "page": 1,
        "limit": 10
    }
}
```

#### `POST /passbooks/:passbook_id/transactions` 🔒 - Add Transaction

**Request**

```json
{
    "amount": 1500.00,
    "transaction_date": "2023-12-31T14:48:00.000Z",
    "type": "CREDIT",
    "party_name": "Aditya Gupta",
    "description": "ice cream contribution",
    "tags": "vacation,food,fun",
}
```
**Responses**
- 201: Transaction added successfully
- 400: Validation error

#### `GET /passbooks/:passbook_id/transactions/:transaction_id` 🔒 - Get Transaction
#### `PATCH /passbooks/:passbook_id/transactions/:transaction_id` 🔒 - Update Transaction
#### `DELETE /passbooks/:passbook_id/transactions/:transaction_id` 🔒 - Delete Transaction

