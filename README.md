# спутник

## Server to server communication

### Create the EC Key

Talking to CloudKit requires authentication. Luckily, there is a command to create the key for you.
  
`sputnik eckey create`

This will create a `eckey.pem` and place it in the `./secrets` folder. As long as you don't touch the `.gitignore`, no secret will be committed.

