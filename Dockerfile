################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application. This often uses a different base
# image from the build stage where the necessary files are copied from the build
# stage.
#
# The example below uses the alpine image as the foundation for running the app.
# By specifying the "latest" tag, it will also use whatever happens to be the
# most recent version of that image when you build your Dockerfile. If
# reproducability is important, consider using a versioned tag
# (e.g., alpine:3.17.2) or SHA (e.g., alpine@sha256:c41ab5c992deb4fe7e5da09f67a8804a46bd0592bfdf0b1847dde0e0889d2bff).
FROM alpine:latest AS final

# Install any runtime dependencies that are needed to run your application.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        ffmpeg \
        && \
        update-ca-certificates

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

#Create application dir
RUN mkdir /home/application
RUN mkdir /home/application/files

# Copy the executable from the "build" stage.
COPY ./bin/app /home/application

# Copy Database migrations to the container
COPY ./database /home/application/database

# Assign file to user and give execute flag to file
RUN chown appuser /home/application/app
RUN chmod +x /home/application/app
RUN chown appuser /home/application
RUN chown appuser /home/application/files

USER appuser

# Expose the port that the application listens on.
EXPOSE 8080

WORKDIR /home/application

# What the container should run when it is started.
ENTRYPOINT [ "/home/application/app" ]
