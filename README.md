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

## System Requirements

As of Release 0.1.0b3:

- Windows 7 64-bit (Windows 10 Recommended)
- Powershell v3

## Getting Started

1. Download and extract the latest DoCD binaries (the zip file) from the [release page](https://github.com/shounakdatta/DoCD/releases)
2. Run the `DoCD Installer 1.0.0.exe` executable as an **administrator**
3. In an **administrative shell**, navigate to your project directory
4. Generate your `DoCD-config.json` configuration file (NOTE: For help, see [here](#understanding-the-`DoCD-config.json`-file))
   1. Enter `docd init` to launch the interactive shell
   2. Enter `generate` to launch the configuration generator utility
   3. Proceed through the prompts, enter `exit` or `Ctrl-c` to finish
5. Build your project
   1. Navigate to your project root directory
   2. Enter `docd build`

## Continuous Deployment

To setup your project to automatically deploy to you dev/prod servers, in your server:

1. Start your project with `docd build` or `docd start`
2. Launch ngrok with `enable-cd`
3. Copy your ngrok public URL
4. Add a new webhook in your Github repository
   1. Navigate to your webhook settings at `<your-github-repo-url>/settings/hooks/new`
   2. Paste `<your-copied-ngrok-url>/github-push-master` in the Payload URL
   3. Set Content Type to `application/json`
   4. Hit `Add Webhook`

NOTE: To disable continuous deployment while your project is running, enter `disable-cd`

With this, any latest change to the active git branch will be immediately deployed to the server.

## Understanding the `DoCD-config.json` File

The `DoCD-config.json` file is a hybrid of the NodeJS `package.json` file and Docker's `Dockerfile`. It stores basic project information, such as Project Name, along with, service installation an build instructions.

### Configuration Fields

#### Base Configuration

| Field              | Type                                      | Description                                                |                              Options                               | Required |
| :----------------- | :---------------------------------------- | :--------------------------------------------------------- | :----------------------------------------------------------------: | :------: |
| ProjectName        | String                                    | The name of your project                                   |                                 -                                  |  False   |
| BasePackageManager | String                                    | The software manager used by DoCD to install your services | `choco` (Windows)<br><br> `brew` (macOS)<br><br> `apt-get` (Linux) |   True   |
| Services           | Array ([Service](#service-configuration)) | The configurations for each of your project services       |                                 -                                  |   True   |

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
