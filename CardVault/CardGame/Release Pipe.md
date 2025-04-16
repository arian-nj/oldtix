- [x] make docker file ✅ 2025-04-16
- [x] docker compose or stack file ✅ 2025-04-16
- [x] read docker swarm docker ✅ 2025-04-16
- [x] docker secrets ✅ 2025-04-16
- [x] git hub actions ✅ 2025-04-16
- [x] auto create image ✅ 2025-04-16
- [x] setup traefik ✅ 2025-04-16
- [x] on push to release automate deploy ✅ 2025-04-16
- [ ] export pck action
- [ ] upload pack
- [ ] manage versioning

# RELEASE PIP LINE
- [x] learn git branching ✅ 2025-03-17
- [x] run action on change to release branch ✅ 2025-03-17
- [x] wtf is SSH Key ✅ 2025-03-17
- [x] setup vps ✅ 2025-03-17
- [x] compile core go code ✅ 2025-03-18
- [x] copy file to vps ✅ 2025-03-18
- [ ] how to run it automatically
- [ ] cors
docker seems to make more sense people say it's bad 
how dockerize my project
i think the solution in docker compose
(later me: Yep!)

# Smooth Update
start new version
route requests to new version
send kill signal to old version
shut down old version

# VPS NOTES
- installed Fail2Ban to prevent brute force attacks (no longer needed with ssh Key)
- ufw is installed 80 , 443 and Open-ssh is open

# Version
major.minor.patch
major == Launcher needs Update from stores
minor == breaking change server and pck update
patch == pck update
### Reference
https://semver.org/

- Store the latest `.pck` version number in a metadata file (e.g., `version.json`) on S3.

```json
{
  "version": "1.0.3",
  "pck_url": "https://your-bucket.s3.amazonaws.com/game-1.0.3.pck"
}
```

- archive old versions for ==rollback==