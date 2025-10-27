FROM mirror.gcr.io/oven/bun:1.3.0

WORKDIR /app
RUN mkdir packages
COPY package.json ./
COPY packages/*/package.json ./packages/
COPY packages ./packages
RUN bun install
WORKDIR /app/packages/forge
CMD bun run ./src/forge.ts
