# URL Screenshot Service

📸 Capture screenshots for all websites corresponding to the received URLs. (Server)

⏲️ Consistently refreshing every captured image. (Scheduler)

🌥️ Utilizing Cloudflare's R2 object storage to store images.

## What is this for

In order to preview external links within my website, I have implemented this. As we are aware, not all websites properly include `og` tags...

## Quick Start

```bash
git clone https://github.com/SkyeYoung/url-screenshot-service.git && cd url-screenshot-service

# create volume dir, config
mkdir log screenshot && cp template.config.json config.json

# edit what U want
vim config.json 

sudo docker compose up -d
```
