#!/usr/bin/env bash
# Publish GoReleaser's CLI cask with Homebrew's declarative install steps.
# CASK_DRY_RUN=1 prints the cask without accessing the tap or requiring a token.
set -euo pipefail

VERSION="${1:?usage: publish-cask.sh <version>}"
CASK_NAME="devkube"

case "$VERSION" in
  *SNAPSHOT* | *snapshot* | *-dirty)
    echo "publish-cask: snapshot build ($VERSION), skipping tap update"
    exit 0
    ;;
esac

CLI_CASK="dist/homebrew/Casks/${CASK_NAME}.rb"
if [ ! -f "$CLI_CASK" ] || [ "$(tail -n 1 "$CLI_CASK")" != "end" ]; then
  echo "publish-cask: missing or incomplete generated cask: $CLI_CASK" >&2
  exit 1
fi
if ! grep -Fqx "  version \"${VERSION}\"" "$CLI_CASK" || grep -Eq '^  postflight(_steps)? do' "$CLI_CASK"; then
  echo "publish-cask: regenerate $CLI_CASK with the current GoReleaser configuration for $VERSION" >&2
  exit 1
fi

WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT
CASK_FILE="$WORKDIR/${CASK_NAME}.rb"
sed '$d' "$CLI_CASK" > "$CASK_FILE"
cat >> "$CASK_FILE" <<EOF

  postflight_steps do
    on_macos do
      run "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "{{staged_path}}/${CASK_NAME}"]
    end
  end
end
EOF

if [ "${CASK_DRY_RUN:-0}" = "1" ]; then
  cat "$CASK_FILE"
  exit 0
fi

: "${GITHUB_TOKEN:?GITHUB_TOKEN is required to push the cask to the tap}"
git clone --depth 1 \
  "https://x-access-token:${GITHUB_TOKEN}@github.com/adrianliechti/homebrew-tap.git" \
  "$WORKDIR/tap" >/dev/null 2>&1
mkdir -p "$WORKDIR/tap/Casks"
cp "$CASK_FILE" "$WORKDIR/tap/Casks/${CASK_NAME}.rb"
cd "$WORKDIR/tap"
git add "Casks/${CASK_NAME}.rb"
if git diff --cached --quiet; then
  echo "publish-cask: cask already up to date, nothing to push"
  exit 0
fi
git -c user.name="Adrian Liechti" -c user.email="adrian@localhost" \
  commit -m "${CASK_NAME} ${VERSION}" >/dev/null
git push origin HEAD >/dev/null 2>&1
echo "publish-cask: pushed ${CASK_NAME} ${VERSION} to adrianliechti/homebrew-tap"
