#!/bin/sh

set -u

ME=$(basename $0)

USAGE="Usage: ${ME} {options}

Builds a license metadata specification and outputs it to stdout or {outfile}.

The available options are:

-k kind...              license kinds
-c condition...         license conditions
-p package...           license package name
-n notice...            license notice file
-d dependency...        license metadata file dependency
-s dependency...        source (input) dependency
-t target...            built targets
-i target...            installed targets
-r root...              root directory of project
-m target:installed...  map dependent targets to their installed names
-is_container           preserved dependent target name when given
-o outfile              output file
"

# Global flag variables
license_kinds=
license_conditions=
license_package_name=
license_notice=
license_deps=
source_deps=
targets=
installed=
installmap=
is_container=false
ofile=
roots=

# Global variables
depfiles=" "
effective_conditions=


# Exits with a message.
#
# When the exit status is 2, assumes a usage error and outputs the usage message
# to stderr before outputting the specific error message to stderr.
#
# Parameters:
#   Optional numeric exit status (defaults to 2, i.e. a usage error.)
#   Remaining args treated as an error message sent to stderr.
die() {
  lstatus=2
  case "${1:-}" in *[^0-9]*) ;; *) lstatus="$1"; shift ;; esac
  case "${lstatus}" in 2) echo "${USAGE}" >&2; echo >&2 ;; esac
  if [ -n "$*" ]; then
    echo -e "$*\n" >&2
  fi
  exit $lstatus
}


# Sets the flag variables based on the command-line.
#
# invoke with: process_args "$@"
process_args() {
  lcurr_flag=
  while [ "$#" -gt '0' ]; do
    case "${1}" in
      -h)
        echo "${USAGE}"
        exit 0
        ;;
      -k)
        lcurr_flag=kind
        ;;
      -c)
        lcurr_flag=condition
        ;;
      -p)
        lcurr_flag=package
        ;;
      -n)
        lcurr_flag=notice
        ;;
      -d)
        lcurr_flag=dependency
        ;;
      -s)
        lcurr_flag=source
        ;;
      -t)
        lcurr_flag=target
        ;;
      -i)
        lcurr_flag=install
        ;;
      -m)
        lcurr_flag=installmap
        ;;
      -o)
        lcurr_flag=ofile
        ;;
      -r)
        lcurr_flag=root
        ;;
      -is_container)
        lcurr_flag=
        is_container=true
        ;;
      -*)
        die "Unknown flag: \"${1}\""
        ;;
      *)
        case "${lcurr_flag}" in
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
          source)
            source_deps="${source_deps}${source_deps:+ }${1}"
            ;;
          target)
            targets="${targets}${targets:+ }${1}"
            ;;
          install)
            installed="${installed}${installed:+ }${1}"
            ;;
          installmap)
            installmap="${installmap}${installmap:+ }${1}"
            ;;
          root)
            root="${1}"
            while [ -n "${root}" ] && ! [ "${root}" == '.' ] && \
                ! [ "${root}" == '/' ]; \
            do
              if [ -d "${root}/.git" ]; then
                roots="${roots}${roots:+ }${root}"
                break
              fi
              root=$(dirname "${root}")
            done
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

## Returns the effective license conditions for the current license metadata.
##
## If a module is restricted or links in a restricted module, the effective
## license has a restricted condition.
#calculate_effective_conditions() {
#  lconditions="${license_conditions}"
#  case "${license_conditions}" in
#    *restricted*) : do nothing ;;
#    *)
#       for d in ${depfiles}; do
#         if cat "${d}" | egrep -q 'effective_condition\s*:.*restricted' ; then
#           lconditions="${lconditions}${lconditions:+ }restricted"
#           break
#         fi
#       done
#     ;;
#  esac
#  echo "${lconditions}"
#}


## Returns the effective license conditions for the current license metadata.
##
## If a module is restricted or links in a restricted module, the effective
## license has a restricted condition.
#calculate_effective_license_text() {
#  ltexts=$(( \
#      for t in ${license_notice}; do
#        echo "\"${t}:${license_package_name}/$(basename ${t})\""
#      done
#      for d in ${depfiles}; do
#        cat "${d}" | awk '
#          BEGIN {
#            pn = ""
#          }
#          $1 == "license_package_name:" {
#            pn = $2
#            sub(/^"/,"",pn)
#            sub(/"$/,"",pn)
#          }
#          $1 == "license_text:" {
#            gsub(/^"|"$/, "", $2)
#            bn = $2
#            sub("^.*[/]", "", bn)
#            print "\""$2":"pn"/"bn"\""
#          }
#          $1 == "effective_license_text:" {
#            print $2
#          }
#        '
#      done
#  ) | awk '!seen[$0]++')
#  echo "${ltexts}"
#}


process_args "$@"

if [ -n "${ofile}" ]; then
  # truncate the output file before appending results
  : >"${ofile}"
else
  ofile=/dev/stdout
fi

