#!/usr/bin/env bash

set -e
set -o pipefail

OS="$(uname -s)"
ARCH="$(uname -m)"

PROGRAM_NAME="archy"
REPOSITORY="xargs-dev/${PROGRAM_NAME}"

WORKDIR="$(mktemp -d /tmp/bindl-bootstrap-XXXXX)"

ARCHIVE="${PROGRAM_NAME}_0.1.2_${OS}_${ARCH}.tar.gz"

function log() {
  echo -e "[\e[1;34mbootstrap\e[0m] $1"
}

function prompt() {
  read -p "Proceed? (y/N) " answer
  if [ $answer != "y" ]; then
    echo "Aborted: only 'y' is accepted answer to continue (received '${answer}')"
    exit 1
  fi
}

log "Hello! The sole purpose of my existence is to bootstrap bindl."
log "I have found myself in ${ARCH} machine running ${OS}."
log "I expect the archive to be named ${ARCHIVE}."

prompt

log "Working in ${WORKDIR}"
pushd "${WORKDIR}" >/dev/null
  log "Downloading (1/2): checksums.txt"
  curl --silent --location --remote-name "https://github.com/${REPOSITORY}/releases/latest/download/checksums.txt"

  log "Downloading (2/2): ${ARCHIVE}"
  curl --silent --location --remote-name "https://github.com/${REPOSITORY}/releases/latest/download/${ARCHIVE}"

  log "Verifying checksums"
  shasum --algorithm 256 --check checksums.txt --ignore-missing

  tar -xzf ${ARCHIVE} ${PROGRAM_NAME}

  log "Printing program version"
  ./${PROGRAM_NAME} version
popd >/dev/null

trap "rm -r ${WORKDIR}" EXIT

mv ${WORKDIR}/${PROGRAM_NAME} .
log "Done! The binary is in current working directory."
log "Move them to a directory accessible in PATH to start using seamlessly."
