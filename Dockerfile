# Node.js 22 LTS on Alpine Linux
FROM node:22-alpine AS runner

WORKDIR /app

# Install OpenSSL & libc6-compat for Prisma engine on Alpine
RUN apk add --no-cache openssl libc6-compat

# Copy server package definitions
COPY server/package*.json ./
COPY server/prisma ./prisma/

# Install dependencies and generate Prisma client
RUN npm ci --include=dev && npx prisma generate

# Copy server application code
COPY server/src ./src

# Create uploads directory for persistent storage
RUN mkdir -p uploads

# Expose backend port
EXPOSE 5000

ENV NODE_ENV=production
ENV PORT=5000

# Start SkillNODE Realtime Conduit API
CMD ["node", "src/server.js"]
