-- Initial Schema for Helios PaaS
-- Version: 1
-- Description: Sets up the core tables for users, projects, applications, and deployments.

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Teams table to group users
CREATE TABLE "teams" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

-- Users table for authentication and authorization
CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" varchar NOT NULL UNIQUE,
  "password_hash" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

-- Junction table for many-to-many relationship between users and teams
CREATE TABLE "team_members" (
    "team_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar NOT NULL DEFAULT 'member', -- e.g., 'admin', 'member'
    PRIMARY KEY (team_id, user_id),
    CONSTRAINT fk_team FOREIGN KEY(team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Projects to organize applications
CREATE TABLE "projects" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "team_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT fk_team FOREIGN KEY(team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Applications table, the core entity that gets deployed
CREATE TABLE "applications" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "project_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "git_repository" varchar NOT NULL,
  "git_branch" varchar NOT NULL DEFAULT 'main',
  "current_backend" varchar NOT NULL DEFAULT 'docker_compose', -- 'docker_compose' or 'k3s'
  "created_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT fk_project FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Deployments table to track each deployment attempt
CREATE TABLE "deployments" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "application_id" uuid NOT NULL,
  "git_commit_sha" varchar,
  "image_uri" varchar,
  "status" varchar NOT NULL DEFAULT 'pending', -- e.g., pending, building, deploying, succeeded, failed
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT fk_application FOREIGN KEY(application_id) REFERENCES applications(id) ON DELETE CASCADE
);

-- Create indexes for foreign keys to improve query performance
CREATE INDEX ON "team_members" ("team_id");
CREATE INDEX ON "team_members" ("user_id");
CREATE INDEX ON "projects" ("team_id");
CREATE INDEX ON "applications" ("project_id");
CREATE INDEX ON "deployments" ("application_id");
CREATE INDEX ON "deployments" ("status");

-- Add a trigger to update the 'updated_at' column on deployments
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_deployments_updated_at
BEFORE UPDATE ON deployments
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();