# Release Guide

This guide explains how to create releases for the k8s-controller project.

## Automated Release Process

The project uses GitHub Actions to automatically build and release binaries and Docker images when you push a version tag.

### Creating a Release

1. **Ensure your code is ready for release**
   - All tests pass
   - Code is merged to main branch
   - CHANGELOG is updated (optional, will be auto-generated)

2. **Create and push a version tag**
   ```bash
   # Create a new tag (use semantic versioning)
   git tag v1.0.0
   
   # Push the tag to trigger the release workflow
   git push origin v1.0.0
   ```

3. **Monitor the release**
   - Go to GitHub Actions tab to monitor the release workflow
   - The workflow will automatically:
     - Build binaries for multiple platforms (Linux, macOS, Windows)
     - Build and push Docker images to GitHub Container Registry
     - Create a GitHub release with all artifacts
     - Generate release notes from git commits

### Release Artifacts

Each release includes:

- **Binaries**: Pre-compiled binaries for multiple platforms
  - `k8s-controller-linux-amd64`
  - `k8s-controller-linux-arm64`
  - `k8s-controller-darwin-amd64` (macOS Intel)
  - `k8s-controller-darwin-arm64` (macOS Apple Silicon)
  - `k8s-controller-windows-amd64.exe`

- **Docker Images**: Multi-platform Docker images
  - `ghcr.io/OWNER/REPO:v1.0.0` (specific version)
  - `ghcr.io/OWNER/REPO:latest` (latest stable release)

- **Kubernetes Manifests**: Ready-to-deploy Kubernetes YAML
  - `k8s-controller.yaml` (includes all necessary resources)

- **Checksums**: SHA256 checksums for all binaries
  - `checksums.txt`

### Version Tagging Convention

Use semantic versioning (SemVer) for tags:

- **Major version** (`v2.0.0`): Breaking changes
- **Minor version** (`v1.1.0`): New features, backward compatible
- **Patch version** (`v1.0.1`): Bug fixes, backward compatible
- **Pre-release** (`v1.0.0-rc1`, `v1.0.0-beta1`): Pre-release versions

### Pre-release vs Stable Release

- **Stable releases**: Tags without suffixes (e.g., `v1.0.0`) are marked as stable
- **Pre-releases**: Tags with suffixes (e.g., `v1.0.0-rc1`) are marked as pre-release

### Manual Release (if needed)

If you need to create a release manually:

1. Go to GitHub Releases page
2. Click "Create a new release"
3. Choose your tag or create a new one
4. Fill in the release title and description
5. Upload any additional files if needed
6. Publish the release

### Troubleshooting

**Release workflow fails:**
- Check GitHub Actions logs for detailed error messages
- Ensure you have proper permissions for GitHub Container Registry
- Verify the tag format follows semantic versioning

**Docker push fails:**
- Ensure `GITHUB_TOKEN` has `packages: write` permission
- Check if the repository is private and container registry is configured correctly

**Binary build fails:**
- Check Go version compatibility
- Ensure all dependencies are properly declared in `go.mod`

### Docker Usage Examples

```bash
# Pull and run the latest release
docker pull ghcr.io/OWNER/REPO:latest
docker run --rm ghcr.io/OWNER/REPO:latest --help

# Run a specific version
docker pull ghcr.io/OWNER/REPO:v1.0.0
docker run --rm ghcr.io/OWNER/REPO:v1.0.0 serve

# Run with custom configuration
docker run -v $(pwd)/config.yaml:/app/k8s-config.yaml \
  ghcr.io/OWNER/REPO:latest serve
```

### Kubernetes Deployment

```bash
# Deploy the latest release
kubectl apply -f https://github.com/OWNER/REPO/releases/latest/download/k8s-controller.yaml

# Deploy a specific version
kubectl apply -f https://github.com/OWNER/REPO/releases/download/v1.0.0/k8s-controller.yaml
```
