[![CircleCI](https://circleci.com/bb/q231950/sputnik/tree/master.svg?style=svg)](https://circleci.com/bb/q231950/sputnik/tree/master) ![Go Report Card](https://goreportcard.com/badge/github.com/q231950/sputnik)

# 🛰 спутник

## Talk to  CloudKit. Server-to-server in Go.

> **Sputnik** enables you to connect to CloudKit from within your Golang package using the Server-to-server communication that CloudKit provides.

### Signing Requests

Sputnik manages the most cumbersome part of CloudKit communication - the signing of your requests. For more information on signing have a look in the [Managing the Signing Identity](https://github.com/q231950/sputnik/wiki/Managing-the-Signing-Identity) section of the [Wiki](https://github.com/q231950/sputnik/wiki).

### Usage

You can use Sputnik either from [the command line](https://github.com/q231950/sputnik/wiki/Sending-Requests#the-sputnik-binary) or [as a package](https://github.com/q231950/sputnik/wiki/Sending-Requests#the-sputnik-package). For more information about requests have a look in the [Sending Requests](https://github.com/q231950/sputnik/wiki/Sending-Requests) section of the [Wiki](https://github.com/q231950/sputnik/wiki).

[Baikonur](https://github.com/q231950/baikonur) uses Sputnik to insert cities into CloudKit:

```go
keyManager := keymanager.New()

config := requesthandling.RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.bish"}
requestManager := requesthandling.New(config, &keyManager, "public")

request, error := requestManager.PostRequest("records/modify", json)
client := &http.Client{}
response, error := client.Do(request)
```

## State

Even though this library works fine for [Baikonur](https://github.com/q231950/baikonur), please keep in mind that it's a [0.0.2](https://github.com/q231950/sputnik/releases).


![Gemeinfrei, <a href="https://commons.wikimedia.org/w/index.php?curid=229349">Link</a>](resources/331px-Sputnik-stamp-ussr.jpg)
