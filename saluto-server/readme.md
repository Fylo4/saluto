## Setting Up the Server

just jotting some things I had to configure on the server - not all things - to make it easier when i want to make another server (i.e. dev server) or open up to federation.

Create a Lightsail server. Just bare OS (none of those template), Ubuntu LTS. Mine is 512 MB ram, 2 vCPUs, 20GB SSD.

Make an IPv4 static public IP. This is what I'm using for now to connect the client to the server; I assume at some point in the future I'll set up a DNS and give it a proper address.

Add an IPv4 firewall: Custom / TCP / 8080 / Any IPv4 address (0.0.0.0/0)

Make a Lightsail DB. PostgreSQL. Mine is 1 GB RAM, 2 vCPUs, 40 GB SSD.

On the Networking tab, it can be set to public mode if you want to run the server on localhost for dev purposes, but for security it is best if public mode is disabled.

You'll need to set up .env files. I have .env.local and .env.prd in my cmd file. (Insert directions on how to set it up)

You'll need to have the server private key installed to SSH into it, upload files to it, etc. On the server, connect tab, there's a button to download SSH key. Put that in a safe location (for the sake of future commands, mine is in my ~/Keys/ directory). chmod it to 600.

The server is expecting an .env file for configuration. Upload .env.prd to the server using `scp -i ~/Keys/saluto-key.pem cmd/.env.prod ubuntu@<ip>:/home/ubuntu/.env.prod`

(Actually, you probably don't need .env.prd, do the following instead)

I am now making a systemd service file to provide the env variables, run the app on launch, reboot the app on crash, etc.

`sudo nano /etc/systemd/system/app.service`

paste the following, with your credentials:
```
[Unit]
Description=Go API Server
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/app
Restart=always
Environment=APP_ENV=production
Environment=DB_USER=dbmaster
Environment=DB_PASS=yourpassword
Environment=DB_HOST=ls-xxxx.rds.amazonaws.com
Environment=DB_PORT=5432
Environment=DB_NAME=dbmaster
Environment=DB_SSLMODE=require

[Install]
WantedBy=multi-user.target
```
save and exit
Reload systemd: `sudo systemctl daemon-reload`
Start the service: `sudo systemctl start app`
Enable on boot: `sudo systemctl enable app`
Status check: `sudo systemctl status app`


## Running

cd into the cmd folder and do `APP_ENV=local go run .`

(It actually sets APP_ENV to local by default, so you could just do `go run .`)

To compile the program, go to the cmd folder and run:
`GOOS=linux GOARCH=amd64 go build -o ../build/app`

To upload the built file to the server:
`scp -i ~/Keys/saluto-key.pem ../build/app ubuntu@<ip>:/home/ubuntu/app` (You need to have the key file)

To run on the server:
`APP_ENV=prod ./app`
