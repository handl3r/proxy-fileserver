# Proxy file server

Proxy file server integrates with Google Drive for caching file with high frequency download

# Instruction Build and Run

Requirement: Git, golang(version > 1.12), docker, openssl, windows or linux os

## 1. Generate key pair rsa

* Command

```shell script
openssl genrsa -out private.pem 512
openssl rsa -in private.pem -outform PEM -pubout -out public.pem
```

You can use 1024,... for length of key instead of 512, but can not small then 512 Long key make long token

* Output:
    * private.pem is secret key use for your website to generate token or use for the proxy server generate token by api
    * public.pem is key use for proxy for validating token purpose

## 2. Build

- In your project folder :

```shell script
git clone https://github.com/handl3r/proxy-fileserver
cd proxy-fileserver
```

```shell script
# in windows
go build -o proxy-fileserver.exe cmd/main.go 
# or in linux
go build -o proxy-fileserver cmd/main.go
```

- Fix .env file with your config:
    * PROXY_SERVER_ENV= dev or prod, prod will come with log in json format
    * SHARED_ROOT_FOLDER_LOCAL= parent folder of shared-folder on server
    * SHARED_ROOT_FOLDER= shared-folder name your shared folder on Google Drive, same name with mirror folder on server
    * SHARED_ROOT_FOLDER_ID= your shared folder id on Google Drive, example: 1eQ-6Ftg1G1ewA_-7puGpV_YCjHkpmkVe
    * CACHE_TIME_LOCAL_FILE_SYSTEM= unit minute, time to delete a cached file on server from last download time
    * AUTH_PUBLIC_KEY= path to your public key that generated by openssl, example certificates/public512.pem
    * PRIVATE_KEY_LOCATION=certificates/private512.pem on proxy server
    * EXPIRED_TIME_TOKEN=10M expired for token
    * CYCLE_TIME_CLEANER= unit minute, cycle time cleaner wake up to clean expired files
    * MYSQL_USER= username database
    * MYSQL_PASSWORD= password database
    * MYSQL_PORT=3306 port database mysql
    * MYSQL_HOST= host database
    * MYSQL_DATABASE= name database
    * HTTP_PORT=8080 port http
    * REQUIRED_TOKEN=ON (ON/OFF) default is ON, if OFF, you do not need token to access resources
    * GOOGLE_APPLICATION_CREDENTIALS= path to credential file (json) of service account cloud google
    * CREDENTIAL_GOOGLE_OAUTH2_FILE=certificates/credentials.json google oauth drive credential
    * TOKEN_GOOGLE_OAUTH2_FILE=certificates/token.json path to save token
    * GOOGLE_OAUTH2_ENABLE=ON default is ON // if set to OFF, use must config service account
    * INTERACTIVE_MODE=OFF default is off. Set to ON when use want to interact with terminal to exchange google access token

- .env and binary file must be in the same folder

## 3. Setup

* On Cloud Google:
    * Create service account or create OAuth client ID if you want to use GOOGLE_OAUTH2_ENABLE and get certificate in json file (1)
    * Enable Drive api

* On Google Drive:
    * Files storage in 'shared-folder'
    * Share your 'shared-folder' with service account of google cloud if you use service account
    * Get ID of 'shared-folder' on url

* On server:
    * Create parent folder: example 'temp'
    * Create 'shared-folder' inside 'temp'
    * Put your public key and certificate from (1) on somewhere: example 'certificates/cer.json', '
      certificates/public512.pem'
    * Define your .env file (see .env.example)
  
* Use exg tool if you set INTERACTIVE_MODE=OFF to pre-generate token. [exg-tool](additional-tools/google_token_exchange)

Example structure of tree folder tree:
[example-tree-folder](assets/example-folder-tree.png)

### 4. Run

On workspace:

#### Run service

```shell
./proxy-fileserver
```

#### Access api on port 8080:

* Request example:
  ```text
  http://localhost:8080/avt.jpg?token=dsahdha.dsad.ewuegud
  ```
* Response example:
    * 204 for no file exist
    * 500 for system error
    * 200 for success response

### Note

#### For test purpose: use this script for generate token:

```go
package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

func GetCertificate(location string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
func main() {
	privateKey, err := GetCertificate("private512.pem")
	if err != nil {
		panic(err)
	}
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := jwt.MapClaims{}
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(1 * time.Minute).Unix()
	t.Claims = claims
	token, err := t.SignedString(privateKey)
	if err != nil {
		log.Printf("Error when generate token: %s", err)
	}
	fmt.Println(token)
}

```

### API:

#### API Download file:

* Request:

```shell
curl 'http://localhost:8080/shared-folder/avt.jpg?token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTU0NzE2MTYsImlhdCI6MTYxNTQ3MTAxn0.jyX3RIaENdI6JTZdziN3c86cvpqj2M7hpFZTuCATMqtU8uzbs9tLjev21Gng9xwSikb5nY4BcCQRtx9ie29SwQ' \
  -H 'Upgrade-Insecure-Requests: 1' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36' \
  --compressed
```

* Response:
    * Get file successfully:
        * HTTP Code: 200
    * Invalid token:
        * HTTP Code: 401
    * File not found:
        * HTTP Code: 404
    * System Error:
        * HTTP Code: 500
        * Body: "System error. Please contact admin!"

#### API get token:

* Request:
  ```shell
  curl --location --request POST 'ip:8080/auth'
  ```
* Response:

```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTUyNzE4MDMsImlhdCI6MTYxNTI3MTc0M30.ERHFO74v31F6n1psU94qT5mL4G7WMUbiOYnZsdGeIqmpSuJ1DhZvmRSORkZsYFRJcmCbjMJgr6Ukq0-pBHES3g"
}
```

#### API verify token:

* Request:

```shell
curl --location --request POST 'localhost:8080/verify' \
--header 'Content-Type: application/json' \
--data-raw '{
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTU0NzE2MTYsImlhdCI6MTYxNTQ3MTAxNn0.jyX3RIaENdI6JTZdziN3c86cvpqj2M7hpFZTuCATMqtU8uzbs9tLjev21Gng9xwSikb5nY4BcCQRtx9ie29SwQ"
}'
```

* Response:
    * Valid:
        * HTTP Code: 200
        * Body: null
    * Invalid:
        * HTTP Code: 401
        * Body: null

### FAQ

* If any error log with Google Drive, re-share drive folder with service account and wait some minutes or maybe some
  hours. Google Drive need time to re-index your share-across-domain and longer with old files but immediately on new
  files.

* Delete all records in database and all file in shared-folder before re-run

* More information to config cloud google:
  * https://developers.google.com/drive/api/v3/quickstart/go
  
* Confuse about INTERACTIVE_MODE:
  * When interactive mode set to OFF, you must had access token and refresh token storage in TOKEN_GOOGLE_OAUTH2_FILE
  * Encourage to use google_token_exchange tool to pre-generate token, then use INTERACTIVE_MODE=OFF in proxy server