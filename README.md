# flash-chat-server
> for teaching purpose !

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
Or build it for your desired environment
```bash
env GOOS=target-OS GOARCH=target-architecture go build main.go
```
You can find target OS and architectures
[Here](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04)

----------
## Usage
There are 4 http methods
1. `POST /register`
```json
// Request
{
    "name":"Parsa",
    "username":"apm",
    "password":"secret"
}
// Response
{
  "code": "200",
  "message": "welcome Parsa",
  "content": {
    "token": "jwt-token"
  }
}
```
2. `POST /login`
```json
// Request
{
    "name":"Parsa",
    "username":"apm",
    "password":"secret"
}
// Response
{
  "code": "200",
  "message": "",
  "content": {
      "token": "jwt-token"
  }
}
```
3. `GET /users` to get list of all users
```json
// add Authorization header like:
// Authorization: Bearer jwt-token
{
  "code": "200",
  "message": "",
  "content": [
    {
      "name": "Parham",
      "username": "popo"
    },
    {
      "name": "Parsa",
      "username": "apm"
    }
  ]
}
```
4. `GET /chats/:username` to start a chat with a user, this request will open a ws connection to send and receive message through websocket
```json
// add Authorization header like:
// Authorization: Bearer jwt-token

// if I want to start chat with "apm", url should be:
// http://<host:port>/chats/apm

// WS send message as raw string
// "hello world!"

// WS received message structure
{
  "from": "apm",
  "body": "hello world!",
  "sent_at": 1630121139  // UTC unix time 
}
```
