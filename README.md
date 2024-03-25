` 23-03-2024 `

# Passbook App API

Plan to create api using golang for the passbook app [frontend](https://github.com/akashsharma99/passbook-app-web)

click [here](NOTES.md) for entity and functional requirement list.

base URL: `http://api.passbook.app/v1/`

## Authentication

```
Authorization: Bearer <token>
```

## Generic Responses

### Success Response

```json
{
    "status": "success",
    "code": 200,
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
    "code": 400,
    "message": "error message",
    "data": null,
    "meta": {}
}
```
## Endpoints

#### `POST /users` - Register User

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
- 400: Validation error
- 409: User already exists

#### `POST /users/login` - Login User

**Request**

```json
{
    "email": "john.doe@email.com",
    "password": "password"
}
```
**Responses**
- 200: User logged in successfully
```json
{
    "status": "success",
    "code": 200,
    "message": "User logged in successfully",
    "data":{
        "access_token": "jwt_token",
        "refresh_token": "refresh_token",
        "user": {
            "user_id": "sdvgsrgregf345f34rt34t5",
            "username": "john_doe",
            "email": "john.doe@email.com"
        }
    },
    "meta": {}
}
```
- 401: Unauthorized
#### `GET /users/me` 🔒 - Get User Profile
#### `PATCH /users/me` 🔒 - Reset password

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

