#! /bin/bash
set -eu
function die() { format=$1; shift; printf "$format\n" $@; exit 1; }
args="$@"
(($#==7)) || die "Signapk needs 7 args"
[[ "${1:-}" == "-a" ]] || die "Bad args to signapk: %s\n" "$args"
shift 2
[[ "$1" == "--align-file-size" ]] || die "Bad args to signapk: %s" "$args"
shift 1
[[ -e "$1" ]] || die "signapk public key %s does not exist" $1
shift
[[ -e "$1" ]] || die "signapk private key %s does not exist" $1
shift
[[ -e "$1" ]] || die "APK %s does not exist" $1
cp $1 $2
