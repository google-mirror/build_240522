#! /bin/bash

# Usage:
#    generate-inheritance.sh product [MODEL] [DEVICE]
# Generate .dot file for product inheritance graph
# Inheritance pairs are read from stdin, each line contains
# <module>:<submod> pair.
# The output graph contains only the modules inherited from the given
# product.

set -eu
declare -r top="${1?specify product}"
declare -r model="${2:-}"
declare -r device="${3:-}"

# Read the edges, mapping each module to the list of submodules it inherits from.
declare -A deps
while read pair; do
	[[ ! "$pair" =~ (.+):(.+) ]] || deps["${BASH_REMATCH[1]}"]+=" ${BASH_REMATCH[2]}"
done

# Find all the reachable nodes.
declare -a nodes
declare -a candidates=( "$top" )

declare -A ready
while (( ${#candidates[@]} )); do
	n="${candidates[0]}"
	candidates=( "${candidates[@]:1}" )
	[[ ! -v ready["$n"] ]] || continue
	candidates+=( ${deps["$n"]:-} )
	ready["$n"]=1
	nodes+=( "$n" )
done

# Generate
printf "digraph {\ngraph [ ratio=.5 ];\n"
for n in "${nodes[@]}"; do
	for d in ${deps[$n]:-}; do
		printf '"%s" -> "%s"\n' $d $n
	done
done

# Top-level node has distinct color and its label displays model and device if present
printf '"%s" [ label="%s" style="filled" fillcolor="orange" colorscheme="svg" fontcolor="darkblue" ]\n' \
  "$top" "${top%/*}/\n${top##*/}\n\n$model\n$device"

for n in "${nodes[@]:1}"; do
	label="${n%/*}/\n${n##*/}\n"
	label+="\n\n"
	case $n in
	  build/make/target/product/*)
		color=beige ;;
	  vendor/*)
		color=labenderblush ;;
	  *)
		color=white ;;
	esac
	printf '"%s" [ label="%s" style="filled" fillcolor="%s" colorscheme="svg" fontcolor="darkblue" ]\n' \
	  $n $label $color
done
printf "}\n"