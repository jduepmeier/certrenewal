#!/usr/bin/env bash
#
# Will be called from go test

set -eo pipefail

BASEDIR="${1:-}"

main() {
    if [[ -z "${BASEDIR}" ]]
    then
        ( >&2 echo "missing base dir. Usage: $0 <basedir>" )
        exit 1
    fi

    DIR="${BASEDIR}/ssh-keys"

    mkdir -p "${DIR}"

    KEY_TYPES=(
        "rsa"
        "ed25519"
    )

    for type in "${KEY_TYPES[@]}"
    do
        private_key="${DIR}/${type}"
        public_key="${DIR}/${type}.pub"
        cert_key="${DIR}/${type}-cert.pub"
        ca_private_key="${DIR}/ca_${type}"
        ca_public_key="${DIR}/ca_${type}.pub"

        # CA key. 
        ssh-keygen -t "${type}" -f "${ca_private_key}" -q -N ""
        # First generate private key.
        ssh-keygen -t "${type}" -f "${private_key}" -q -N ""
        # Sign +30 days
        ssh-keygen -s "${ca_private_key}" -I "30days" -V "-30d:+30d" "${public_key}"
        mv "${cert_key}" "${DIR}/${type}-cert-30days.pub"
        ssh-keygen -s "${ca_private_key}" -I "1day" -V "-1d:+1d" "${public_key}"
        mv "${cert_key}" "${DIR}/${type}-cert-1day.pub"
        ssh-keygen -s "${ca_private_key}" -I "forever" "${public_key}"
        mv "${cert_key}" "${DIR}/${type}-cert-forever.pub"
    done
}
main "$@"