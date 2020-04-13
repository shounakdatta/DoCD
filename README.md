<p align="center">
  <img src='./docs/media/DoCD-logo.png?raw=true' width='40%'>
</p>

---

DoCD (pronounced 'Docked') is a containerization alternative to Docker. DoCD uses system relevant package managers for dependency management and has built-in continuous-deployment support with git webhooks.

DoCD is ideal for projects with multiple service dependencies in the same repo
(i.e. React + Flask + MongoDB etc.), developed in machines not powerful enough
to run a dedicated VM.

It is also useful in environments where implementation of continuous-deployment tools (i.e. Jenkins) is not practical and the developer needs a quick
way to automatically deploy code to development/production servers.

## Expected Project Directory Structure

```
my-project/
  logs/
    service-1.log
    service-2.log
  service-1/
    index.js
  service-2/
    hello.py
  DoCD-config.json
```

## Getting Started

1. Download and extract the latest DoCD binaries from the [release page](https://github.com/shounakdatta/DoCD/releases)
2. Download and install the latest version of [Chocolatey](https://chocolatey.org/install)
3. (Optional) Add the location of your extracted binaries to your [User PATH environment variable](https://helpdeskgeek.com/windows-10/add-windows-path-environment-variable/)
4. In an **administrative shell**, navigate to your project directory
5. Generate your `DoCD-config.json` configuration file
   1. Enter `docd init` to launch the interactive shell
   2. Enter `generate` to launch the configuration generator utility
   3. Proceed through the prompts, enter `exit` or `Ctrl-c` to finish
6. Add your service installation and build commands in the newly generated `DoCD-config.json` file (NOTE: For help, see [here](#understanding-the-`DoCD-config.json`-file))

## Understanding the `DoCD-config.json` File

The `DoCD-config.json` file is a hybrid of the NodeJS `package.json` file and Docker's `Dockerfile`. It stores basic project information, such as Project Name, along with, service installation an build instructons.

### Configuration Fields

#### Base Configuration

| Field              | Type                                      | Description                                                |                          Options                           | Required |
| :----------------- | :---------------------------------------- | :--------------------------------------------------------- | :--------------------------------------------------------: | :------: |
| ProjectName        | String                                    | The name of your project                                   |                             -                              |  False   |
| BasePackageManager | String                                    | The software manager used by DoCD to install your services | `choco` (Windows)<br> `brew` (macOS)<br> `apt-get` (Linux) |   True   |
| Services           | Array ([Service](#service-configuration)) | The configurations for each of your project services       |                             -                              |   True   |

#### Service Configuration

| Field                | Type                                      | Description                                                        |                             Options                             | Required |
| :------------------- | :---------------------------------------- | :----------------------------------------------------------------- | :-------------------------------------------------------------: | :------: |
| ServiceName          | String                                    | The name of the service package as defined by Chocolatey           | Search for your package [here](https://chocolatey.org/packages) |   True   |
| PackageManager       | String                                    | Service specific package managers                                  |                     (i.e. `pip` for python)                     |  False   |
| Path                 | String                                    | Path to the directory where you store the code for your service    |  (i.e. [`./service-1`](#expected-project-directory-structure))  |  False   |
| LogFilePath          | String                                    | Path to the file where you want to store the logs for your service |                  (i.e. `./logs/service-1.log`)                  |   True   |
| InstallationCommands | Array ([Command](#command-configuration)) | Commands to install and configure the service                      |                                -                                |   True   |
| BuildCommands        | Array ([Command](#command-configuration)) | Commands to compile and launch the service                         |                                -                                |   True   |

#### Command Configuration

| Field       | Type           | Description                                                     |                                 Options                                  | Required |
| :---------- | :------------- | :-------------------------------------------------------------- | :----------------------------------------------------------------------: | :------: |
| Directory   | String         | The relative directory in which you want to execute the command | (i.e. `\\service-1`) - NOTE: Can be a blank string for project directory |   True   |
| Command     | String         | The bash command you want to run                                |                    (i.e. `flask run --host=0.0.0.0`)                     |   True   |
| Environment | Array (String) | Any environment variables required for your command             |                       (i.e. `FLASK_APP=hello.py`)                        |  False   |

### Sample Configuration File

```json
{
  "ProjectName": "simple-flask-react-webapp",
  "BasePackageManager": "choco",
  "Services": [
    {
      "ServiceName": "python",
      "Path": "service-2",
      "PackageManager": "pip",
      "LogFilePath": "./logs/service-2.log",
      "InstallationCommands": [
        {
          "Directory": "",
          "Command": "python -m venv venv"
        },
        {
          "Directory": "\\venv\\Scripts",
          "Command": "cmd.exe /C activate.bat"
        },
        {
          "Directory": "\\service-2",
          "Command": "pip install -e ."
        }
      ],
      "BuildCommands": [
        {
          "Directory": "\\service-2",
          "Command": "flask run",
          "Environment": ["FLASK_APP=hello.py", "FLASK_DEBUG=1"]
        }
      ]
    }
  ]
}
```

## License

### Code

All DoCD source code is licensed under the [MIT License](./LICENSE).

### Logo

Copyright to Shounak Datta, 2020.
