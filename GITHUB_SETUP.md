# Push to GitHub - Step by Step Guide

## 🚀 Quick Steps

### 1. Create a GitHub Repository

1. Go to [GitHub.com](https://github.com) and sign in
2. Click the **"+"** icon in the top right → **"New repository"**
3. Repository name: `collab-platform-backend` (or your preferred name)
4. Description: "Real-Time Collaboration Platform Backend - Built in Go"
5. Choose **Public** or **Private**
6. **DO NOT** initialize with README, .gitignore, or license (we already have these)
7. Click **"Create repository"**

### 2. Push Your Code

Run these commands in your project directory:

```bash
# Initialize git (if not already done)
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: Real-Time Collaboration Platform Backend"

# Add your GitHub repository as remote
# Replace YOUR_USERNAME with your GitHub username
git remote add origin https://github.com/YOUR_USERNAME/collab-platform-backend.git

# Push to GitHub
git branch -M main
git push -u origin main
```

### 3. Alternative: Using SSH (if you have SSH keys set up)

```bash
git remote add origin git@github.com:YOUR_USERNAME/collab-platform-backend.git
git push -u origin main
```

## 📝 Detailed Commands

### Check Git Status
```bash
git status
```

### See What Will Be Committed
```bash
git add -n .
```

### Commit with a Message
```bash
git commit -m "Your commit message"
```

### View Remote Repositories
```bash
git remote -v
```

### Push Changes
```bash
git push origin main
```

## 🔐 Authentication

### Option 1: Personal Access Token (Recommended)
1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate new token (classic)
3. Select scopes: `repo` (full control)
4. Copy the token
5. Use it as password when pushing (username = your GitHub username)

### Option 2: GitHub CLI
```bash
# Install GitHub CLI, then:
gh auth login
git push
```

### Option 3: SSH Keys
1. Generate SSH key: `ssh-keygen -t ed25519 -C "your_email@example.com"`
2. Add to GitHub: Settings → SSH and GPG keys → New SSH key
3. Use SSH URL: `git@github.com:USERNAME/REPO.git`

## ✅ Verify Push

After pushing, check your GitHub repository:
- All files should be visible
- README.md should display
- Code should be accessible

## 🔄 Future Updates

To push future changes:

```bash
git add .
git commit -m "Description of changes"
git push
```

## 🎯 Repository Settings (Optional)

After pushing, consider:
- Adding topics/tags: `go`, `websocket`, `real-time`, `collaboration`, `backend`
- Adding description
- Enabling GitHub Actions (if you add CI/CD)
- Adding a license file
- Setting up branch protection rules

