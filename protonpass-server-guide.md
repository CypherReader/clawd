# Proton Pass CLI - Server Environment Guide

## Server Setup Complete ✅

Proton Pass CLI is configured for **headless server environments**.

## Key Configuration

### Filesystem Key Storage
Since we're on a server without a keyring, we use:
```bash
export PROTON_PASS_KEY_PROVIDER=fs
```
- Stores encryption key in `~/.local/share/proton-pass-cli/.session/local.key`
- Required for containers, headless servers, and CI/CD environments
- **Less secure** than OS keyring but works everywhere

### Environment Setup
Load the server environment before using `pass-cli`:
```bash
source ~/.protonpass-env
```

## Login Process for Servers

### Step 1: Prepare Environment
```bash
# Load server configuration
source ~/.protonpass-env

# Force logout if previously used keyring
pass-cli logout --force
```

### Step 2: Login
```bash
pass-cli login
```

### Login Challenges on Servers:
1. **No browser** - Proton Pass CLI typically opens a browser for OAuth
2. **Possible solutions**:
   - **Device code flow** (if supported)
   - **Manual token** (if available)
   - **SSH tunnel** to local machine with browser

## Alternative: Manual Credential Storage

If Proton Pass login isn't feasible on the server, we can use **secure file storage**:

### Option A: Encrypted Environment File
```bash
# Store credentials in encrypted .env file
echo "GITHUB_TOKEN=your_token_here" > .env
chmod 600 .env  # Restrict permissions
```

### Option B: Age Encryption
```bash
# Install age encryption tool
apt-get install -y age

# Encrypt secret
echo "your_token_here" | age -e -r "age1qy..." > secret.age

# Decrypt when needed
age -d secret.age
```

### Option C: GPG Encryption
```bash
# Encrypt with GPG
echo "your_token_here" | gpg --encrypt --recipient "your@email.com" > secret.gpg

# Decrypt
gpg --decrypt secret.gpg
```

## Git Token Storage Examples

### Using Proton Pass (if login works):
```bash
# Store token
pass-cli item create \
  --vault "Personal" \
  --title "GitHub Token" \
  --field "token=ghp_xxxxxxxx" \
  --field "purpose=Git repository access"

# Retrieve token
GIT_TOKEN=$(pass-cli item get "GitHub Token" --field token)
```

### Using Secure Files:
```bash
# Store in encrypted file
echo "ghp_xxxxxxxx" | age -e -r "age1qy..." > .git-token.age

# Use in scripts
GIT_TOKEN=$(age -d .git-token.age)
git clone https://$GIT_TOKEN@github.com/user/repo.git
```

## Persistence Notes

✅ **Will persist:**
- `pass-cli` binary
- Configuration files
- Setup scripts
- This guide

🔄 **May need:**
- Re-authentication (tokens expire)
- Environment reloading

⚠️ **Security considerations:**
- Filesystem key storage is less secure
- Regular token rotation recommended
- Use minimal permissions

## Troubleshooting

### "Error creating client features"
```bash
# Ensure filesystem key provider is set
export PROTON_PASS_KEY_PROVIDER=fs
pass-cli logout --force
pass-cli login
```

### "Command not found"
```bash
export PATH=$PATH:~/.local/bin
```

### Login fails (no browser)
Consider alternative storage methods or check if device code flow is available.

## Next Steps

1. **Try login**: `source ~/.protonpass-env && pass-cli login`
2. **If login fails**: We'll use secure file storage
3. **Store Git token**: Using whichever method works
4. **Test Git operations**: With the retrieved token

What Git repository are we working with? I can help set up the authentication once we have the token stored.