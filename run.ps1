# Real-Time Collaboration Platform - Local Run Script (PowerShell)

Write-Host "🚀 Starting Real-Time Collaboration Platform Backend" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "❌ Docker is not running. Please start Docker and try again." -ForegroundColor Red
    exit 1
}

# Check if docker-compose is available
if (-not (Get-Command docker-compose -ErrorAction SilentlyContinue)) {
    Write-Host "❌ docker-compose is not installed. Please install docker-compose and try again." -ForegroundColor Red
    exit 1
}

Write-Host "📦 Starting services with Docker Compose..." -ForegroundColor Yellow
docker-compose up -d

Write-Host ""
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

Write-Host ""
Write-Host "✅ Services started!" -ForegroundColor Green
Write-Host ""
Write-Host "📡 API Server: http://localhost:8080" -ForegroundColor Cyan
Write-Host "🔍 Health Check: http://localhost:8080/health" -ForegroundColor Cyan
Write-Host ""
Write-Host "📊 Services:" -ForegroundColor Cyan
Write-Host "   - PostgreSQL: localhost:5432"
Write-Host "   - Redis: localhost:6379"
Write-Host "   - Backend API: localhost:8080"
Write-Host ""
Write-Host "📝 To view logs: docker-compose logs -f" -ForegroundColor Yellow
Write-Host "🛑 To stop: docker-compose down" -ForegroundColor Yellow
Write-Host ""

