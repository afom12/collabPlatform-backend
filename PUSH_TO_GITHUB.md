# 🚀 Push to GitHub - Simple Guide

## Step 1: Create GitHub Repository

1. Go to [github.com](https://github.com) and sign in
2. Click the **"+"** icon → **"New repository"**
3. Name: `collab-platform-backend` (or your choice)
4. Description: "Real-Time Collaboration Platform Backend - Built in Go"
5. Choose **Public** or **Private**
6. **DO NOT** check "Initialize with README" (we already have files)
7. Click **"Create repository"**

## Step 2: Run These Commands

Open PowerShell in your project directory and run:

```powershell
# Make sure you're in the project directory
cd "C:\Users\nunus\OneDrive\Desktop\GO project"

# Initialize git (if not done)
git init

# Add all project files
git add .

# Create first commit
git commit -m "Initial commit: Real-Time Collaboration Platform Backend"

# Add your GitHub repository (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/collab-platform-backend.git

# Set main branch
git branch -M main

# Push to GitHub
git push -u origin main
```

## Step 3: Authentication

When you run `git push`, GitHub will ask for credentials:

**Option A: Personal Access Token (Recommended)**
1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token (classic)
3. Select scope: `repo` (full control)
4. Copy the token
5. When prompted:
   - Username: Your GitHub username
   - Password: Paste the token (not your GitHub password)

**Option B: GitHub CLI**
```powershell
# Install GitHub CLI first, then:
gh auth login
git push
```

## ✅ Verify

After pushing, visit your repository on GitHub:
- `https://github.com/YOUR_USERNAME/collab-platform-backend`
- You should see all your files!

## 🔄 Future Updates

To push changes later:

```powershell
git add .
git commit -m "Description of your changes"
git push
```

## 📝 Quick Copy-Paste Commands

Replace `YOUR_USERNAME` with your actual GitHub username:

```powershell
git init
git add .
git commit -m "Initial commit: Real-Time Collaboration Platform Backend"
git remote add origin https://github.com/YOUR_USERNAME/collab-platform-backend.git
git branch -M main
git push -u origin main
```

## ⚠️ Troubleshooting

**If git init was done in wrong directory:**
```powershell
# Remove incorrect git repo
Remove-Item -Recurse -Force .git

# Re-initialize in correct directory
cd "C:\Users\nunus\OneDrive\Desktop\GO project"
git init
```

**If remote already exists:**
```powershell
git remote remove origin
git remote add origin https://github.com/YOUR_USERNAME/collab-platform-backend.git
```

**If you see "fatal: not a git repository":**
```powershell
cd "C:\Users\nunus\OneDrive\Desktop\GO project"
git init
```

