#!/bin/bash
#
# wrapper to convert MAKEFLAGS to Ninja command line options
#
# getopt(1) does not handle --longoptions "jobserver-fds:" correctly
#echo "MAKEFLAGS:  '${MAKEFLAGS}'"
_makeflags="${MAKEFLAGS/--jobserver-fds=*,* /}"
# fix non-empty MAKEFLAGS that starts without "-"
if [[ -n "${_makeflags}" ]] &&
   [[ ! "${_makeflags}" =~ ^- ]]; then
    _makeflags="-${_makeflags}"
fi
#echo "MAKEFLAGS:  '${_makeflags}'"

# parse and loop over MAKEFLAGS options
_parsed=( $(getopt --quiet --unquoted --options "jkl:w" -- "${_makeflags}") )
_i=0
_additional_args=()
while /bin/true; do
    #echo "ARGUMENT ${_i}: '${_parsed[${_i}]}'"
    case "${_parsed[${_i}]}" in
	-j) # Ninja with patch for PR#1139
	    # parallelism limited by GNU make jobserver
	    _additional_args+=( -j10000 )
	    ;;
	-k)
	    _additional_args+=( -k0 )
	    ;;
	-l)
	    let _i=_i+1
	    _additional_args+=( -l${_parsed[${_i}]} )
	    ;;
	--) # end of command line parameters
	    break
	    ;;
	*) # ignore everything else
	    ;;
    esac
    let _i=_i+1
done

# and over to Ninja...
echo $0: "$@" "${_additional_args[@]}"
exec "$@" "${_additional_args[@]}"
