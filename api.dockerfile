### Node Cache Stage
FROM node:24-alpine AS dps_base_node_cache
# Configure path
WORKDIR /src/
RUN chmod -R 777 /src/
# Install app dependecies
USER node
COPY ./package.json .
# Install dependencies
RUN npm install


### Go
FROM golang:1.26-bookworm

ENV GO111MODULE=on
ENV G_PATH=/app
ENV G_USER=api

# Install dependencies
RUN apt update && apt install -yq \
    git make gcc libc6-dev ca-certificates sudo curl \
    && apt clean \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://deb.nodesource.com/setup_24.x | sudo -E bash - \
    && apt update \
    && apt install -y nodejs
RUN npm install -g npm@latest

# Create app dir
WORKDIR ${G_PATH}
COPY . .

COPY --chown=${G_USER}:${G_USER} --from=dps_base_node_cache /src/node_modules ./node_modules
COPY --chown=${G_USER}:${G_USER} --from=dps_base_node_cache /src/package.json ./package.json
COPY --chown=${G_USER}:${G_USER} --from=dps_base_node_cache /src/package-lock.json ./package-lock.json

# Install nx
RUN npm install -g nx
# Install codegangsta/gin
RUN go install github.com/codegangsta/gin@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN nx tidy api

CMD [ "nx", "dev", "api" ]

EXPOSE 8080
