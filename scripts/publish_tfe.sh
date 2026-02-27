#!/usr/bin/env bash
set -euo pipefail

PROVIDER_NAME="atlassian"
REGISTRY_NAME="private"
PROTOCOLS='["5.0","6.0"]'

die() {
  echo "error: $*" >&2
  exit 1
}

need_env() {
  local key=$1
  local val=${!key-}
  if [[ -z "${val}" ]]; then
    die "Missing required env: ${key}"
  fi
}

if ! command -v jq >/dev/null 2>&1; then
  die "jq is required"
fi

if [[ $# -lt 1 ]]; then
  die "usage: scripts/publish_tfe.sh v0.1.0|0.1.0"
fi

tag="$1"
if [[ "${tag}" != v* ]]; then
  tag="v${tag}"
fi
version="${tag#v}"

TFE_HOST="${TFE_HOST:-app.terraform.io}"
TFE_HOST="${TFE_HOST%%/}"
TFE_BASE="https://${TFE_HOST}/api/v2"

need_env TFE_TOKEN
need_env TFE_ORG
GPG_KEY_ID="${TFE_GPG_KEY_ID:-${GPG_KEY_ID-}}"
if [[ -z "${GPG_KEY_ID}" ]]; then
  die "Missing GPG_KEY_ID (or TFE_GPG_KEY_ID)"
fi

need_env GITHUB_TOKEN
GITHUB_OWNER="${GITHUB_OWNER:-surajsinghrajput}"
GITHUB_REPO="${GITHUB_REPO:-terraform-provider-atlassian}"

tmpdir=$(mktemp -d)
trap 'rm -rf "${tmpdir}"' EXIT

echo "Using tag=${tag}, version=${version}, org=${TFE_ORG}, host=${TFE_HOST}"

echo "Fetching GitHub release..."
release_json="${tmpdir}/release.json"
curl -fsSL \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/tags/${tag}" \
  -o "${release_json}"

asset_url_for() {
  local name=$1
  jq -r --arg NAME "${name}" '.assets[] | select(.name == $NAME) | .url' "${release_json}" | head -n1
}

shasums_name="terraform-provider-atlassian_${version}_SHA256SUMS"
shasums_sig_name="${shasums_name}.sig"

shasums_url="$(asset_url_for "${shasums_name}")"
shasums_sig_url="$(asset_url_for "${shasums_sig_name}")"

if [[ -z "${shasums_url}" || "${shasums_url}" == "null" ]]; then
  die "Release missing ${shasums_name}"
fi
if [[ -z "${shasums_sig_url}" || "${shasums_sig_url}" == "null" ]]; then
  die "Release missing ${shasums_sig_name}"
fi

echo "Downloading SHASUMS and signature..."
curl -fsSL \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "Accept: application/octet-stream" \
  "${shasums_url}" -o "${tmpdir}/${shasums_name}"

curl -fsSL \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "Accept: application/octet-stream" \
  "${shasums_sig_url}" -o "${tmpdir}/${shasums_sig_name}"

echo "Creating provider (idempotent)..."
create_provider_body=$(cat <<EOF
{
  "data": {
    "type": "registry-providers",
    "attributes": {
      "name": "${PROVIDER_NAME}",
      "namespace": "${TFE_ORG}",
      "registry-name": "${REGISTRY_NAME}"
    }
  }
}
EOF
)

provider_resp="${tmpdir}/provider.json"
provider_code=$(curl -sS -o "${provider_resp}" -w '%{http_code}' \
  -H "Authorization: Bearer ${TFE_TOKEN}" \
  -H "Content-Type: application/vnd.api+json" \
  -X POST "${TFE_BASE}/organizations/${TFE_ORG}/registry-providers" \
  -d "${create_provider_body}")

if [[ "${provider_code}" != "200" && "${provider_code}" != "201" && "${provider_code}" != "422" ]]; then
  echo "Response:" >&2
  sed -e 's/$/\\n/' "${provider_resp}" | head -c 300 >&2
  die "Create provider failed (status ${provider_code})"
fi

echo "Creating provider version ${version}..."
create_version_body=$(cat <<EOF
{
  "data": {
    "type": "registry-provider-versions",
    "attributes": {
      "version": "${version}",
      "key-id": "${GPG_KEY_ID}",
      "protocols": ${PROTOCOLS}
    }
  }
}
EOF
)

version_resp="${tmpdir}/version.json"
version_code=$(curl -sS -o "${version_resp}" -w '%{http_code}' \
  -H "Authorization: Bearer ${TFE_TOKEN}" \
  -H "Content-Type: application/vnd.api+json" \
  -X POST "${TFE_BASE}/organizations/${TFE_ORG}/registry-providers/${REGISTRY_NAME}/${TFE_ORG}/${PROVIDER_NAME}/versions" \
  -d "${create_version_body}")

if [[ "${version_code}" != "200" && "${version_code}" != "201" ]]; then
  echo "Response:" >&2
  sed -e 's/$/\\n/' "${version_resp}" | head -c 300 >&2
  die "Create version failed (status ${version_code})"
fi

shasums_upload=$(jq -r '.data.links["shasums-upload"]' "${version_resp}")
shasums_sig_upload=$(jq -r '.data.links["shasums-sig-upload"]' "${version_resp}")

if [[ -z "${shasums_upload}" || "${shasums_upload}" == "null" ]]; then
  die "Version response missing shasums-upload link"
fi
if [[ -z "${shasums_sig_upload}" || "${shasums_sig_upload}" == "null" ]]; then
  die "Version response missing shasums-sig-upload link"
fi

echo "Uploading SHASUMS..."
curl -fsS -X PUT -H "Content-Type: application/octet-stream" \
  --data-binary @"${tmpdir}/${shasums_name}" \
  "${shasums_upload}"

echo "Uploading SHASUMS signature..."
curl -fsS -X PUT -H "Content-Type: application/octet-stream" \
  --data-binary @"${tmpdir}/${shasums_sig_name}" \
  "${shasums_sig_upload}"

platforms_url="${TFE_BASE}/organizations/${TFE_ORG}/registry-providers/${REGISTRY_NAME}/${TFE_ORG}/${PROVIDER_NAME}/versions/${version}/platforms"

echo "Creating platforms and uploading binaries..."
while read -r sha filename; do
  [[ -z "${sha}" || -z "${filename}" ]] && continue
  [[ "${filename}" != terraform-provider-atlassian_"${version}"_*_*.zip ]] && continue

  asset_url="$(asset_url_for "${filename}")"
  if [[ -z "${asset_url}" || "${asset_url}" == "null" ]]; then
    die "Release missing asset ${filename}"
  fi

  base="${filename%.zip}"
  tmp="${base%_*}"
  goarch="${base##*_}"
  goos="${tmp##*_}"

  platform_body=$(cat <<EOF
{
  "data": {
    "type": "registry-provider-platforms",
    "attributes": {
      "os": "${goos}",
      "arch": "${goarch}",
      "shasum": "${sha}",
      "filename": "${filename}"
    }
  }
}
EOF
)

  platform_resp="${tmpdir}/platform-${goos}-${goarch}.json"
  platform_code=$(curl -sS -o "${platform_resp}" -w '%{http_code}' \
    -H "Authorization: Bearer ${TFE_TOKEN}" \
    -H "Content-Type: application/vnd.api+json" \
    -X POST "${platforms_url}" \
    -d "${platform_body}")

  if [[ "${platform_code}" != "200" && "${platform_code}" != "201" ]]; then
    echo "Response:" >&2
    sed -e 's/$/\\n/' "${platform_resp}" | head -c 300 >&2
    die "Create platform ${filename} failed (status ${platform_code})"
  fi

  upload_url=$(jq -r '.data.links["provider-binary-upload"]' "${platform_resp}")
  if [[ -z "${upload_url}" || "${upload_url}" == "null" ]]; then
    die "No provider-binary-upload link for ${filename}"
  fi

  echo "Uploading ${filename} (${goos}/${goarch})..."
  curl -fsSL \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "Accept: application/octet-stream" \
    "${asset_url}" -o "${tmpdir}/${filename}"

  curl -fsS -X PUT -H "Content-Type: application/octet-stream" \
    --data-binary @"${tmpdir}/${filename}" \
    "${upload_url}"
done < "${tmpdir}/${shasums_name}"

echo "Published ${PROVIDER_NAME} ${version} to TFE private registry (org=${TFE_ORG})"

