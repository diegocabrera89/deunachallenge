# ms-payment-core

Microservice code in GO. This microservice manages processor transactions

## Pre-steps
If this is your first time using go and GoLand as an IDE you should check the following tutorial:
*   [Installing and configuring Go](https://go.dev/doc/install)

If this is your first time using go and VS Code as an IDE you should check the following tutorial:
*   [Visual Studio Code extension for Go](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code)

Install AWS CLI
*   [Install or update the latest version of the AWS CLI](https://docs.aws.amazon.com/es_es/cli/latest/userguide/getting-started-install.html)

Configure the AWS CLI and profile with name --> dev
*   [Configure the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)

---

## Clone a repository

If you use HTTP
```shel 
git clone [url repository]
```

if you use SSH, configure and execute bellow command
```shell
git clone git@github.com:diegocabrera89/ms-payment-core.git
```
---

## Execute commands

This commands will download our dependencies

```shel
$ npm i
$ go mod download
$ go mod tidy
$ go mod vendor
```

---

## Build project
To build the project run the following command
```shel
    make build
```
---

## Deploy

To deploy the project run the following command
```shel
    make deploy
```
