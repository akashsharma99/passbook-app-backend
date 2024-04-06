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
## Auth Endpoints


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

#### `POST /auth/refresh` - Refresh access token

Expected to be called when the access token is expired. The refresh should be present in the http-only cookie of the request headers.

**Responses**
- 400: Bad request (missing refresh token in cookie)
- 401: Unauthorized (invalid refresh token)
- 500: Internal failures
- 200: Access token refreshed successfully
```json
{
    "status": "success",
    "message": "Token refreshed successfully",
    "data": {
        "access_token": "jwt_token"
    }
}
```
Also as part of token rotation strategy, the server will send a new refresh token in the http-only cookie.
For more info on how refresh tokens work I would suggest reading [this](https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/) article by Auth0 team.

## User Endpoints

> All endpoints marked with the 🔒 symbol require you to pass the access_token as a Bearer token in Authorization header.

#### `GET /users/me` 🔒 - Get User Profile
```bash
curl -X GET http://api.domain.app/v1/users/me -H "Authorization : Bearer <token>"
```

## Passbook Endpoints

#### `POST /passbooks` 🔒 - Create Passbook

**Request**

```json
{
    "bank_name": "Bank of Zelda",
    "account_number": "123512",
    "total_balance": 1024.45,
    "nickname": "salary"
}
```
**Responses**
- 201: Passbook created successfully
```json
{
    "data": {
        "passbook": {
            "passbook_id": "217c0dc1-cd9a-4562-825c-376b0da8a96e",
            "user_id": "3aaff7dd-91f3-4eab-8b26-b4ddbe68e5a5",
            "bank_name": "Bank of Zelda",
            "account_number": "123512",
            "total_balance": 1024.45,
            "nickname": "salary",
            "created_at": "2024-04-02T00:22:09.134347+05:30",
            "updated_at": "2024-04-02T00:22:09.134347+05:30"
        }
    },
    "status": "success"
}
```
- 400: Validation errors or account already exists
- 500: Internal failures

#### `GET /passbooks` 🔒 - Get All Passbooks created by logged in user

**Responses**
- 200: Passbooks fetched successfully
```json
{
    "data": {
        "passbooks": [
            {
                "passbook_id": "217c0dc1-cd9a-4562-825c-376b0da8a96e",
                "user_id": "3aaff7dd-91f3-4eab-8b26-b4ddbe68e5a5",
                "bank_name": "Bank of Zelda",
                "account_number": "123512",
                "total_balance": 1024.45,
                "nickname": "salary",
                "created_at": "2024-04-02T00:22:09.134347+05:30",
                "updated_at": "2024-04-02T00:22:09.134347+05:30"
            },
            {
                "passbook_id": "2aaff5dd-61f3-4eab-8b26-b4ddbe68e5a5",
                "user_id": "3aaff7dd-91f3-4eab-8b26-b4ddbe68e5a5",
                "bank_name": "Bowser Bank",
                "account_number": "900000123512",
                "total_balance": 1221024.45,
                "nickname": "savings",
                "created_at": "2024-04-02T00:22:09.134347+05:30",
                "updated_at": "2024-04-02T00:22:09.134347+05:30"
            }
        ]
    },
    "status": "success"
}
```
#### `GET /passbooks/:passbook_id` 🔒 - Get Passbook

**Responses**
- 200: Passbook fetched successfully
```json
{
    "data": {
        "passbook": {
            "passbook_id": "217c0dc1-cd9a-4562-825c-376b0da8a96e",
            "user_id": "3aaff7dd-91f3-4eab-8b26-b4ddbe68e5a5",
            "bank_name": "Bank of Zelda",
            "account_number": "123512",
            "total_balance": 1024.45,
            "nickname": "salary",
            "created_at": "2024-04-02T00:22:09.134347+05:30",
            "updated_at": "2024-04-02T00:22:09.134347+05:30"
        }
    },
    "status": "success"
}
```
- 404: Passbook not found
- 500: Internal failures

#### `DELETE /passbooks/:passbook_id` 🔒 - Delete Passbook
**Responses**
- 404: Passbook not found
- 200: Passbook deleted successfully
```json
{
    "status": "success",
    "message": "Passbook deleted successfully"
}
```
#### `PATCH /passbooks/:passbook_id` 🔒 - Update Passbook
**Request**

```json
{
    "passbook_id": "217c0dc1-cd9a-4562-825c-376b0da8a96e",
    "user_id": "3aaff7dd-91f3-4eab-8b26-b4ddbe68e5a5",
    "bank_name": "Bank of Zelda",
    "account_number": "123512",
    "total_balance": 2024.45,
    "nickname": "salary old"
}
```
**Responses**

- 200: Passbook updated successfully
```json
{
    "passbook_id": "217c0dc1-cd9a-4562-825c-376b0da8a96e",
    "user_id": "3aaff7dd-91f3-4eab-8b26-b4ddbe68e5a5",
    "bank_name": "Bank of Zelda",
    "account_number": "123512",
    "total_balance": 2024.45,
    "nickname": "salary old",
    "created_at": "2024-04-02T00:22:09.134347+05:30",
    "updated_at": "2024-05-22T01:02:09.134347+05:30"
}
```
- 404: Passbook not found
- 403: Forbidden if user_id on reqeust body and token user_id do not match
- 500: Internal failures

## Transaction Endpoints

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

