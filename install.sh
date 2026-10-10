#!/bin/sh
# Download the latest ogp binary from GitHub Releases.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/yuzneri/HugoOGPImageGenerator/main/install.sh | sh
#
# Environment variables:
#   INSTALL_DIR  Directory to place the "ogp" binary in (default: current directory)
#   VERSION      Release tag to download, e.g. "v19" (default: latest release)
set -eu

REPO="yuzneri/HugoOGPImageGenerator"
INSTALL_DIR="${INSTALL_DIR:-.}"

fail() {
	echo "error: $*" >&2
	exit 1
}

for cmd in curl tar uname mktemp; do
	command -v "$cmd" >/dev/null 2>&1 || fail "$cmd is required"
done

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
linux | darwin) ;;
*) fail "unsupported OS: $os (download manually from https://github.com/$REPO/releases/latest)" ;;
esac

arch=$(uname -m)
case "$arch" in
x86_64 | amd64) arch=amd64 ;;
aarch64 | arm64) arch=arm64 ;;
*) fail "unsupported architecture: $arch (download manually from https://github.com/$REPO/releases/latest)" ;;
esac

# Resolve the tag by following the /releases/latest redirect instead of the
# REST API, which avoids both a JSON parser dependency and API rate limits.
if [ -z "${VERSION:-}" ]; then
	latest_url=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "https://github.com/$REPO/releases/latest")
	VERSION=${latest_url##*/}
	case "$VERSION" in
	v[0-9]*) ;;
	*) fail "could not determine the latest release (resolved to $latest_url)" ;;
	esac
fi

# Both the archive and the binary inside it carry the version in their names.
binary="ogp${VERSION}-${os}-${arch}"
url="https://github.com/$REPO/releases/download/$VERSION/$binary.tar.gz"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $url"
curl -fsSL -o "$tmp/ogp.tar.gz" "$url"
tar -xzf "$tmp/ogp.tar.gz" -C "$tmp" "$binary"

mkdir -p "$INSTALL_DIR"
mv "$tmp/$binary" "$INSTALL_DIR/ogp"
chmod 755 "$INSTALL_DIR/ogp"

echo "Installed ogp $VERSION to $INSTALL_DIR/ogp"
