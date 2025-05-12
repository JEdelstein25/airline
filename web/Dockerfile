FROM node:22-alpine

# Set working directory
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@9.15.1 --activate

# Copy package.json and related files
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install

# Copy application code
COPY . .

# Build application if not in development mode
RUN if [ "$NODE_ENV" != "development" ]; then pnpm build; fi

# Expose port that Vite actually uses
EXPOSE 5179

# Start the application in development mode
CMD ["pnpm", "dev", "--host", "0.0.0.0"]