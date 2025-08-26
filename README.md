# <img src="https://github.com/user-attachments/assets/c915c0fd-37a0-489f-a122-3f158348ef59" alt="CheckBag Logo" style="width: 49%;">

![CheckBag Demo](https://github.com/user-attachments/assets/d92c5049-40ef-474f-ae05-7f73a6c6ef6a)

## What is CheckBag?

CheckBag is an analytics dashboard designed for homelab owners, providing insights into who's accessing your services and from where.

## Why?

A dashboard like Cloudflare's is helpful, but doesn't clearly tell you what traffic is actually allowed through, nor does it give you much for filtering options with the collected data. CheckBag changes this and aims to be a simple to use dashboard that provides meaningful insights into what your network is doing.

## How?

To collect these insights, CheckBag is a proxy that sits between your reverse proxy (ex. [NGINX Proxy Manager](https://nginxproxymanager.com/)) and your services, and uses [Valkey](https://valkey.io/) to store rolling analytics. These analytics can be queried for and turned into a dashboard accessible through your browser.

# Installation

CheckBag is deployed via [Docker](https://www.docker.com/), and requires little configuration to get up and going.

### Step 1: Download the repository

1. Open your terminal
2. Ensure you have git and docker installed
3. Run `git clone https://github.com/benjaminRoberts01375/CheckBag`

### Step 2: Setup your secrets

1. Inside the CheckBag folder you just made, create a folder called `.secrets` (including the leading period).
2. Create two files: `valkey.json` and `valkey.conf`.
3. In `valkey.json`, use this format and replace the password with your own (note: the quotes should remain):

```json
{
	"password": "YOUR-PASSWORD-HERE"
}
```

4. In `valkey.conf`, use this format and replace the password with the same password from `valkey.json` (note the _lack_ of quotes here):

```dotfile
requirepass YOUR-PASSWORD-HERE
port 6379
protected-mode yes
tcp-keepalive 30
maxmemory 4GB
```

### Step 3: Running CheckBag

1. Now that CheckBag is configured, you can run `docker compose up -d --build` to get CheckBag running.
2. Point your reverse proxy at CheckBag. You'll need to know the IP address of the system CheckBag's running on to get this working. Here's an example of this with NGINX Proxy Manager:

<div style="display: flex;">
<img src="https://github.com/user-attachments/assets/47792884-9067-467f-9369-ef1624b56e27" alt="NGINX Configuration Details Screen" style="width: 49%;">
<img src="https://github.com/user-attachments/assets/69e81892-7248-4bb6-81d4-87e06caa0a9e" alt="NGINX Configuration Custom Locations Screen" style="width: 49%;">
</div>

In this example, you'll need to use a Custom Location at the location `/` for your proxy host. An example can be found below. Replace `YOUR-IP-HERE` and `YOUR-PORT-HERE` with your own, keeping the rest the same for all proxy hosts.

```conf

location / {
    proxy_set_header Host $host;
    proxy_pass http://YOUR-IP-HERE:YOUR-PORT-HERE/api/service/;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_read_timeout 86400s;
    proxy_send_timeout 86400s;
    proxy_connect_timeout 86400s;
}
```

# Compatibility

CheckBag has been tested with CloudFlare for the domain provider and proxy, which provides headers for some information like country of origin. CheckBag may not be out of the box compatible with other proxy hosts, and may require some additional tuning in your reverse proxy. It's highly recommended to add an issue for such problems.
