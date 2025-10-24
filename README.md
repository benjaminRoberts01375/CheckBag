# <img src="https://github.com/user-attachments/assets/c915c0fd-37a0-489f-a122-3f158348ef59" alt="CheckBag Logo" style="width: 49%;">

![CheckBag Demo](https://github.com/user-attachments/assets/74f02bba-1fe1-45d9-a63c-a482cc3a1e6a)

## What is CheckBag?

CheckBag is an analytics dashboard designed for homelab owners, providing insights into who's accessing your services and from where.

## Why?

A dashboard like Cloudflare's is helpful, but doesn't clearly tell you what traffic is actually allowed through, nor does it give you much for filtering options with the collected data. CheckBag changes this and aims to be a simple to use dashboard that provides meaningful insights into what your network is doing.

## How?

To collect these insights, CheckBag is a proxy that sits between your reverse proxy (ex. [NGINX Proxy Manager](https://nginxproxymanager.com/)) and your services, and uses [Valkey](https://valkey.io/) to store rolling analytics. These analytics can be queried for and turned into a dashboard accessible through your browser.

# Installation

## Manual Installation

CheckBag is deployed via [Docker](https://docs.docker.com/desktop/setup/install/linux/), and requires a little configuration to get up and going.

### Step 1: Install Docker

[Docker](https://docs.docker.com/desktop/setup/install/linux/) is used to "containerize" CheckBag to ensure all of its assets are accounted for. CheckBag is built for a Linux deployment on a NAS or similar server, which typically run some form of Linux.

### Step 2: Downloading Files

1. Go to [the releases page](https://github.com/benjaminRoberts01375/CheckBag/releases) and find the latest version of CheckBag.
2. Download `docker-compose.yml` and `example.env`.
3. Move the files to a folder that you can find again later, and don't mind sticking around.
4. Rename `example.env` to `.env`. Note: this may make the file disappear, so you may need to show hidden files. On Linux it's usually `ctrl + h` or use `ls -a`, macOS is `cmd + shift + .`, and Windows is `Win + h` to show hidden files.

### Step 3: Configure CheckBag

Open `.env`, and you'll see some options. Most notably you'll need to add a secure password to `CACHE_PASSWORD` since this will be used to secure access to collected data. The remaining options can stay the same if you'd like, or can be updated.

### Step 4: Ready for Launch

1. Open a terminal or command line window at the directory you saved your CheckBag files to.
2. Run `docker compose up -d` (`-d` lets you reuse your terminal if you still want it), and CheckBag will launch. You can access it on the WebUI port specified in the `.env` file.


## Unraid Installation

### Step 1: The Unraid UI
- Open the Unraid UI.
- On the top bar, go to Apps.

### Step 2: Install Valkey
- Search for `Valkey`. There should be a single result.
- Click on Valkey, and install. Note the Valkey Port.

### Step 3: Install CheckBag
- Back in the Apps section, search for `CheckBag`.
- Click on CheckBag, and install.
- Set the Valkey IP and port. The IP can be found by going to the top bar > Docker. Find the Valkey row and use the Container IP.
- You can leave the Valkey Password blank, and the CheckBag Data directory and WebUI and Proxy Port the same unless there are conflicts.
- Hit Apply to install.

# Adding Services

Point each of your endpoints listed in your reverse proxy at CheckBag. You'll need to know the IP address of the system CheckBag's running on to get this working. Here's an example of this with NGINX Proxy Manager:

<div style="display: flex;">
<img src="https://github.com/user-attachments/assets/6893ed75-6ccc-4d82-a679-609a68806938" alt="NGINX Configuration Details Screen" style="width: 49%;">
<img src="https://github.com/user-attachments/assets/f2526657-4bdd-4eff-b75b-9d5e5438a2ad" alt="NGINX Configuration Custom Locations Screen" style="width: 49%;">
</div>

In this example, you'll need to use a Custom Location at the location `/` for your proxy host. The NGINX config can be found below, paste it in the advanced configuration by hitting the gear icon in the top right. Replace `CHECKBAG-IP` and `CHECKBAG-PORT` with your own. This will need to be done for all sub-domains.

```conf
location / {
    proxy_set_header Host $host;
    proxy_pass http://CHECKBAG-IP:CHECKBAG-PORT/api/service/;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_read_timeout 86400s;
    proxy_send_timeout 86400s;
    proxy_connect_timeout 86400s;
}
```

# Compatibility

- CheckBag has been tested with CloudFlare for the domain provider and proxy, which provides headers for some information like country of origin. CheckBag may not be out of the box compatible with other proxy hosts, and may require some additional tuning in your reverse proxy. It's highly recommended to add an issue for such problems.
- The provided Docker Image in the release page is built for Linux x86/ARM.
- If you're using CloudFlare, ensure your domain has Rules > Settings > `Remove "X-Powered-By" header` and `Remove visitor IP headers` disabled, and `Add visitor location headers` enabled.
