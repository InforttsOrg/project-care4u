#!/usr/bin/env bash
# 🧬 Care4u Release Validation Gatekeeper.
# Verifies the Go microservices compile and their tests pass, the compose stack and
# the CDN publisher contract are intact, and no credentials are tracked — then bumps
# the patch version. Use "--test-only" to validate without bumping (CI-safe).
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"
VERSION_FILE="$PROJECT_DIR/.version"

# Fail closed on unknown/incoherent arguments: only the release path (no args) and
# "--test-only" are valid; a stray second arg used to be silently ignored.
TEST_ONLY=0
case "${1:-}" in
  "") ;;
  "--test-only") TEST_ONLY=1 ;;
  "-h"|"--help")
    echo "Usage: $0 [--test-only]"
    echo "  (no args)     validate and bump the patch version"
    echo "  --test-only   validate only; leaves .version untouched (CI-safe)"
    exit 0
    ;;
  *)
    echo "validate-release.sh: unknown argument '$1' (expected --test-only)" >&2
    exit 2
    ;;
esac
if [ "$#" -gt 1 ]; then
  echo "validate-release.sh: unexpected extra argument '$2' (only --test-only is accepted)" >&2
  exit 2
fi

echo "🧬 Launching Care4u Release Validation Gatekeeper..."

if [ ! -f "$VERSION_FILE" ]; then
  echo "1.0.0" > "$VERSION_FILE"
fi
CURRENT_VERSION=$(tr -d '[:space:]' < "$VERSION_FILE")
# Strict: the bump below is pure bash arithmetic, so a malformed or zero-padded value
# must be rejected here rather than emit garbage like "abc..1" at the very end.
if ! printf '%s' "$CURRENT_VERSION" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "validate-release.sh: malformed .version '$CURRENT_VERSION' (expected MAJOR.MINOR.PATCH)" >&2
  exit 1
fi
if printf '%s' "$CURRENT_VERSION" | grep -Eq '(^|\.)0[0-9]'; then
  echo "validate-release.sh: .version '$CURRENT_VERSION' has a zero-padded component" >&2
  exit 1
fi
echo "📍 Current Version: $CURRENT_VERSION"

# ---------------------------------------------------------------------------
# Go services: compile every module (catches non-compiling code even in packages
# with no tests) and run its unit tests.
# ---------------------------------------------------------------------------
if ! command -v go >/dev/null 2>&1; then
  echo "validate-release.sh: Go toolchain not found — cannot verify the services." >&2
  exit 1
fi

