## Overview
This is a fairly simple Go + Sveltekit philosophy-based application for the purpose of testing out how to deploy to GKE.

## Application Architecture
Backend
 - Golang with standard HTTP library
 - Slog library for structured logging
 - Service pattern - repository -> service -> controller layer connect to each other with dependency injection
 - A gRPC session service which handles all user sessions, separate from main backend

Data
 - SQLC to generate type-safe queries
 - Postgres database
 - Migrations managed by Go Goose

Frontend
 - Svelte 5 Typescript and Sveltekit 
 - TailwindCSS

## Cloud Architecture
Containers
 - Dockferfile.backend
 - Dockerfile.migrate - the container used by a migration job