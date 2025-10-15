# Deployment Guide for 4me Todos

This guide covers deploying the 4me Todos application to production.

## Overview

The application consists of three main components:

1. **Frontend** - Vue.js application (deploy to Vercel)
2. **Backend** - Go API server (deploy to Railway, Render, or DigitalOcean)
3. **Database** - PostgreSQL (use Supabase or separate hosting)

## Prerequisites

- GitHub account
- Vercel account
- Railway/Render/DigitalOcean account
- Supabase account
- Google Cloud Console project (for OAuth)

## Step 1: Setup Supabase

### Database Setup

1. Create a Supabase project at <https://supabase.com>
2. Note your database connection string from Settings > Database
3. The database will be used by your backend API

### Storage Setup

1. Go to Storage in your Supabase dashboard
2. Create a new bucket named `4me-attachments`
3. Set bucket to **public** or configure Row Level Security policies:

   ```sql
   -- Allow authenticated users to upload
   CREATE POLICY "Allow authenticated uploads"
   ON storage.objects FOR INSERT
   TO authenticated
   WITH CHECK (bucket_id = '4me-attachments');

   -- Allow public read access
   CREATE POLICY "Allow public read"
   ON storage.objects FOR SELECT
   TO public
   USING (bucket_id = '4me-attachments');

   -- Allow users to delete their own files
   CREATE POLICY "Allow authenticated delete"
   ON storage.objects FOR DELETE
   TO authenticated
   USING (bucket_id = '4me-attachments');
   ```

4. Note your Supabase URL and anon key from Settings > API

## Step 2: Setup Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable Google+ API
4. Go to Credentials > Create Credentials > OAuth 2.0 Client ID
5. Configure OAuth consent screen
6. Add authorized redirect URIs:
   - Development: `http://localhost:8080/api/auth/google/callback`
   - Production: `https://your-backend-domain.com/api/auth/google/callback`
7. Note your Client ID and Client Secret

## Step 3: Deploy Backend

### Option A: Railway

1. Push your code to GitHub
2. Go to <https://railway.app>
3. Click "New Project" > "Deploy from GitHub repo"
4. Select your repository
5. Railway will auto-detect the Go application
6. Add environment variables in the Variables tab:

   ```
   DATABASE_URL=your-supabase-postgres-url
   JWT_SECRET=your-random-secret-key-here
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   GOOGLE_REDIRECT_URL=https://your-backend-url.railway.app/api/auth/google/callback
   SUPABASE_URL=https://your-project.supabase.co
   SUPABASE_KEY=your-supabase-anon-key
   FRONTEND_URL=https://your-frontend.vercel.app
   PORT=8080
   ```

7. Set custom start command: `./backend/cmd/api/main.go`
8. Deploy

### Option B: Render

1. Go to <https://render.com>
2. Click "New +" > "Web Service"
3. Connect your GitHub repository
4. Configure:
   - **Name**: 4me-backend
   - **Environment**: Go
   - **Build Command**: `cd backend && go build -o bin/server cmd/api/main.go`
   - **Start Command**: `./backend/bin/server`
5. Add environment variables (same as above)
6. Create Web Service

### Option C: DigitalOcean App Platform

1. Go to DigitalOcean > Apps
2. Create App from GitHub
3. Select repository and branch
4. Configure:
   - **Type**: Web Service
   - **Build Command**: `cd backend && go build -o bin/server cmd/api/main.go`
   - **Run Command**: `./backend/bin/server`
5. Add environment variables
6. Deploy

## Step 4: Deploy Frontend to Vercel

### Automatic Deployment

1. Push your code to GitHub
2. Go to <https://vercel.com>
3. Click "New Project"
4. Import your GitHub repository
5. Configure build settings:
   - **Framework Preset**: Vite
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist`
6. Add environment variables:

   ```
   VITE_API_URL=https://your-backend-url.railway.app/api
   VITE_GOOGLE_CLIENT_ID=your-google-client-id
   ```

7. Deploy

### Using Vercel CLI

```bash
cd frontend
npm install -g vercel
vercel login
vercel
```

Follow the prompts and set environment variables when asked.

## Step 5: Update Google OAuth Redirect URIs

After deployment, update your Google OAuth settings:

1. Go back to Google Cloud Console > Credentials
2. Edit your OAuth 2.0 Client ID
3. Add production redirect URI:
   - `https://your-backend-url.railway.app/api/auth/google/callback`
