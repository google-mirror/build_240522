#!/bin/bash

set -u

ME=$(basename $0)

USAGE="Usage: ${ME} {options}

Builds a license metadata specification and outputs it to stdout.

The available options are:

-k kind...         license kinds
-c condition...    license conditions
-p package...      license package name
-n notice...       license notice file
-d dependency...   license metadata file dependency
-t target...       targets
-is_container      preserved dependent target name when given
-o outfile         output file
"

# Global flag variables
license_kinds=
license_conditions=
license_package_name=
license_notice=
license_deps=
targets=
is_container=false
ofile=


# Global variables
declare -A depfiles
effective_kind=

# work around bug where "${#array[@]}" traps if never set
depfiles["x"]=808      # set an arbitrary index
unset depfiles["x"]    # delete it to make depfiles empty again


# Exits with a message.
#
# When the exit status is 2, assumes a usage error and outputs the usage message
# to stderr before outputting the specific error message to stderr.
#
# Parameters:
#   Optional numeric exit status (defaults to 2, i.e. a usage error.)
#   Remaining args treated as an error message sent to stderr.
function die() {
  local status=2
  case "${1:-}" in *[^0-9]*) ;; *) status="$1"; shift ;; esac
  case "${status}" in 2) echo "${USAGE}" >&2; echo >&2 ;; esac
  if [ -n "$*" ]; then
    echo -e "$*\n" >&2
  fi
  exit $status
}


# Sets the flag variables based on the command-line.
#
# invoke with: process_args "$@"
function process_args() {
  local curr_flag=
  while [ "$#" -gt '0' ]; do
    case "${1}" in
      -k)
        curr_flag=kind
        ;;
      -c)
        curr_flag=condition
        ;;
      -p)
        curr_flag=package
        ;;
      -n)
        curr_flag=notice
        ;;
      -d)
        curr_flag=dependency
        ;;
      -t)
        curr_flag=target
        ;;
      -o)
        curr_flag=ofile
        ;;
      -is_container)
        is_container=true
        ;;
      -*)
        die "Unknown flag: \"${1}\""
        ;;
      *)
        case "${curr_flag}" in
          kind)
            license_kinds="${license_kinds}${license_kinds:+ }${1}"
            ;;
          condition)
            license_conditions="${license_conditions}${license_conditions:+ }${1}"
            ;;
          package)
            license_package_name="${license_package_name}${license_package_name:+ }${1}"
            ;;
          notice)
            license_notice="${license_notice}${license_notice:+ }${1}"
            ;;
          dependency)
            license_deps="${license_deps}${license_deps:+ }${1}"
            ;;
          target)
            targets="${targets}${targets:+ }${1}"
            ;;
          ofile)
            if [ -n "${ofile}" ]; then
              die "Output file -o appears twice as \"${ofile}\" and \"${1}\""
            fi
            ofile="${1}"
            ;;
          *)
            die "Must precede argument \"${1}\" with type flag."
            ;;
        esac
        ;;
    esac
    shift
  done
}

function most_restrictive() {
  local kind="UNKNOWN"
  while [ "$#" -gt '0' ]; do
    case "${kind}" in
      UNKNOWN)
        case "${1}" in
          RESTRICTED*) kind="RESTRICTED" ;;
          *) kind="${1}" ;;
        esac
        ;;
      UNENCUMBERED)
        case "${1}" in
          UNKNOWN) : do nothing ;;
          RESTRICTED*) kind="RESTRICTED" ;;
          *) kind="${1}" ;;
        esac
        ;;
      NOTICE)
        case "${1}" in
          UNKNOWN|UNENCUMBERED) : do nothing ;;
          RESTRICTED*) kind="RESTRICTED" ;;
          *) kind="${1}" ;;
        esac
        ;;
      RECIPROCAL)
        case "${1}" in
          UNKNOWN|UNENCUMBERED|NOTICE) : do nothing ;;
          RESTRICTED*) kind="RESTRICTED" ;;
          *) kind="${1}" ;;
        esac
        ;;
      RESTRICTED*)
        case "${1}" in BY_EXCEPTION_ONLY|NOT_ALLOWED) kind="${1}" ;; esac ;;
      BY_EXCEPTION_ONLY)
        case "${1}" in NOT_ALLOWED) kind="${1}" ;; esac ;;
      NOT_ALLOWED)
        break
        ;;
    esac
    shift
  done
  echo "${kind}"
}

