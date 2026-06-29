# Use the official Microsoft Playwright image (Noble/Ubuntu 24.04 based)
# This has Playwright, browsers, and system dependencies pre-installed.
# Note: As of June 2026, this image uses Python 3.12.
FROM mcr.microsoft.com/playwright/python:v1.60.0-noble

# Set work directory
WORKDIR /e2e

# Copy requirements first to leverage Docker cache
COPY requirements.txt .

# Install Python dependencies
# Playwright is already in the system, but we install other test dependencies (pytest, etc.)
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the test files
COPY . .

# Set default environment variables
ENV BASE_URL=http://localhost
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Default command to run tests in headless mode
CMD ["pytest", "--base-url", "http://localhost", "--headless"]