4. Save

## Step 6: Update CORS Settings (Backend)

Ensure your backend's CORS middleware allows requests from your Vercel frontend:

The `FRONTEND_URL` environment variable should be set to your Vercel deployment URL.

## Step 7: Test the Deployment

1. Visit your Vercel URL
2. Test user registration
3. Test login
4. Test Google OAuth
5. Create a project
6. Create boards and tasks
7. Test drag-and-drop
8. Upload a file
9. Add comments

## Environment Variables Reference

### Backend Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@host:5432/db` |
| `JWT_SECRET` | Secret key for JWT tokens | Random 32+ character string |
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | From Google Console |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | From Google Console |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | `https://api.yourapp.com/api/auth/google/callback` |
| `SUPABASE_URL` | Supabase project URL | `https://xxx.supabase.co` |
| `SUPABASE_KEY` | Supabase anon key | From Supabase dashboard |
| `FRONTEND_URL` | Frontend URL for CORS | `https://yourapp.vercel.app` |
| `PORT` | Server port | `8080` |

### Frontend Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `https://api.yourapp.com/api` |
| `VITE_GOOGLE_CLIENT_ID` | Google OAuth Client ID | From Google Console |

## Monitoring & Maintenance

### Backend Logs

- **Railway**: View logs in the Railway dashboard
- **Render**: Check logs in the Render dashboard
- **DigitalOcean**: Use the DigitalOcean logs viewer

### Database Backups

Supabase provides automatic daily backups. You can also:

1. Go to Supabase > Database > Backups
2. Create manual backups before major changes
3. Download backups for local storage

### Scaling

- **Railway**: Auto-scales based on usage
- **Render**: Configure auto-scaling in settings
- **DigitalOcean**: Adjust instance size as needed

### SSL/HTTPS

All platforms provide automatic SSL certificates:

- Vercel: Automatic for custom domains
- Railway/Render/DO: Automatic HTTPS

## Custom Domain Setup

### Frontend (Vercel)

1. Go to your project settings
2. Click "Domains"
3. Add your custom domain
4. Update DNS records as instructed
5. Wait for DNS propagation

### Backend (Railway/Render)

1. Go to Settings > Domains
2. Add custom domain
3. Update DNS with provided CNAME
4. Update `GOOGLE_REDIRECT_URL` environment variable
5. Update Google OAuth redirect URIs

## Troubleshooting

### Common Issues

1. **CORS errors**: Check `FRONTEND_URL` in backend environment variables
2. **OAuth not working**: Verify redirect URIs in Google Console
3. **Database connection failed**: Check `DATABASE_URL` format
4. **File uploads failing**: Verify Supabase storage bucket is public
5. **API calls failing**: Ensure `VITE_API_URL` doesn't have trailing slash

### Health Checks

Backend health endpoint: `GET /health`

Should return:

```json
{
  "status": "ok"
}
```

## Rollback Procedure

### Vercel

1. Go to Deployments
2. Find the previous working deployment
3. Click "..." > "Promote to Production"

### Railway/Render

1. Go to Deployments
2. Revert to previous successful deployment
3. Or redeploy from a specific Git commit

## Cost Estimates

### Free Tier (Good for personal use)

- **Vercel**: Free for hobby projects
- **Railway**: $5/month credit free
- **Render**: Free tier available
- **Supabase**: Free tier (500MB database, 1GB storage)

### Paid (For production/scaling)

- **Vercel Pro**: $20/month
- **Railway**: Pay as you go (~$10-50/month)
- **Render**: $7/month for starter instance
- **Supabase**: $25/month for Pro tier

## Security Checklist

- [ ] Change JWT_SECRET to a strong random value
- [ ] Use HTTPS for all endpoints
- [ ] Enable Supabase Row Level Security
- [ ] Set up database backups
- [ ] Configure proper CORS origins
- [ ] Use environment variables for all secrets
- [ ] Enable rate limiting (if needed)
- [ ] Set up monitoring and alerts

## Next Steps

1. Monitor application performance
2. Set up error tracking (e.g., Sentry)
3. Add analytics (e.g., Google Analytics, Plausible)
4. Configure domain emails
5. Set up CI/CD pipelines
6. Add automated tests

---

**Need Help?**

Check the logs first, then review this guide. Most issues are related to environment variables or CORS configuration.
