-- Cryo DB Schema v1.1 (Secure Zero-Trust Model)

-- USERS TABLE
CREATE TABLE users (
    id UUID PRIMARY KEY,
    account_number TEXT NOT NULL UNIQUE,                  -- hashed with secret key (HMAC)
    encrypted_recovery TEXT,                       -- optional, CSE backup blob
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- WALLETS TABLE
CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    owner_ref TEXT NOT NULL,                       -- HMAC(account_hash + 'wallet_owner')
    owner_type TEXT NOT NULL CHECK (owner_type IN ('user', 'merchant')),
    currency TEXT NOT NULL,
    encrypted_seed TEXT,                           -- CSE only, optional depending on config
    wallet_type TEXT NOT NULL CHECK (wallet_type IN ('static', 'hot', 'cold', 'ota')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (owner_ref) REFERENCES users(account_number) ON DELETE CASCADE -- conditional on owner_type
);

-- MERCHANTS TABLE
CREATE TABLE merchants (
    id UUID PRIMARY KEY,
    account_ref TEXT NOT NULL,                     -- hashed account_number + 'merchant_owner'
    merchant_name TEXT NOT NULL,
    metadata TEXT,                                 -- encrypted metadata (optional)
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- REFUNDS TABLE
CREATE TABLE refunds (
    id UUID PRIMARY KEY,
    txn_hash TEXT NOT NULL,                        -- links to original transaction
    merchant_hash TEXT NOT NULL,                   -- HMAC(account_hash + 'refund_merchant')
    amount NUMERIC(36, 18) NOT NULL,
    reason TEXT,
    rfnd_status TEXT NOT NULL CHECK (rfnd_status IN ('requested', 'reviewing', 'approved', 'sent', 'rejected')),
    status_detail TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- TRANSACTIONS TABLE
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    owner_hash TEXT NOT NULL,                      -- sender hashed with HMAC(account hash)
    destination_encrypted TEXT NOT NULL,
    destination_hash TEXT NOT NULL,
    txn_type TEXT NOT NULL CHECK (txn_type IN ('user', 'merchant', 'refund')),
    refund_id_ref UUID,                            -- optional FK to refunds
    currency TEXT NOT NULL,
    amount NUMERIC(36, 18) NOT NULL,
    txn_status TEXT NOT NULL CHECK (txn_status IN ('invoice', 'pending', 'confirmed', 'failed')),
    external_ref TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    expiration DATE,                               -- for refunds / invoice expiry

    FOREIGN KEY (refund_id_ref) REFERENCES refunds(id) ON DELETE SET NULL
);

-- -- USER TAGS TABLE (Work in progress)
-- CREATE TABLE user_tags (
--     id UUID PRIMARY KEY,
--     user_id_ref UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     user_tag_hash TEXT NOT NULL,                   -- hash(global_key + user_tag)
--     user_tag_blind TEXT NOT NULL,                  -- optional blind signature or alternate hash
--     status TEXT DEFAULT 'active',
--     created_at TIMESTAMP DEFAULT NOW(),
--     updated_at TIMESTAMP,
--     tag_type TEXT
-- );
