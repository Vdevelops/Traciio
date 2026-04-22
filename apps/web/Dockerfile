# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Enable corepack untuk support packageManager field
RUN corepack enable && corepack prepare pnpm@9.0.0 --activate

# Copy root package files untuk monorepo setup
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages ./packages

# Copy web app package files
COPY apps/web/package.json ./apps/web/
COPY apps/web/next.config.ts ./apps/web/
COPY apps/web/tsconfig.json ./apps/web/
COPY apps/web/components.json ./apps/web/
COPY apps/web/postcss.config.mjs ./apps/web/
COPY apps/web/eslint.config.mjs ./apps/web/

# Copy web app source code
COPY apps/web/app ./apps/web/app
COPY apps/web/src ./apps/web/src
COPY apps/web/public ./apps/web/public

# Install dependencies
RUN pnpm install --frozen-lockfile

# Build application
WORKDIR /app/apps/web
RUN pnpm build

# Production stage
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# Copy standalone output from builder
COPY --from=builder /app/apps/web/.next/standalone ./
COPY --from=builder /app/apps/web/.next/static ./apps/web/.next/static
COPY --from=builder /app/apps/web/public ./apps/web/public

# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Change ownership
RUN chown -R nextjs:nodejs /app

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

WORKDIR /app/apps/web

CMD ["node", "server.js"]

