FROM node:24-alpine

WORKDIR /app

# Install dependencies separately to leverage Docker cache
COPY package.json package-lock.json* ./
RUN npm install

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Source code will be mounted via volume for development
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["npm", "run", "dev"]
