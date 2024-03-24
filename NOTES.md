## Entities

### User
- user_id ( string primary key)
- username
- email
- password_hash ( store hashed password including salt )
- created_at timestampz
- updated_at timestampz


### Passbook
- passbook_id
- user_id (foreign key to users table)
- bank_name
- account_number ( bank_name + account_number should be unique )
- total_balance
- nickname

### Transaction
- transaction_id
- amount
- transaction_date
- transaction_type (CREDIT/DEBIT)
- party_name
- description
- created_at
- updated_at
- tags ( comma separated list of tags. max 3 tags per transaction)
- passbook_id (foreign key to passbooks table)
- user_id (foreign key to users table)


## Requirements

### Users
- User can register with username, email and password.
- User can login with email and password.

### Passbooks
- User can create multiple passbooks with bank name, account number and zero balance.
- User can view all passbooks.
- user can delete a passbook which also deletes all transactions for that passbook.

### Transactions
- User can add transactions to a passbook with amount, transaction_date, transaction_type, party name, description and tags.
- User can edit a transaction.
- User can delete a transaction.
- User can view all transactions for a passbook.
- User can filter transactions by
    - party name
    - tags
    - transaction_type