echo "🔨 Building & testing Go microservices..."
FOUND_SERVICE=0
for d in "$PROJECT_DIR"/services/*/; do
  [ -f "$d/go.mod" ] || continue
  FOUND_SERVICE=1
  name="$(basename "$d")"
  echo "📦 Service: $name"
  if [ -n "$(cd "$d" && gofmt -l .)" ]; then
    echo "❌ Error: gofmt reports unformatted files in services/$name:" >&2
    (cd "$d" && gofmt -l .) >&2
    exit 1
  fi
  if ! (cd "$d" && go build ./...); then
    echo "❌ Error: go build failed in services/$name!" >&2
    exit 1
  fi
  if ! (cd "$d" && go test ./...); then
    echo "❌ Error: Go unit tests failed in services/$name!" >&2
    exit 1
  fi
done
if [ "$FOUND_SERVICE" -eq 0 ]; then
  echo "validate-release.sh: no Go service module found under services/." >&2
  exit 1
fi
echo "✓ Go build + unit tests: PASS"

# ---------------------------------------------------------------------------
# Compose stack parses (a broken docker-compose.yml blocks every deploy).
# ---------------------------------------------------------------------------
echo "🐳 Validating docker-compose.yml..."
if [ -f "$PROJECT_DIR/docker-compose.yml" ] && docker compose version >/dev/null 2>&1; then
  if ! docker compose -f "$PROJECT_DIR/docker-compose.yml" config --quiet; then
    echo "❌ Error: docker-compose.yml is invalid!" >&2
    exit 1
  fi
  echo "✓ docker-compose.yml: PASS"
else
  echo "  (docker compose unavailable or no compose file — skipped)"
fi

# ---------------------------------------------------------------------------
# CDN publisher fail-closed contract.
# ci/upload_to_hf.py is what every installed OTA client reads. The 140-line
# scaffold template pins versionCode 1, advertises apk_url before upload, and lets a
# version-only publish OVERWRITE the live manifest (stripping apk_url/sha256 from
# installed clients) — silent data loss nothing else in this repo would notice.
# Runs against a STUBBED huggingface_hub in a temp dir; the real CDN is never touched,
# and --token short-circuits get_token() before it reads the live credential caches.
# ---------------------------------------------------------------------------
echo "📡 Verifying CDN publisher is fail-closed..."
if [ -f "$PROJECT_DIR/ci/upload_to_hf.py" ]; then
  if [ ! -f "$PROJECT_DIR/ci/requirements.txt" ]; then
    echo "validate-release.sh: ci/upload_to_hf.py has no ci/requirements.txt pinning huggingface_hub" >&2
    exit 1
  fi
  PY_CMD="$(command -v python3 || true)"
  if [ -z "$PY_CMD" ]; then
    echo "validate-release.sh: python3 not found — cannot verify the CDN publisher." >&2
    exit 1
  fi
  "$PY_CMD" - "$PROJECT_DIR" <<'PY'
import os, pathlib, subprocess, sys, tempfile

root = pathlib.Path(sys.argv[1])
src = root / "ci" / "upload_to_hf.py"
bad = []

with tempfile.TemporaryDirectory() as td:
    td = pathlib.Path(td)
    script = td / "upload_to_hf.py"
    script.write_bytes(src.read_bytes())
    # A stub that FAILS LOUDLY if the publisher reaches the CDN: reaching it at all
    # means a guard above it failed to fire.
    (td / "huggingface_hub.py").write_text(
        "def _boom(*a, **k):\n"
        "    raise AssertionError('publisher reached the CDN: a fail-closed guard did not fire')\n"
        "class HfApi:\n"
        "    def __init__(self, *a, **k): _boom()\n"
        "def create_repo(*a, **k): _boom()\n"
        "def upload_file(*a, **k): _boom()\n"
    )
    env = dict(os.environ, PYTHONPATH=str(td))

    def run(label, extra):
        p = subprocess.run(
            [sys.executable, str(script), "--slug", "care4u",
             "--repo", "rttss/infortts-gate-must-not-be-used", "--token", "gate-dummy-token",
             *extra],
            capture_output=True, text=True, env=env, cwd=td, timeout=60)
        if p.returncode != 2:
            detail = (p.stderr or p.stdout).strip().splitlines()
            bad.append(f"{label}: expected exit 2, got {p.returncode}"
                       + (f" — {detail[-1]}" if detail else ""))
        elif "reached the CDN" in (p.stdout + p.stderr):
            bad.append(f"{label}: publisher contacted the CDN — the guards did not fire")
        return p

    apk = td / "gate.apk"
    apk.write_bytes(b"not a real apk")

    # 1. No artifact at all: would overwrite the live manifest with one carrying no
    #    apk_url/sha256 — the silent data-loss case. Must refuse.
    run("publish with no --apk and no --patch", ["--version", "1.2.0"])
    # 2. A requested artifact not on disk: a manifest advertising a 404.
    run("publish with a non-existent --apk", ["--version", "1.2.0", "--apk",
                                              str(td / "does-not-exist.apk")])
    # 3. A version with no '+<build>' yields no positive Android versionCode, and
    #    versionCode 0 reads as a ROLLBACK to every OTA client comparing versions.
    run("publish an artifact with no build number", ["--version", "1.2.0", "--apk", str(apk)])

for line in bad:
    print(f"  {line}")
if bad:
    print(f"  {len(bad)} CDN publisher regression(s)")
sys.exit(1 if bad else 0)
PY
  echo "✓ CDN publisher refuses artifact-less and unversioned publishes (3/3)"
fi

# ---------------------------------------------------------------------------
# Secret guard: credential files and build caches must never be tracked.
# ---------------------------------------------------------------------------
echo "🔐 Checking for tracked credentials..."
if git -C "$PROJECT_DIR" rev-parse --git-dir >/dev/null 2>&1; then
  TRACKED="$(git -C "$PROJECT_DIR" ls-files)"
  LEAKED="$(printf '%s\n' "$TRACKED" | grep -E '(^|/)(key\.properties|local\.properties|google-services\.json|GoogleService-Info\.plist)$|(^|/)\.env$|\.(jks|keystore|p12|pem)$' || true)"
  if [ -n "$LEAKED" ]; then
    echo "validate-release.sh: credential files are tracked and must never be committed:" >&2
    printf '  %s\n' $LEAKED >&2
    exit 1
  fi
  CACHES="$(printf '%s\n' "$TRACKED" | grep -E '(^|/)(\.gradle|\.dart_tool|build)/|\.iml$' || true)"
  if [ -n "$CACHES" ]; then
    echo "validate-release.sh: build caches are tracked:" >&2
    printf '  %s\n' $CACHES >&2
    exit 1
  fi
  echo "✓ No tracked credential files or build caches"
else
  echo "  (not a git checkout — skipped)"
fi

# ---------------------------------------------------------------------------
# Env declaration: every env key the Go services read must be documented in
# .env.example (the deploy contract docker-compose/dev.sh follow).
# ---------------------------------------------------------------------------
echo "🔎 Checking Go env declarations..."
PY_CMD="$(command -v python3 || true)"
if [ -n "$PY_CMD" ]; then
  "$PY_CMD" - "$PROJECT_DIR" <<'PY'
import glob, os, re, sys

root = sys.argv[1]
pat = re.compile(r'(?:getEnv|Getenv|LookupEnv|GetenvOrDefault)\(\s*"([A-Z][A-Z0-9_]*)"')
names = {}
for path in glob.glob(os.path.join(root, "services", "**", "*.go"), recursive=True):
    with open(path) as f:
        body = f.read()
    for m in pat.finditer(body):
        names.setdefault(m.group(1), set()).add(os.path.relpath(path, root))
declared = open(os.path.join(root, ".env.example")).read()
missing = {k: v for k, v in names.items()
           if not re.search(r"(?m)^\s*#?\s*" + re.escape(k) + r"\s*=", declared)}
if missing:
    for k, v in sorted(missing.items()):
        print(f"  undeclared environment variable {k} (read by {', '.join(sorted(v))})")
    sys.exit(1)
print(f"✓ {len(names)} environment variable(s) declared in .env.example")
PY
else
  echo "  (python3 unavailable — skipped)"
fi

if [ "$TEST_ONLY" -eq 1 ]; then
  echo "🚀 Validation Succeeded (--test-only). .version left untouched at $CURRENT_VERSION."
  exit 0
fi

IFS='.' read -r major minor patch <<< "$CURRENT_VERSION"
NEXT_PATCH=$((patch + 1))
NEXT_VERSION="$major.$minor.$NEXT_PATCH"
echo "✨ Bumping version to: $NEXT_VERSION"
echo "$NEXT_VERSION" > "$VERSION_FILE"
echo "🚀 Validation Succeeded. Care4u is production-ready."
exit 0
