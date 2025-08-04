#!/bin/bash

# OPM Deployment Script
# Usage: ./deploy.sh

set -e  # Exit on error

echo "🚀 Starting OPM deployment..."

# Change to script directory
cd "$(dirname "$0")"

# Pull latest changes
echo "📥 Pulling latest changes from git..."
git pull origin main

# Build the server
echo "🔨 Building Go server..."
cd server
go build -o opm-server main.go

# Restart the service
echo "🔄 Restarting OPM service..."
sudo systemctl restart opm

# Wait a moment for service to start
sleep 2

# Check service status
echo "✅ Checking service status..."
sudo systemctl status opm --no-pager

# Test the health endpoint
echo "🏥 Testing health endpoint..."
if curl -f -s http://localhost:8080/health > /dev/null; then
    echo "✅ Health check passed!"
else
    echo "❌ Health check failed!"
    echo "Checking logs..."
    sudo journalctl -u opm -n 20 --no-pager
fi

echo "🎉 Deployment complete!"