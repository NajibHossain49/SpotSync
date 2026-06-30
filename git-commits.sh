#!/usr/bin/env bash
# Helper: initialise the repo with meaningful, progressive commits.
# Run this ONCE inside the project folder after copying the files.
# It stages files in logical groups so your history shows real development.
set -e

git init
git add go.mod .gitignore .env.example
git commit -m "chore: initialise Go module and project config"

git add config/ database/
git commit -m "feat: add config loader and database connection with pooling"

git add models/
git commit -m "feat: add GORM models for users, zones and reservations"

git add dto/
git commit -m "feat: add DTOs and standard response wrappers with validation tags"

git add utils/errors.go utils/jwt.go utils/validator.go
git commit -m "feat: add JWT helpers, custom errors and validator"

git add repository/user_repository.go repository/zone_repository.go
git commit -m "feat: add user and zone repositories with dynamic availability"

git add repository/reservation_repository.go
git commit -m "feat: implement concurrency-safe reservation with FOR UPDATE row lock"

git add service/
git commit -m "feat: add auth, zone and reservation service layers (business logic)"

git add middleware/
git commit -m "feat: add JWT auth middleware and admin role guard"

git add handler/
git commit -m "feat: add HTTP handlers with centralized error handling"

git add main.go .air.toml Dockerfile
git commit -m "feat: wire dependency injection, routes and Docker build"

git add README.md api-tests.http
git commit -m "docs: add project README and API test collection"

echo "Done. 12 commits created. Now: git remote add origin <your-repo-url> && git push -u origin main"
