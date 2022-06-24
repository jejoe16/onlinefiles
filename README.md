# client-digital-id

There are two source packages a front end and a backend the front end is built in golang and the front end is built in vue.

To get started you will need 2 API keys, Sendgrid (emails), AWS (media files), once obtained you will add them to the golang project `config.yaml`
When ready you will `go build` and then publish to your server and then point your nginx configs to the package locations and you're live. 

## Internal

- Postgres
- go version go1.16.3
- Bootstrap 5.0
- Fontawsome

## Mobile

- Flutter 2.2.1
  
## External

Current server IP: `34.168.73.50`
When switching servers your DNS must point to the external IP of the VM

The server is running nginx and postgres

A API proxy for the backend and a front end www config. 
Along with Certbot for SSL

## API Doc

<https://wes-kay.github.io/api_id_sure/>
API Package: <https://github.com/swaggo/swag>

## Local machine Dependencies to build and run project

- Postgres
- Golang
- Docker
- Make
- NGINX

packages included:

- fresh
- Sendgrid
- Amazon web service sdk

## Connecting to local database

`docker exec -it digital_id_container psql -U postgres`

To kill docker `wsl --shutdown`

## Connecting to external server
`sudo nano ~/.ssh/authorized_keys`
Add your public key to the server and `ssh keeling_wesley@35.199.161.65`

## Connecting to postgres on external server

`sudo -u postgres psql postgres`

## Creating Database

To install the db on the deployment server:

- log into postgres `sudo -u postgres psql postgres`
- Create database: `create database did_data;`
- switch to: `\c did_data`
- Run SQL: `\i '/usr/bin/idsure/migration/db/init.sql';`

## Build

If you're on windows you will have to run `set GOOS=linux`
To go back to local testing on windows `set GOOS=windows`
Then `go build`

## NGINX settings

You will have to set up the reverse proxy and install postgres on the server and set up the ENV variables

## Connecting to server

Request the Admin to add you to the list of users and `ssh username@35.199.161.65`

## Running on the server
once the build has been uploaded you may have to run:
`chmod a+x golang_asset_engine` For proper file permissions on linux

<https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/>

start `sudo systemctl start didservice`
stop `sudo service didservice stop`

## Email server

Twilio send grid, you must update the DNS records if you spin up a new server, this will be found on the GCP External IP

## Admin

The only acessable way to assign a user (admin) as an Admin is to do it through the database for security, once admin you will be able to access the `/admin` pathing in the application (web only):
`UPDATE account SET role = 3 WHERE id = n`

## System

API Schema gen: `swag init --parseDependency --parseInternal`