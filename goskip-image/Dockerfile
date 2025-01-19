FROM oven/bun:1.1.44

WORKDIR /app

COPY package.json bun.lockb ./

RUN bun install

EXPOSE 8080
EXPOSE 8081

CMD ["bun", "run", "skip.ts"]
