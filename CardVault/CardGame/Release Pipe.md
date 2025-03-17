# Release

# Server Vps
- installed Fail2Ban to prevent brute force attacks  
- ufw is installed 80 , 443 and Open-ssh is open

## notes
- Store the latest `.pck` version number in a metadata file (e.g., `version.json`) on S3.

```json
{
  "version": "1.0.3",
  "pck_url": "https://your-bucket.s3.amazonaws.com/game-1.0.3.pck"
}
```

- archive old versions for ==rollback==
- 