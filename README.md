# Over-The-Air (OTA) Server

Simple OTA server

### Environments

```.dotenv
DB_USERNAME=db_username
DB_PASSWORD=db_password
DB_HOST=db_ip_address
DB_PORT=db_port
DB_INSTANCE=db_instance
AUTO_MIGRATE=1 # 1 = auto migrate when startup, on production it should be 0 or remove
ROOT_USER=basic_user # i'm using basic auth 
ROOT_SECRET=basic_secret 
GOOGLE_CREDENTIALS=/path/to/credentials.json # i'm using google cloud storage to storage .ipa or apk
GCS_BUCKET=bucket-name # bucket name 
PORT=8081
HOST=example.com # the application's host
STATIC_PATH=/path/to/spa/build/folder
```

### TODO

- [ ] Apply auth (Google)
- [ ] Download IPA from GCS as a stream
- [x] Build docker
- [x] Integrate GitHub Actions
