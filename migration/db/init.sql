CREATE EXTENSION pgcrypto;
CREATE TABLE account_status_lookup(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
INSERT INTO account_status_lookup(name) VALUES('Active');
INSERT INTO account_status_lookup(name) VALUES('Banned');
CREATE TABLE account_role_lookup(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
INSERT INTO account_role_lookup(name) VALUES('Unverified');
INSERT INTO account_role_lookup(name) VALUES('User');
INSERT INTO account_role_lookup(name) VALUES('Certifier');
INSERT INTO account_role_lookup(name) VALUES('Admin');
CREATE TABLE account(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT now(),
    failed_login_attempts SMALLINT DEFAULT 0,
    email TEXT NOT NULL UNIQUE,
    -- TODO: Salt and hash
    password TEXT DEFAULT NULL,
    profile_image TEXT DEFAULT NULL,
    name TEXT DEFAULT NULL,
    address TEXT DEFAULT NULL,
    country TEXT DEFAULT NULL,
    phone TEXT DEFAULT NULL,
    passport_number TEXT DEFAULT NULL,
    date_of_birth TIMESTAMPTZ DEFAULT NULL,
    kin TEXT DEFAULT NULL,
    rank TEXT DEFAULT NULL,
    licence TEXT DEFAULT NULL,
    ship_experience INT DEFAULT NULL,
    experience INT DEFAULT NULL,
    nationality TEXT DEFAULT NULL,
    status INT REFERENCES account_status_lookup(id) DEFAULT 1,
    role INT REFERENCES account_role_lookup(id) DEFAULT 1,
    temp BOOLEAN DEFAULT false, 
    mobile_user BOOLEAN DEFAULT false,
    email_validated BOOLEAN DEFAULT false,
    profile_completed BOOLEAN DEFAULT false,
    profile_validated BOOLEAN DEFAULT false
);
CREATE TABLE account_profile(
    id SERIAL PRIMARY KEY UNIQUE,
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT now(),
    profile_image TEXT DEFAULT NULL,
    passport_image_one TEXT DEFAULT NULL,
    passport_image_two TEXT DEFAULT NULL,
    name TEXT DEFAULT NULL,
    address TEXT DEFAULT NULL,
    country TEXT DEFAULT NULL,
    phone TEXT DEFAULT NULL,
    passport_number TEXT DEFAULT NULL,
    date_of_birth TIMESTAMPTZ DEFAULT NULL,
    kin TEXT DEFAULT NULL,
    rank TEXT DEFAULT NULL,
    licence TEXT DEFAULT NULL,
    ship_experience INT DEFAULT NULL,
    experience INT DEFAULT NULL,
    nationality TEXT DEFAULT NULL,
    fk_account_id INT REFERENCES account(id) NOT NULL,
    completed BOOLEAN DEFAULT false
);
CREATE TABLE email_verification(
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    code UUID NOT NULL,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT now(),
    id SERIAL PRIMARY KEY,
    fk_account_id INT REFERENCES account(id),
    used BOOLEAN DEFAULT false
);
CREATE TABLE session(
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    expire TIMESTAMPTZ DEFAULT now(),
    fk_account_id INT REFERENCES account(id) UNIQUE,
    ip TEXT DEFAULT NULL
);
CREATE TABLE course_category_lookup(
    id SERIAL PRIMARY KEY, 
    name TEXT NOT NULL
);
INSERT INTO course_category_lookup(name) VALUES('STCW training');
INSERT INTO course_category_lookup(name) VALUES('Onboard training');
INSERT INTO course_category_lookup(name) VALUES('Non STCW training');
INSERT INTO course_category_lookup(name) VALUES('Company training');
INSERT INTO course_category_lookup(name) VALUES('Shorebased training');
CREATE TABLE course(
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT now(),
    id SERIAL PRIMARY KEY,
    course_name TEXT NOT NULL, 
    certificate_name TEXT NOT NULL,
    description TEXT DEFAULT NULL,
    additional_description TEXT DEFAULT NULL,
    expire INT NOT NULL,
    -- number INT NOT NULL,
    fk_account_id INT REFERENCES account(id),
    fk_course_category_id INT REFERENCES course_category_lookup(id),
    image_source TEXT DEFAULT NULL
);
CREATE TABLE certificate_type_lookup(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
INSERT INTO certificate_type_lookup(name) VALUES('Active');
INSERT INTO certificate_type_lookup(name) VALUES('Void');
INSERT INTO certificate_type_lookup(name) VALUES('Revoked');
INSERT INTO certificate_type_lookup(name) VALUES('Account Deleted');
CREATE TABLE certificate(
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT NULL,
    id SERIAL PRIMARY KEY,
    fk_course_id INT REFERENCES course(id),
    fk_account_id INT REFERENCES account(id),
    type INT REFERENCES certificate_type_lookup(id) DEFAULT 1,
    issued TIMESTAMPTZ DEFAULT now(),
    activated BOOL DEFAULT false,
    accessed INT DEFAULT 0
    -- Has the user completed the certificate
    -- completed BOOLEAN DEFAULT false,
    -- -- Has the certifier validated and vertified the info
    -- verified BOOLEAN DEFAULT false
);
CREATE TABLE certificate_accessed(
    fk_certificate_id INT REFERENCES certificate(id),
    created TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE terms(
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT now(),
    id SERIAL PRIMARY KEY,
    fk_account_id INT REFERENCES account(id),
    terms TEXT NOT NULL,
    confirmed TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE pdf(
    
);
CREATE TABLE alert(
    uid UUID DEFAULT gen_random_uuid() UNIQUE,
    created TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ DEFAULT now(),
    id SERIAL PRIMARY KEY,
    fk_account_id INT REFERENCES account(id),
    value TEXT NOT NULL,
    seen BOOLEAN DEFAULT false
);