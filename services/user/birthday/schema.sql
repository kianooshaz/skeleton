-- Birthday Service Database Schema
-- This file contains the SQL schema for the birthday service tables.
-- Run this manually in your PostgreSQL database to create the required tables.
-- Create birthdays table
CREATE TABLE IF NOT EXISTS birthdays (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    date_of_birth DATE NOT NULL,
    age INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    -- Constraints
    CONSTRAINT birthdays_user_id_unique UNIQUE (user_id),
    CONSTRAINT birthdays_age_positive CHECK (age >= 0),
    CONSTRAINT birthdays_age_reasonable CHECK (age <= 200),
    CONSTRAINT birthdays_date_not_future CHECK (date_of_birth <= CURRENT_DATE)
);
-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_birthdays_user_id ON birthdays (user_id);
CREATE INDEX IF NOT EXISTS idx_birthdays_age ON birthdays (age);
CREATE INDEX IF NOT EXISTS idx_birthdays_birth_month ON birthdays (
    EXTRACT(
        MONTH
        FROM date_of_birth
    )
);
CREATE INDEX IF NOT EXISTS idx_birthdays_created_at ON birthdays (created_at);
-- Add comments for documentation
COMMENT ON TABLE birthdays IS 'Stores user birthday information and calculated age';
COMMENT ON COLUMN birthdays.id IS 'Unique identifier for the birthday record';
COMMENT ON COLUMN birthdays.user_id IS 'Foreign key reference to the user (unique per user)';
COMMENT ON COLUMN birthdays.date_of_birth IS 'User date of birth';
COMMENT ON COLUMN birthdays.age IS 'Calculated age based on date of birth';
COMMENT ON COLUMN birthdays.created_at IS 'Timestamp when the record was created';
COMMENT ON COLUMN birthdays.updated_at IS 'Timestamp when the record was last updated';