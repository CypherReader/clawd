# Proton Pass CLI Guide

## Setup Complete ✅

Proton Pass CLI is installed at: `~/.local/bin/pass-cli`

## Basic Commands

### 1. Login (First Time)
```bash
pass-cli login
```
- This will open a browser for authentication
- Since we're on a server, you may need to use device code flow

### 2. List Vaults
```bash
pass-cli vault list
```

### 3. Store a Secret (Git Token Example)
```bash
pass-cli item create \
  --vault "Personal" \
  --title "GitHub Personal Access Token" \
  --field "token=ghp_xxxxxxxxxxxxxxxxxxxx" \
  --field "url=https://github.com" \
  --field "username=your-username" \
  --note "Git token for repository access"
```

### 4. Retrieve a Secret
```bash
# Get specific field
pass-cli item get "GitHub Personal Access Token" --field token

# Get all fields
pass-cli item get "GitHub Personal Access Token"
```

### 5. List Items in a Vault
```bash
pass-cli item list --vault "Personal"
```

## Using Secrets in Scripts

### Method 1: Direct output
```bash
GIT_TOKEN=$(pass-cli item get "GitHub Personal Access Token" --field token)
```

### Method 2: Environment variable injection
```bash
# Store in .env file
pass-cli item get "GitHub Personal Access Token" --field token > .git-token
```

## Persistence Notes

✅ **Will persist between sessions:**
- Installed binary (`~/.local/bin/pass-cli`)
- Configuration files (`~/.config/proton-pass/`)
- This guide file

🔄 **May need re-authentication:**
- Session tokens expire
- Run `pass-cli login` again if needed

## Security Best Practices

1. **Never commit secrets** to version control
2. **Use `.gitignore`** for any local secret files
3. **Regularly rotate** access tokens
4. **Use different vaults** for different purposes (Personal, Work, etc.)

## Troubleshooting

### "Command not found"
```bash
export PATH=$PATH:~/.local/bin
```

### Authentication issues
```bash
# Check login status
pass-cli whoami

# Re-login
pass-cli logout
pass-cli login
```

### Server environment (no browser)
Use device code flow if available, or manual token setup.