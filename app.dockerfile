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


### App Stage
FROM node:24-bookworm
# Setting up envs
ENV N_PATH=/app/
ENV N_USER=node
# Install nodemon for production
RUN npm install -g nx
# Create app dir
WORKDIR ${N_PATH}
RUN chmod -R 777 ${N_PATH}
# Adding files to project
COPY . .

COPY --chown=${N_USER}:${N_USER} --from=dps_base_node_cache /src/node_modules ./node_modules
COPY --chown=${N_USER}:${N_USER} --from=dps_base_node_cache /src/package.json ./package.json
COPY --chown=${N_USER}:${N_USER} --from=dps_base_node_cache /src/package-lock.json ./package-lock.json
# Set user
# USER ${N_USER}
# Start the app
CMD [ "nx", "dev", "dashboard" ]

EXPOSE 3000
