# flash-chat-server
## Simple golang P2P chat server

**Optional:** You can change configurations by adding .env file beside the *`runnable`* file or *`main.go`* file.

.env example:
```
JWT_KEY=your-jwt-encryption-key
JWT_EXP=60 // jwt expiration time in minutes

HTTP_HOST=":8080" // http serving port
```

Download and run your OS compatible latest runnable file from release page

[Release Page](https://github.com/apm-dev/flash-chat-server/releases)

Or you can clone the project and run with go run command

```bash
go mod tidy
go run main.go
```