FROM node:20-alpine
WORKDIR /app

RUN rm -rf node_modules/.vite-temp

COPY package.json package-lock.json ./
RUN npm install

# default command runs vite dev server
CMD ["npm", "run", "dev"]
