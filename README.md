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
<set you environment variable>
ENTRYPOINT ["/bin/dev_server"]
```

The base image of your project should be linux based.

This Dockerfile creates the Dockerimage for your project, copies the devmode executable into it, and sets some environment variables. The environment variables are used to configure the devmode executable. The following environment variables are supported:
- WATCH_FOLDER: The folder that should be watched for changes. Default value: ```/app/dev```
- INCLUDE_FILES: A regex that is used to filter the files that should be watched. Default value: ```.*```
- START_SERVER_SCRIPT: The path to the script that is executed to start the server. Default value: ```/app/scripts/start.sh```
- KILL_SERVER_SCRIPT: The path to the script that is executed to kill the server. Default value: ```/app/scripts/kill.sh```
- PID_FILE: The path to the file that contains the PID of the server. Default value: ```/tmp/server.pid```
- WAIT_FOR_SERVER_KILL: The wait timeout while waiting for server shutdown. It accepts values like ```10s```, ```100ms```. Default value: ```100ms```
- TIMEOUT_BETWEEN_CHECKS: The wait time between cecking the source files for changes. It accepts values like ```10s```, ```100ms```. Default value: ```500ms```

Also, the start script should store the server's PID in the file PID_FELI, so that the devmode executable can kill the server when it is restarted.