# spit out the license metadata file content
(
  echo 'license_package_name: "'"${license_package_name}"'"'
  for r in ${roots}; do
    echo 'root: "'"${r}"'"'
  done
  for kind in ${license_kinds}; do
    echo 'license_kind: "'"${kind}"'"'
  done
  for condition in ${license_conditions}; do
    echo 'license_condition: "'"${condition}"'"'
  done
  for f in ${license_notice}; do
    echo 'license_text: "'"${f}"'"'
  done
  echo "is_container: ${is_container}"
  for t in ${targets}; do
    echo 'built: "'"${t}"'"'
  done
  for i in ${installed}; do
    echo 'installed: "'"${i}"'"'
  done
  for m in ${installmap}; do
    echo 'install_map: "'"${m}"'"'
  done
  for s in ${source_deps}; do
    echo 'source: "'"${s}"'"'
  done
) >>"${ofile}"
depfiles=" $(echo $(echo ${license_deps} | tr ' ' '\n' | sort -u)) "
#effective_conditions=$(calculate_effective_conditions)
#for condition in ${effective_conditions}; do
#  echo 'effective_condition: "'"${condition}"'"'
#done >>"${ofile}"
#while read t; do
#  if [ -n "${t}" ]; then
#    echo "effective_license_text: ${t}"
#  fi
#done < <(calculate_effective_license_text) >>"${ofile}"
#effective_license_text=$(calculate_effective_license_text)
#for text in ${effective_license_text}; do
#  echo 'effective_license_text: "'"${text}"'"'
#done >>"${ofile}"
#subdeps=
for dep in ${depfiles}; do
  echo 'dep: "'"${dep}"'"'
#  if ${is_container}; then
#    subdeps="${subdeps}${subdeps:+ }"$(cat "${dep}" | \
#      awk '
#        function dequote() {
#          sub(/^"/,"",$2)
#          sub(/"$/,"",$2)
#        }
#        $1 == "dep:" {
#          dequote()
#          print $2
#        }
#        $1 == "sub_dep:" {
#          dequote()
#          print $2
#        }
#     '
#    )
#  fi
#  echo 'dep {'
#  cat "${dep}" | \
#    awk -v name="${dep}" '
#      function strip_type() {
#        $1 = ""
#        gsub(/^\s+/, "", $0)
#      }
#      BEGIN {
#        print "  name: \""name"\""
#        indep=0
#      }
#      $1 == "dep" {
#        indep=1
#      }
#      $1 == "}" {
#        indep=0
#      }
#      $1 == "license_package_name:" && !indep {
#        strip_type()
#        print "  license_package_name:"$0
#      }
#      $1 == "root:" && !indep {
#        print "  root: "$2
#      }
#      $1 == "license_kind:" && !indep {
#        print "  license_kind: "$2
#      }
#      $1 == "license_condition:" && !indep {
#        print "  license_condition: "$2
#      }
#      $1 == "is_container:" && !indep {
#        print "  is_container: "$2
#      }
#      $1 == "license_text:" && !indep {
#        strip_type()
#        print "  license_text:"$0
#      }
#      $1 == "effective_license_text:" && !indep {
#        strip_type()
#        print "  effective_license_text:"$0
#      }
#      $1 == built:" && !indep {
#        print "  built: "$2
#      }
#      $1 == "installed:" && !indep {
#        print "  installed: "$2
#      }
#      $1 == "install_map:" && !indep {
#        print "  install_map: "$2
#      }
#      $1 == "source:" && !indep {
#        print "  source: "$2
#      }
#      $1 == "name:" && indep {
#        print "  sub_dep: "$2
#      }
#  '
#  # The restricted license kind is contagious to all linked dependencies.
#  dep_conditions=$(echo $(
#      cat "${dep}" | awk '
#        BEGIN {
#          indep=0
#        }
#        $1 == "dep" {
#          indep=1
#        }
#        $1 == "}" {
#          indep=0
#        }
#        $1 == "effective_condition:" && !indep {
#          $1 = ""
#          sub(/^\s*/, "")
#          gsub(/"/, "")
#          print
#        }
#      '
#  ))
#  for condition in ${dep_conditions}; do
#    echo '  effective_condition: "'${condition}'"'
#  done
#  if ! ${is_container}; then
#    case "${dep_conditions}" in
#      *restricted*) : already restricted -- nothing to inherit ;;
#      *)
#        case "${effective_conditions}" in
#          *restricted*)
#            # "contagious" restricted infects everything linked to restricted
#            echo '  effective_condition: "restricted"'
#            ;;
#        esac
#        ;;
#    esac
#  fi
#  echo '}'
done >>"${ofile}"
#if ${is_container}; then
#  subdeps=" $(echo $(echo ${subdeps} | tr ' ' '\n' | sort -u)) "
#  for sd in ${subdeps}; do
#    case "${depfiles}" in
#      *" ${sd} "*) : do nothing ;;
#      *) echo 'sub_dep: '"${sd}" ;;
#    esac
#  done >>"${ofile}"
#fi
