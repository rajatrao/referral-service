CREATE TABLE IF NOT EXISTS programs (
    id text PRIMARY KEY,
    name text,
    title text,
    is_active boolean,
    created_at int,
    updated_at int
);

CREATE TABLE IF NOT EXISTS members (
    id text PRIMARY KEY,
    first_name text NOT NULL,
    last_name text,
    email text unique NOT NULL,
    program_id text NOT NULL,
    referral_code text unique NOT NULL,
    is_active boolean,
    created_at int,
    updated_at int,
    CONSTRAINT fk_program FOREIGN KEY (program_id) REFERENCES programs(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS referrals (
    id text PRIMARY KEY,
    first_name text,
    last_name text,
    email text,
    phone text,
    referral_code text NOT NULL,
    status text CHECK (status IN ('pending', 'qualified', 'approved', 'denied')),
    created_at int,
    updated_at int,
    CONSTRAINT fk_member FOREIGN KEY (referral_code) REFERENCES members(referral_code)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);