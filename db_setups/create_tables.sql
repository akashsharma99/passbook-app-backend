-- create users table
create table
  passbook_app.users (
    user_id uuid primary key DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash text NOT NULL,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
  );
-- create passbooks table
create table
  passbook_app.passbooks (
    passbook_id uuid primary key DEFAULT gen_random_uuid(),
    user_id uuid references passbook_app.users(user_id) not null,
    bank_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    total_balance DECIMAL(10,2) NOT NULL,
    nickname VARCHAR(255) NOT NULL,
    constraint unique_bank_account unique (bank_name, account_number)
  );
-- create transactions table
create table
  passbook_app.transactions (
    transaction_id UUID primary key DEFAULT gen_random_uuid(),
    amount DECIMAL(10,2) NOT NULL,
    transaction_date timestamp with time zone not null,
    transaction_type VARCHAR(50) NOT NULL,
    party_name VARCHAR(255) not null,
    description text,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    tags VARCHAR(512),
    passbook_id uuid references passbook_app.passbooks(passbook_id) not null,
    user_id uuid references passbook_app.users(user_id) not null
  );