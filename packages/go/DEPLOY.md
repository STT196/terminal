# Deployment Guide

## Quick Start

### 1. Build the Docker image
```bash
docker build -t terminal-ssh .
```

### 2. Deploy with Docker Compose
```bash
docker compose up -d
```

## Manual Deployment (VPS)

### Prerequisites
- Ubuntu 22.04+ VPS
- Docker and Docker Compose installed
- Domain pointing to your VPS IP

### Setup Steps

1. **Install dependencies**
```bash
sudo apt update
sudo apt install -y docker.io docker-compose nginx certbot python3-certbot-nginx
sudo systemctl enable docker
sudo systemctl start docker
```

2. **Clone and build**
```bash
git clone <your-repo>
cd terminal/packages/go
docker build -t terminal-ssh .
```

3. **Configure domain**
Edit `nginx.conf` and replace `yourdomain.com` with your actual domain.

4. **Get SSL certificate**
```bash
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

5. **Deploy nginx config**
```bash
sudo cp nginx.conf /etc/nginx/sites-available/terminal
sudo ln -s /etc/nginx/sites-available/terminal /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

6. **Run the container**
```bash
docker run -d \
  -p 2222:2222 \
  -v /host_key.pem:/keys/host_key.pem:ro \
  -e SSH_HOST_KEY=/keys/host_key.pem \
  your-image-name
```

## Access

- **Web**: https://yourdomain.com
- **SSH**: `ssh -p 2222 user@yourdomain.com`

## Update Deployment

```bash
cd packages/go
git pull
docker build -t terminal-ssh .
docker stop terminal-ssh
docker rm terminal-ssh
docker run -d \
  --name terminal-ssh \
  --restart unless-stopped \
  -p 2222:2222 \
  -p 8000:8000 \
  terminal-ssh:latest
```

## Troubleshooting

### Check logs
```bash
# Container logs
docker logs -f terminal-ssh

# Nginx logs
sudo tail -f /var/log/nginx/terminal-error.log
```

### Verify services
```bash
# Check if ports are open
sudo netstat -tlnp | grep -E ':(80|443|2222|8000)'

# Test SSH locally
ssh -p 2222 localhost

# Test HTTP
curl http://localhost:8000
```
