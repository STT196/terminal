# Hosting Guide — Terminal Portfolio

## Overview

Your portfolio runs as a Docker container exposing an SSH server. Visitors connect via `ssh your-domain.com` and browse your projects in the terminal.

**Architecture:**
```
Visitor → ssh :22 → Docker container (terminal-ssh :2222 internal)
Admin   → ssh :2222 → Server SSH (for management)
```

---

## Initial Setup (First Time)

### 1. Build the Docker image

```bash
cd packages/go
docker build -t terminal-ssh .
```

### 2. Save and upload to your VPS

```bash
# Save image to tar
docker save -o terminal-ssh.tar terminal-ssh

# Upload to server (replace with your server IP and port)
scp -P 2222 terminal-ssh.tar root@YOUR_SERVER_IP:~/ssh-terminal/
```

### 3. Upload config files

```bash
# Upload docker-compose.yml
scp -P 2222 docker-compose.yml root@YOUR_SERVER_IP:~/ssh-terminal/

# Create data directory and upload projects
ssh -p 2222 root@YOUR_SERVER_IP "mkdir -p ~/ssh-terminal/data"
scp -P 2222 data/projects.json root@YOUR_SERVER_IP:~/ssh-terminal/data/
```

### 4. Deploy on the server

```bash
ssh -p 2222 root@YOUR_SERVER_IP

cd ~/ssh-terminal
docker load -i terminal-ssh.tar
docker compose up -d
```

### 5. Verify it works

```bash
ssh YOUR_SERVER_IP
```

You should see your portfolio. Press `q` to quit.

---

## Adding or Editing Projects

Edit the projects file on your server — no rebuild needed:

```bash
ssh -p 2222 root@YOUR_SERVER_IP
nano ~/ssh-terminal/data/projects.json
```

Example format:
```json
[
  {
    "name": "Project Name",
    "year": "2025",
    "technologies": "Go, Docker, AWS",
    "description": "• Key achievement one.\n• Key achievement two.\n• Tech highlights."
  }
]
```

Then restart:
```bash
docker compose -f ~/ssh-terminal/docker-compose.yml restart
```

---

## Updating the Application Code

When you make code changes, rebuild and redeploy:

```bash
# 1. Rebuild locally
cd packages/go
docker build -t terminal-ssh .

# 2. Save and upload
docker save -o terminal-ssh.tar terminal-ssh
scp -P 2222 terminal-ssh.tar root@YOUR_SERVER_IP:~/ssh-terminal/

# 3. Reload on server
ssh -p 2222 root@YOUR_SERVER_IP "cd ~/ssh-terminal && docker load -i terminal-ssh.tar && docker compose up -d"
```

---

## Port Configuration

| Port | Purpose |
|------|---------|
| **22** → container:2222 | Public portfolio (visitors connect here) |
| **2222** | Server SSH admin access |

To change ports, edit `docker-compose.yml`:

```yaml
ports:
  - "YOUR_PORT:2222"    # External:Internal
```

---

## Troubleshooting

**Container not starting:**
```bash
ssh -p 2222 root@YOUR_SERVER_IP
cd ~/ssh-terminal
docker compose logs
```

**Port 22 already in use:**
Check if your server's SSH is on port 22:
```bash
ss -tlnp | grep :22
```
If yes, move it to another port in `/etc/ssh/sshd_config`:
```bash
sed -i 's/^#Port 22/Port 2222/' /etc/ssh/sshd_config
systemctl restart sshd
```

**Projects not showing:**
Verify the file exists and is valid JSON:
```bash
cat ~/ssh-terminal/data/projects.json | python3 -m json.tool
```

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
