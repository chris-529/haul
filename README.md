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
- Receipt image upload from demo UI
- AI receipt parsing into structured JSON
- PostgreSQL schema completed
- Recipe ingredient join table designed for many-to-many recipe/ingredient relationships
- Auth and login functionality with bcrypt password hashing
- JWT-based authentication with protected route middleware
- Saving and retrieving receipts by authenticated user ID

### In Progress

- Frontend auth and receipt dashboard
- Additional security measures
- Completing remaining REST API handlers

## Roadmap

- Async receipt image upload and AI parsing
- Grocery inventory tracking
- Meal recommendations based on available ingredients
- Cost reduction by replacing Gemini API with OCR + lightweight LLM
- Redis integration for caching and/or background job support
- Deploy app 

## Database Schema

<img src="docs/db-diagram.png" width="400" />
