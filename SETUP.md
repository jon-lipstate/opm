# Odin Package Registry Setup Guide

This guide will walk you through setting up the Odin Package Registry for local development.

## Prerequisites

- Go 1.24 or higher
- Node.js 22 or higher  
- PostgreSQL 17 or higher
- pnpm (for the client)

## Database Setup

1. Create a PostgreSQL database:
```bash
createdb opm
```

2. Run the schema migration:
```bash
psql -d opm -f server/schema-mvp.sql
```

## GitHub OAuth Setup

1. Go to [GitHub Settings > Developer settings > OAuth Apps](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the application details:
   - **Application name**: Odin Registry Dev (or your preferred name)
   - **Homepage URL**: `http://localhost:9000`
   - **Authorization callback URL**: `http://localhost:8080/api/auth/github/callback`
4. Click "Register application"
5. Copy the **Client ID**
6. Click "Generate a new client secret" and copy the **Client Secret**

> **Note**: For production, you'll need to create a separate OAuth app with your production URLs.

## Server Setup

1. Navigate to the server directory:
```bash
cd server
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and fill in your configuration:
```env
# Server Configuration
HOST=http://localhost
PORT=8080
ENV=development
FRONTEND_URL=http://localhost:9000

# Database
DATABASE_URL=postgres://username:password@localhost:5432/odin_registry?sslmode=disable

# JWT Secret (generate a random string)
JWT_SECRET=your-super-secret-jwt-key-here

# GitHub OAuth (from the previous step)
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=api/auth/github/callback
```

4. Install dependencies:
```bash
go mod download
```

5. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Client Setup

1. Navigate to the client directory:
```bash
cd client
```

2. Install dependencies:
```bash
pnpm install
```

3. Start the development server:
```bash
pnpm dev
```

The client will start on `http://localhost:9000`

## Verifying Your Setup

1. Open `http://localhost:9000` in your browser
2. Click "Sign In" in the top right
3. You should be redirected to GitHub for authentication
4. After authorizing, you should be redirected back and logged in

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running: `pg_isready`
- Check your connection string in `.env`
- Verify the database exists: `psql -l | grep odin_registry`

### OAuth Redirect Issues
- Make sure the callback URL in GitHub matches exactly: `http://localhost:8080/api/auth/github/callback`
- Check that both `HOST` and `PORT` in `.env` match your setup
- Ensure the frontend URL is correct for CORS

### Port Conflicts
- If port 8080 is taken, change `PORT` in server `.env`
- If port 9000 is taken, the Quasar dev server will automatically find another port

## Discord OAuth

OAuth is not yet implemented for discord; we'll be integrating with the discord bot to enable slash commands directly from discord, so its usage may be slightly different than the all web github version.

- [Discord Developer Portal](https://discord.com/developers/applications)
