# Over-The-Air (OTA) Server

OTA Server is a simple server that allows users to install iOS apps (beta) directly on devices running the iOS operating
system.

<p align="center">
<img alt="demo" src="https://i.imgur.com/NP9duGF.png" width="50%" />
</p>

### Introduction

OTA Server has a built-in Notification system, if any build is created, 
it will send a notification to a channel, a group on Telegram. 
This is very beneficial when integrating OTA Serer into CI/CD.

I use GCS (Google Cloud Storage) to store builds, using GCS ensures that builds are always available. 
Every time there is a download request, 
the OTA Server will open a connection directly from GCS to the client, 
instead of downloading from GCS and then sending it to the client. 
This will reduce the download time on the client side.

### Authentication
The OTA Server has a built-in authentication system. It simply uses JWT (Json Web Token) to authenticate the user. 

But one point to note is that the OTA installation method for iOS uses the `itms-services` protocol to download the configuration file (.plist). 

Basically, the `itms-services` protocol cannot understand JWT, so a simpler authentication mechanism is needed. 
The OTA Server will generate 1 code to exchange with customers **(exchange_code)**. 

The Exchange Code will be generated when the user logs into the system, and it has the same effect as the Access Token.
There are **2 APIs** that use Exchange Code authentication: the API to download the configuration file, and the API to download the uploaded files.

### Environments
Below are the environment variables that OTA Server uses.

```.dotenv
# MySQL config variables
DB_USERNAME=db_username
DB_PASSWORD=db_password
DB_HOST=db_ip_address
DB_PORT=db_port
DB_INSTANCE=db_instance
AUTO_MIGRATE=1 # 1 = auto migrate when startup, on production it should be 0 or remove

# redis config variables
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=password
REDIS_DB=0

# storage config variables
GOOGLE_CREDENTIALS=/path/to/credentials.json # i'm using google cloud storage to storage .ipa or apk
GCS_BUCKET=bucket-name # bucket name

# server config variables 
PORT=8081
HOST=example.com # the application's host
STATIC_PATH=/path/to/spa/build/folder
SECRET=your-secret

# notification config variables
TELEGRAM_BOT_TOKEN=your_telegram_bot_token # use @BotFather to create new bot and get token
TELEGRAM_GROUP_ID=your_group_id # https://stackoverflow.com/questions/33858927/how-to-obtain-the-chat-id-of-a-private-telegram-channel
```

### Usage
There are several ways to start the OTA Server

#### Binary
Make sure that the environment variables listed above 
have been exported to the environment.
```shell
go mod download && go build -o app . && ./app
```

#### Docker
To make the execute command shorter I usually prepare a file containing the environment variables inside.

```shell
docker build . -t ota-server && \ 
docker run -p 8080:8080 --env-file .env \ 
  -v /path/to/credentials.json:/credential.json:ro \
  ota-server
```

#### Docker-compose
If you use `docker-compose` then `DB_HOST` should be the name of the database service *(database)*, the same goes for `REDIS_ADDR` *(redis)*.

```yaml
version: "3"
services:
  server:
    build: .
    env_file:
      - .env
    ports:
    - 8080:8080
    volumes:
    - /path/to/credentials.json:/credentials.json
  database:
    image: mysql:5.7
    ports:
      - 3306:3306
    volumes:
    - .data:/var/lib/mysql
    environment:
      MYSQL_ROOT_PASSWORD: root_password 
      MYSQL_DATABASE: ota
      MYSQL_USER: ota
      MYSQL_PASSWORD: user_password
  redis:
    image: redis:6.2-alpine
    ports:
      - 6379:6379
```


### TODO

- [ ] ~~Apply auth (Google)~~
- [x] Download IPA from GCS as a stream
- [x] Build docker
- [x] Integrate GitHub Actions
- [x] Telegram Notification
- [x] Built-in Authentication
- [ ] Resume download