function extract_deps() {
  awk '$1 == "dep_name:" { sub(/^"/, "", $2); sub(/"$/, "", $2); print $2 }'
}

function read_deps() {
  local newdeps=$(
    for d in ${license_deps}; do
      case "${d}" in
        *.meta_module) cat "${d}" ;;
        *) echo "${d}" ;;
      esac
    done
  )
  local alldeps=
  local deps=
  local content=
  local mod=
  while [ "${#newdeps}" -gt '0' ]; do
    deps="${newdeps}"
    newdeps=
    for dep in ${deps}; do
      content=$(cat ${dep})
      depfiles["${dep}"]="${content}"
      alldeps="${alldeps}${alldeps:+ }"$(echo "${content}" | extract_deps)
    done
    alldeps=$(for d in ${alldeps}; do echo "${d}"; done | sort -u)
    for d in ${alldeps}; do
      mod="${d}"
      case "${d}" in *.meta_module) mod=$(cat ${d}) ;; esac
      if [ -z "${depfiles[${mod}]+isset}" ]; then
        newdeps="${newdeps}${newdeps:+ }${mod}"
      fi
    done
    alldeps=
  done
}

function calculate_effective_kind() {
  local kind=
  case "${license_kinds}" in
    *RESTRICTED*) echo "RESTRICTED" ;;
    *)
       for d in "${!depfiles[@]}"; do
         kind=$(
             echo "${depfiles[${d}]}" | \
                 awk '$1 == "effective_kind:" { print $2 }' \
         )
         case "${kind}" in
           *RESTRICTED*)
             echo "RESTRICTED"
             break
             ;;
         esac
       done
       case "${kind}" in
         *RESTRICTED*) : do nothing ;;
         *) echo $(most_restrictive ${license_kinds}) ;;
       esac
     ;;
  esac
}


process_args "$@"

if [ -n "${ofile}" ]; then
  echo -n >"${ofile}"
else
  ofile=/dev/stdout
fi

echo 'license_package_name: "'${license_package_name}'"' >>"${ofile}"
echo 'license_kind: "'$(most_restrictive ${license_kinds})'"' >>"${ofile}"
echo 'license_conditions: "'${license_conditions}'"' >>"${ofile}"
for f in ${license_notice}; do
  echo 'license_text: "'${f}'"' >>"${ofile}"
done
read_deps
effective_kind=$(calculate_effective_kind)
echo 'effective_kind: "'${effective_kind}'"' >>"${ofile}"
for t in ${targets}; do
  echo 'target: "'${t}'"' >>"${ofile}"
done
for dep in "${!depfiles[@]}"; do
  echo 'dep {' >>"${ofile}"
  echo '  dep_name: "'${dep}'"' >>"${ofile}"
  echo '  dep_package_name:'$(
      echo "${depfiles[${dep}]}" | \
          awk '
            $1 == "license_package_name:" {
              $1 = ""
              gsub(/^\s*/, "")
              print
            }
          '
  ) >>"${ofile}"
  echo '  dep_license_kind: '$(
      echo "${depfiles[${dep}]}" | awk '$1 == "license_kind:" { print $2 }'
  ) >>"${ofile}"
  echo '  dep_license_conditions:'$(
      echo "${depfiles[${dep}]}" | \
          awk '
            $1 == "license_conditions:" {
              $1 = ""
              gsub(/^\s*/, "")
              print
            }
          '
  ) >>"${ofile}"
  echo "${depfiles[${dep}]}" | awk '
    $1 == "license_text:" {
      printf "  dep_license_text:"
      $1 = ""
      gsub(/^\s*/, "")
      print
    }
  ' >>"${ofile}"
  case "${effective_kind}" in
    RESTRICTED)
      echo '  dep_effective_kind: "RESTRICTED"' >>"${ofile}"
      ;;
    *)
      echo '  dep_effective_kind: '$(
          echo "${depfiles[${dep}]}" | awk '$1 == "effective_kind:" { print $2 }'
      ) >>"${ofile}"
      ;;
  esac
  if ${is_container} || [ -z "${targets}" ]; then
    echo "${depfiles[${dep}]}" | awk '
      $1 == "target:" {
        print "  dep_target: "$2
      }
    ' >>"${ofile}"
  else
    for t in ${targets}; do
      echo '  dep_target: "'${t}'"' >>"${ofile}"
    done
  fi
  echo '}' >>"${ofile}"
done
