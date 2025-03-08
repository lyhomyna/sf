FROM node:22

WORKDIR /usr/src/sf-fe

COPY package.json ./
RUN npm install

COPY . .

CMD ["npm", "run", "start"]
