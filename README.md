# Haul

Haul is a full-stack app for turning grocery receipts into structured inventory using AI, with planned meal suggestions based on available ingredients. Deployment on Vercel is planned.

## Stack

- Go
- React
- Vite
- PostgreSQL
- Gemini API
- Redis planned

## Status

### Done

- REST API routes set up
- Receipt image upload from frontend
- AI receipt parsing into structured JSON
- PostgreSQL schema completed
- Recipe ingredient join table designed for many-to-many recipe/ingredient relationships
- Auth and login functionality with bcrypt password hashing
- JWT-based authentication with protected route middleware
- Saving and retrieving receipts by authenticated user ID
- Receipt dashboard with per-user receipt list and selectable receipt detail view
- Receipt deletion for authenticated users

### In Progress

- Frontend auth polish
- Completing remaining REST API handlers
- Receipt editing/update flow
- Additional security measures

## Roadmap

- Receipt editing and update support
- Dockerized background worker for async receipt image parsing
- Redis-backed job queue for receipt processing
- Grocery inventory tracking
- Meal recommendations based on available ingredients
- Cost reduction by replacing Gemini API with OCR + lightweight LLM
- Deploy app

## Database Schema

<img src="docs/db-diagram.png" width="400" />
