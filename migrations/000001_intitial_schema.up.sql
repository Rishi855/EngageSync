CREATE TYPE user_role AS ENUM ('Admin', 'User');

-- Create TenantRegistry table
CREATE TABLE TenantRegistry (
    TenantID UUID PRIMARY KEY,
    OrgName TEXT NOT NULL,
    SchemaName TEXT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    IsActive BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create GlobalUsers table
CREATE TABLE GlobalUsers (
    GlobalUserID UUID PRIMARY KEY,
    Email TEXT NOT NULL UNIQUE,
    Password TEXT NOT NULL,
    TenantID UUID NOT NULL,
    Role user_role NOT NULL,
    FOREIGN KEY (TenantID) REFERENCES TenantRegistry(TenantID)
);

-- Create SuperAdmin table
CREATE TABLE SuperAdmin (
    SuperAdminID UUID PRIMARY KEY,
    Name TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL
);