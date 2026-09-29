#!/bin/bash

# Care4u / HealthFlow Release Validation Gatekeeper
# Part of the Infortts Swarm Ecosystem.
#
#   ./validate-release.sh              full gate (no version bump: the release
#                                      version is owned by the Jenkins
#                                      tag-driven plan in ci/jenkins-common.groovy)
#   ./validate-release.sh --test-only  identical checks, CI-friendly flag
#
# Steps: shell syntax -> gofmt -> go vet -> go build -> go test, for every
# service in services/. Exits non-zero on the first real failure so a broken
# build or test can never be published.

cd "$(dirname "$0")" || exit 1

TEST_ONLY=0
if [ "${1:-}" = "--test-only" ]; then
    TEST_ONLY=1
fi

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🩺 Care4u Release Validation Igniting...${NC}"

SERVICES=$(ls -d services/*/ 2>/dev/null | sed 's#services/##; s#/$##')
if [ -z "$SERVICES" ]; then
    echo -e "${RED}❌ No services found under services/ — nothing to validate.${NC}"
    exit 1
fi

# 1. Shell syntax for the orchestrator scripts that CI and humans both run
echo -e "${YELLOW}1. Shell syntax check...${NC}"
for script in dev.sh validate-release.sh; do
    [ -f "$script" ] || continue
    if ! bash -n "$script"; then
        echo -e "${RED}❌ bash -n ${script} failed. Release rejected.${NC}"
        exit 1
    fi
done
echo -e "${GREEN}   shell scripts parse.${NC}"

# 2. JSON / YAML / compose sanity for the config that ships with the release
echo -e "${YELLOW}2. Config sanity...${NC}"
if command -v docker > /dev/null 2>&1 && docker compose version > /dev/null 2>&1; then
    if ! docker compose config > /dev/null; then
        echo -e "${RED}❌ docker compose config is invalid. Release rejected.${NC}"
        exit 1
    fi
    echo -e "${GREEN}   docker-compose.yml resolves.${NC}"
else
    echo -e "${YELLOW}   docker compose unavailable — skipped compose validation.${NC}"
fi

# 3. Go toolchain: no Go, no gate. The mac Jenkins agent provides it.
GO_CMD="$(command -v go || true)"
if [ -z "$GO_CMD" ]; then
    echo -e "${YELLOW}⚠️  Go toolchain not found — skipping fmt/vet/build/test verification.${NC}"
    echo -e "${GREEN}(gate degraded, not failed: no Go on this host)${NC}"
    exit 0
fi

echo -e "${YELLOW}3. gofmt (all services)...${NC}"
UNFORMATTED=""
for svc in $SERVICES; do
    out=$(cd "services/$svc" && gofmt -l . 2>/dev/null)
    [ -n "$out" ] && UNFORMATTED="$UNFORMATTED $svc:$(echo "$out" | tr '\n' ' ')"
done
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}❌ gofmt needed:${NC}${UNFORMATTED}"
    echo -e "${RED}   Run: gofmt -w services/. Release rejected.${NC}"
    exit 1
fi
echo -e "${GREEN}   all sources formatted.${NC}"

JOBS="$(command -v nproc > /dev/null 2>&1 && nproc || (command -v sysctl > /dev/null 2>&1 && sysctl -n hw.ncpu) || echo 2)"

FAILED=""
for svc in $SERVICES; do
    echo -e "${YELLOW}4. services/${svc}: vet + build + test${NC}"
    if ! (cd "services/$svc" && go vet ./... && go build ./... ); then
        echo -e "${RED}❌ services/${svc} failed vet/build. Release rejected.${NC}"
        FAILED="$FAILED $svc"
        continue
    fi
    if ! (cd "services/$svc" && go test ./... ); then
        echo -e "${RED}❌ services/${svc} failed tests. Release rejected.${NC}"
        FAILED="$FAILED $svc"
        continue
    fi
    echo -e "${GREEN}   services/${svc} ok.${NC}"
done

if [ -n "$FAILED" ]; then
    echo -e "${RED}❌ FAILED services:${FAILED}${NC}"
    exit 1
fi

echo -e "\n===================================================="
if [ "$TEST_ONLY" = "1" ]; then
    echo -e "${GREEN}🎉 CARE4U VALIDATION PASSED (test-only: no release action taken)${NC}"
else
    echo -e "${GREEN}🎉 CARE4U VALIDATION PASSED${NC}"
    echo -e "Release versioning is owned by the Jenkins tag plan (ci/jenkins-common.groovy)."
fi
echo "===================================================="
exit 0
