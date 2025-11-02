#!/bin/bash
set -e

echo "📦 Packaging TailTunnel for macOS..."

# Clean previous builds
rm -rf TailTunnel.app TailTunnel.zip

# Build the app
echo "🔨 Building app..."
make build-macos-app

# Create zip for distribution
echo "📦 Creating ZIP archive..."
zip -r TailTunnel.zip TailTunnel.app

echo "✅ Package created: TailTunnel.zip"
echo ""
echo "To create a GitHub release:"
echo "  1. Tag your release: git tag -a v1.0.0 -m 'Release v1.0.0'"
echo "  2. Push the tag: git push origin v1.0.0"
echo "  3. Upload TailTunnel.zip to the GitHub release"
echo ""
echo "SHA256 checksum:"
shasum -a 256 TailTunnel.zip
