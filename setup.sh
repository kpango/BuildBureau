#!/bin/bash

# BuildBureau Setup Script
# This script helps you set up BuildBureau quickly

set -e

echo "🏢 BuildBureau Setup"
echo "===================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "Please install Go 1.21 or later from https://go.dev/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✓ .env file created"
    echo ""
    echo "⚠️  Please edit .env and add your API keys:"
    echo "   - GOOGLE_API_KEY (required)"
    echo "   - SLACK_BOT_TOKEN (optional)"
    echo ""
    echo "After editing .env, run this script again."
    exit 0
fi

# Check if GOOGLE_API_KEY is set
source .env
if [ -z "$GOOGLE_API_KEY" ] || [ "$GOOGLE_API_KEY" = "your-gemini-api-key-here" ]; then
    echo "❌ Error: GOOGLE_API_KEY is not set in .env"
    echo "Please edit .env and add your Gemini API key"
    echo "Get one from: https://makersuite.google.com/app/apikey"
    exit 1
fi

echo "✓ Environment variables loaded"
echo ""

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies ready"
echo ""

# Build the application
echo "🔨 Building BuildBureau..."
mkdir -p bin
go build -o bin/buildbureau ./cmd/buildbureau
echo "✓ Build complete"
echo ""

# Run tests
echo "🧪 Running tests..."
go test ./... -v
echo "✓ Tests passed"
echo ""

echo "✅ Setup complete!"
echo ""
echo "To start BuildBureau, run:"
echo "  ./bin/buildbureau"
echo ""
echo "Or use make:"
echo "  make run"
echo ""
echo "For more information, see:"
echo "  - README.md"
echo "  - docs/QUICKSTART.md"
echo "  - docs/ARCHITECTURE.md"
