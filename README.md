# спутник

## Server to server communication

### Create the EC Key

Talking to CloudKit requires authentication. Luckily, there is a command to create the key for you.

`sputnik eckey create`

This will create a `eckey.pem` and place it in the `~/.sputnik/secrets` folder.

### Add EC key to CloudKit Dashboard

You may print out the key in a CloudKit understandable format. Copy the output and paste it as described in [the reference](https://developer.apple.com/library/content/documentation/DataManagement/Conceptual/CloutKitWebServicesReference/SettingUpWebServices/SettingUpWebServices.html#//apple_ref/doc/uid/TP40015240-CH24-SW8)

`sputnik eckey`

### Ping Shelve

[Authenticate Web Service Requests](https://developer.apple.com/library/content/documentation/DataManagement/Conceptual/CloutKitWebServicesReference/SettingUpWebServices/SettingUpWebServices.html#//apple_ref/doc/uid/TP40015240-CH24-SW9)

`[Current date]:[Request body]:[Web service URL subpath]`
