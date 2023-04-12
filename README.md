# Development mode
The goal of this project is to provide a simple docker container, that can be used to watch the source codes of a project and automatically rebuild the project when a change is detected. It is very useful during developing (mainly servers), because you don't have to rebuild the project manually every time you make a change.

# Usage

To use this product, you have to create a Dockerfile in your project directory, similar to this one:
```
FROM taskat/devserver:latest AS devmode
WORKDIR /src/devmode

FROM <base image for your project> AS final
WORKDIR /app
COPY --from=devmode /bin/dev_server /bin/dev_server
COPY <your scripts> /app/scripts
ENV WATCH_FOLDER="/app/dev"
ENV INCLUDE_FILES=""
ENV START_SERVER_SCRIPT=""
ENV KILL_SERVER_SCRIPT=""
ENTRYPOINT ["/bin/dev_server"]
```

The base image of your project should be linux based.

This Dockerfile creates the Dockerimage for your project, copies the devmode executable into it, and sets the environment variables. The environment variables are used to configure the devmode executable. The following environment variables are available:
- WATCH_FOLDER: The folder that should be watched for changes.
- INCLUDE_FILES: A regex that is used to filter the files that should be watched.
- START_SERVER_SCRIPT: The path to the script that is executed to start the server.
- KILL_SERVER_SCRIPT: The path to the script that is executed to kill the server.

Also, the start script should store the server's PID in the file ```/tmp/server.pid```, so that the devmode executable can kill the server when it is restarted.