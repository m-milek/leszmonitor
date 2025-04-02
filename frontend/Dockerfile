# Stage 1: Dependencies
FROM node:18-alpine AS deps
WORKDIR /app

# Install dependencies
COPY package.json package-lock.json* ./
RUN npm ci --only=production --no-audit --no-fund

# Stage 2: Builder
FROM node:18-alpine AS builder
WORKDIR /app

# Copy dependencies
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Ensure public directory exists in builder
RUN mkdir -p ./public

# Build the application
ENV NEXT_TELEMETRY_DISABLED 1
RUN npm install --no-audit --no-fund
RUN npm run build
RUN npm prune --production

# Stage 3: Runner
FROM node:18-alpine AS runner
WORKDIR /app

# Set to production environment
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1
ENV HOSTNAME 0.0.0.0

# Create a non-root user
RUN addgroup --system --gid 1001 nodejs && \
 adduser --system --uid 1001 nextjs

# Copy only necessary files
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Handle public directory properly
RUN mkdir -p ./public
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Set the correct permissions
USER nextjs

# Expose the port
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]