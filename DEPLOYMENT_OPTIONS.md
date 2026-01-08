# Free Deployment Options for FileOrcha

This document outlines free hosting options for deploying your Go file-sharing application.

## 🚀 Recommended Free Hosting Options

### 1. **Fly.io** ⭐ (Highly Recommended)
**Why:** Excellent Go support, generous free tier, Redis support

**Free Tier:**
- 3 shared-cpu-1x VMs (256MB RAM each)
- 160GB outbound data transfer/month
- Free Redis instance available

**Deployment Steps:**
1. Install Fly CLI: `curl -L https://fly.io/install.sh | sh`
2. Sign up: `fly auth signup`
3. Create `fly.toml`:
```toml
app = "your-app-name"
primary_region = "iad"

[build]
  builder = "paketobuildpacks/builder:base"

[env]
  PORT = "8080"
  REDIS_URL = "redis://your-redis-url"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    port = 80
    handlers = ["http"]
    force_https = true

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]
```

4. Deploy: `fly deploy`
5. Add Redis: `fly redis create` (free tier available)

**Pros:**
- Great Go support
- Built-in Redis
- Global edge network
- Easy scaling

**Cons:**
- Need to keep 3 VMs running for free tier
- CLI-based deployment

---

### 2. **Koyeb** ⭐
**Why:** Simple Docker deployment, free tier, Redis support

**Free Tier:**
- 2 nano instances (0.25 vCPU, 256MB RAM)
- 100GB bandwidth/month
- Free Redis instance

**Deployment Steps:**
1. Create `Dockerfile`:
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/download.html .
COPY --from=builder /app/styles.css .
COPY --from=builder /app/script.js .
EXPOSE 8080
CMD ["./main"]
```

2. Push to GitHub
3. Go to [koyeb.com](https://www.koyeb.com)
4. Connect GitHub repo
5. Select Dockerfile
6. Add environment variables:
   - `PORT=8080`
   - `REDIS_URL=your-redis-url`
7. Deploy

**Pros:**
- Very simple UI
- Free Redis included
- Auto-deploy from GitHub
- Good documentation

**Cons:**
- Limited to 2 instances
- May sleep after inactivity

---

### 3. **Zeabur**
**Why:** Easy deployment, free tier, good for beginners

**Free Tier:**
- $5 credit/month (usually enough for small apps)
- Can use external Redis (Upstash free tier)

**Deployment Steps:**
1. Push code to GitHub
2. Go to [zeabur.com](https://zeabur.com)
3. Click "New Project"
4. Import from GitHub
5. Select your repo
6. Add environment variables
7. Deploy

**Pros:**
- Very user-friendly
- Auto-detects Go projects
- Free credits monthly

**Cons:**
- Credits can run out
- Need external Redis (Upstash free tier works)

---

### 4. **Google Cloud Run** ⭐
**Why:** Pay-per-use, generous free tier, serverless

**Free Tier:**
- 2 million requests/month
- 360,000 GB-seconds compute time
- 180,000 vCPU-seconds

**Deployment Steps:**
1. Install gcloud CLI
2. Create `Dockerfile` (same as Koyeb)
3. Build and deploy:
```bash
gcloud run deploy fileorcha \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars PORT=8080,REDIS_URL=your-redis-url
```

4. Use Upstash Redis (free tier) for Redis

**Pros:**
- Very generous free tier
- Auto-scaling
- Pay only for what you use
- Google infrastructure

**Cons:**
- Cold starts possible
- Need Google Cloud account
- More complex setup

---

### 5. **Cyclic.sh**
**Why:** Serverless, free tier, easy deployment

**Free Tier:**
- Unlimited apps
- 100GB bandwidth/month
- Serverless (no cold starts)

**Deployment Steps:**
1. Push to GitHub
2. Go to [cyclic.sh](https://cyclic.sh)
3. Connect GitHub
4. Select repo
5. Add environment variables
6. Deploy

**Pros:**
- True serverless
- No cold starts
- Simple deployment

**Cons:**
- Need external Redis
- Limited bandwidth

---

### 6. **Heroku** (Limited Free Tier)
**Why:** Classic platform, but free tier is limited

**Note:** Heroku removed free tier in 2022, but offers $5/month hobby tier

**Alternative:** Use Heroku's free credits for new accounts

---

## 🔴 Redis Options (Free)

Since your app needs Redis, here are free Redis providers:

### 1. **Upstash Redis** ⭐ (Recommended)
- Free tier: 10,000 commands/day
- Global replication
- Serverless
- [upstash.com](https://upstash.com)

### 2. **Redis Cloud**
- Free tier: 30MB storage
- [redis.com/cloud](https://redis.com/cloud)

### 3. **Fly.io Redis**
- Free with Fly.io deployment
- Included in Fly.io free tier

### 4. **Koyeb Redis**
- Free with Koyeb deployment
- Included in Koyeb free tier

---

## 📋 Quick Comparison

| Platform | Free Tier | Redis Included | Go Support | Ease of Use |
|----------|-----------|----------------|------------|-------------|
| **Fly.io** | ⭐⭐⭐⭐⭐ | ✅ Yes | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Koyeb** | ⭐⭐⭐⭐ | ✅ Yes | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Zeabur** | ⭐⭐⭐⭐ | ❌ No | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Cloud Run** | ⭐⭐⭐⭐⭐ | ❌ No | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Cyclic** | ⭐⭐⭐ | ❌ No | ⭐⭐⭐ | ⭐⭐⭐⭐ |

---

## 🎯 Recommended Setup

**Best Overall:** Fly.io + Fly Redis
- Everything included
- Great performance
- Easy to scale

**Easiest:** Koyeb + Koyeb Redis
- Simplest UI
- Everything in one place
- Good for beginners

**Most Flexible:** Google Cloud Run + Upstash Redis
- Most generous free tier
- Best for scaling
- More control

---

## 📝 Deployment Checklist

Before deploying:

- [ ] Set `CORS_ORIGINS` environment variable
- [ ] Configure Redis connection string
- [ ] Test health check endpoint (`/health`)
- [ ] Verify file size limits work
- [ ] Test password-protected uploads
- [ ] Check rate limiting
- [ ] Monitor logs

---

## 🐳 Dockerfile (Required for most platforms)

Create this `Dockerfile` in your root directory:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary and static files
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/download.html .
COPY --from=builder /app/styles.css .
COPY --from=builder /app/script.js .

EXPOSE 8080

CMD ["./main"]
```

---

## 🔧 Environment Variables

Set these in your hosting platform:

```bash
PORT=8080
REDIS_URL=redis://your-redis-url:6379
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

---

## 📚 Additional Resources

- [Fly.io Go Documentation](https://fly.io/docs/languages-and-frameworks/go/)
- [Koyeb Documentation](https://www.koyeb.com/docs)
- [Google Cloud Run Docs](https://cloud.google.com/run/docs)
- [Upstash Redis](https://upstash.com/docs/redis)

---

## 💡 Tips

1. **Start with Fly.io or Koyeb** - They're the easiest and include Redis
2. **Use Upstash Redis** if your platform doesn't include Redis
3. **Monitor your usage** - Free tiers have limits
4. **Set up alerts** - Know when you're approaching limits
5. **Use CDN** - For static files (styles.css, script.js) if needed

---

**Need help?** Check each platform's documentation or community forums.

