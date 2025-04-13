# Leszmonitor
Lightweight logging and monitoring

## Setup
### Deployment
Leszmonitor's default deployment method is through Docker using `docker-compose`.
To deploy the application, run the following command in the root directory of the project:
```bash
docker-compose up
```
If you want to run the application in the background, you can use the `-d` (--detach) option.

When run using `docker-compose`, the application will log to the `.logs` directory in the root of the project. This is not configurable as of now, but feel free to just edit the `docker-compose.yml` file to change it.

### Development
For development, it's convenient to use the Taskfile provided in the project. This file contains various commands that can be executed using the `task` command. 

See [taskfile.dev](https://taskfile.dev/) for more information on how to install and use it.

To see a list of available tasks, run:
```bash
task --list
```
