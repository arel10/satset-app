-- SatSet Laundry Platform Database Initialization
-- Character set: utf8mb4, Collation: utf8mb4_unicode_ci

CREATE DATABASE IF NOT EXISTS satset_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE satset_db;

-- Set global timezone to UTC
SET GLOBAL time_zone = '+00:00';
SET time_zone = '+00:00';
