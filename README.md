## Required Environment Settings
Create a `.env` file and provide the following settings.
```plaintext
HOST_IP=<local ip>
CERT_CACHE_DIR=certs
DOMAINS=api.site.com,www.site.com
WWW_BACKEND=http://<HOST_IP>:<PORT>
API_BACKEND=http://<HOST_IP>:<PORT>
```