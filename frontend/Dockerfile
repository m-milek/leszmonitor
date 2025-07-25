# Stage 1: Install dependencies (for your build process)
FROM node:18-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund --ignore-scripts

# Stage 2: Build the application
FROM node:18-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Stage 3: Production image - Now using 'serve'
FROM node:18-alpine AS runner
WORKDIR /app

# Install the 'serve' package globally in this stage
# This is a lightweight way to serve static files
RUN npm install --ignore-scripts -g serve \
    && addgroup --system --gid 1001 nodejs \
    && adduser --system --uid 1001 appuser

# Copy ONLY the built 'dist' folder
# We don't need node_modules or package.json here if we're just serving static files
COPY --from=builder /app/dist ./dist

# Set permissions
RUN chown -R appuser:nodejs /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 3000

# Start the application using the 'serve' command
# 'serve -s dist' means serve the 'dist' directory and treat it as a Single Page Application
CMD ["serve", "-s", "dist", "-l", "3000"]