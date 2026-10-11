#!/usr/bin/env bash
# Builds the release archives into dist/, named per the appkit contract:
#   tea-<version>-macos-x86_64.tar.gz and tea-<version>-macos-aarch64.tar.gz,
# each holding tea-<version>-macos-<arch>/{tea,README.md,LICENSE,NOTICE,THIRD_PARTY_NOTICES}.
# The executables are ad-hoc signed where codesign exists (the release runs on macOS).
#
# usage: scripts/package.sh VERSION
set -euo pipefail
cd "$(dirname "$0")/.."

version=${1:?usage: scripts/package.sh VERSION}

# Licenses of everything linked into the executable: the Go runtime and
# standard library, then every module the macOS build imports.
third_party_notices() {
    printf 'tea links the following third-party code. Their licenses follow.\n'
    printf '\n==== Go standard library and runtime (%s)\n\n' "$(go env GOVERSION)"
    cat "$(go env GOROOT)/LICENSE"
    local modules
    modules=$(for goarch in amd64 arm64; do
        GOOS=darwin GOARCH=$goarch go list -deps \
            -f '{{with .Module}}{{if not .Main}}{{.Path}}@{{.Version}}{{"\t"}}{{.Dir}}{{end}}{{end}}' .
    done | sort -u)
    while IFS=$'\t' read -r module dir; do
        local files=()
        for f in "$dir"/LICENSE* "$dir"/LICENCE* "$dir"/COPYING* "$dir"/NOTICE* "$dir"/PATENTS*; do
            [[ -f $f ]] && files+=("$f")
        done
        if ((${#files[@]} == 0)) && [[ $module == github.com/soundadam/* ]]; then
            # soundadam's own modules carry their licenses at the repository root, outside the Go module
            printf '\n==== %s\n\nSame author as tea; licensed as stated in https://%s\n' \
                "$module" "$(cut -d/ -f1-3 <<<"$module")"
            continue
        fi
        if ((${#files[@]} == 0)); then
            echo "no license file in $module ($dir)" >&2
            return 1
        fi
        for f in "${files[@]}"; do
            printf '\n==== %s (%s)\n\n' "$module" "$(basename "$f")"
            cat "$f"
        done
    done <<<"$modules"
}

rm -rf dist
mkdir -p dist
third_party_notices >dist/THIRD_PARTY_NOTICES

for arch in x86_64 aarch64; do
    case $arch in
    x86_64) goarch=amd64 ;;
    aarch64) goarch=arm64 ;;
    esac
    name="tea-$version-macos-$arch"
    mkdir -p "dist/$name"
    CGO_ENABLED=0 GOOS=darwin GOARCH=$goarch \
        go build -trimpath -ldflags "-s -w -X main.version=$version" -o "dist/$name/tea" .
    # Go signs arm64 ad hoc by itself; x86_64 needs it done here.
    if command -v codesign >/dev/null; then
        codesign --force --sign - "dist/$name/tea"
    fi
    cp README.md LICENSE NOTICE dist/THIRD_PARTY_NOTICES "dist/$name/"
    COPYFILE_DISABLE=1 tar -C dist -czf "dist/$name.tar.gz" "$name"
    rm -rf "dist/${name:?}"
done
rm dist/THIRD_PARTY_NOTICES

ls -l dist
