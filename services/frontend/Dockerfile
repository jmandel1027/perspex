FROM node:19.1-buster AS tilt

ENV PACKAGE_PATH="."
ENV WORKING_DIR="/mnt/app"

WORKDIR ${WORKING_DIR}

COPY ${PACKAGE_PATH}/package.json ${WORKING_DIR}/package.json
COPY ${PACKAGE_PATH}/package-lock.json ${WORKING_DIR}/package-lock.json

RUN npm install

COPY ${PACKAGE_PATH} ${WORKING_DIR}

CMD [ "npm", "run", "dev" ]

FROM node:19.1-buster AS common

ENV PACKAGE_PATH="."
ENV WORKING_DIR="/mnt"
ENV RELEASE="not-set"

WORKDIR ${WORKING_DIR}

COPY ${PACKAGE_PATH}/package.json ${WORKING_DIR}/package.json
COPY ${PACKAGE_PATH}/package-lock.json ${WORKING_DIR}/package-lock.json

RUN npm install

COPY ${PACKAGE_PATH} ${WORKING_DIR}

# RUN npm run gen

FROM common as builder

RUN npm run build

# --- deployable image
FROM common AS deployable

# Copy static assets from builder stage
COPY --from=builder /mnt/dist .

CMD [ "npm", "start" ]
