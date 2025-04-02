# Stage 1: Dependencies
FROM node:18-alpine AS deps
WORKDIR /app

# Install dependencies with precise control
COPY package.json package-lock.json* ./
# Install production dependencies only
RUN npm ci --only=production --no-audit --no-fund

# Stage 2: Builder
FROM node:18-alpine AS builder
WORKDIR /app

# Copy production dependencies
COPY --from=deps /app/node_modules ./node_modules
# Copy source files
COPY . .

# Install dev dependencies for build process
RUN npm install --no-audit --no-fund

# Build the application
ENV NEXT_TELEMETRY_DISABLED 1
RUN npm run build

# Remove dev dependencies and keep only production ones
RUN npm prune --production

# Stage 3: Runner (ultra minimal)
FROM node:18-alpine AS runner
WORKDIR /app

# Set production environment
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1
ENV HOSTNAME 0.0.0.0

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
  adduser --system --uid 1001 nextjs && \
  # Clean npm cache and apk cache to reduce image size
  npm cache clean --force && \
  rm -rf /var/cache/apk/*

# Copy only the minimal required files
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Set user and expose port
USER nextjs
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]