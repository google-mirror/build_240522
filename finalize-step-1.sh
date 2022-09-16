#!/bin/bash
# Automation for finalize_branch_for_release.sh.
# Sets up local environment, runs the finalization script and submits the results.
# WIP:
# - does not revert the results of the previous runs,
# - does not submit, only sends to gerrit.

# set -ex

function revert_local_changes() {
    repo forall -c "git checkout . ; git clean -fdx ; git checkout @ ; git b fina-step1 -D ; repo start fina-step1 ; git checkout @ ; git b fina-step1 -D"
}

function finalize() {
    source $top/build/make/finalize_branch_for_release.sh
}

function create_branches_and_cls() {
    repo forall -c 'if [[ $(git status --short) ]]; then repo start fina-step1 ; git add -A . ; git commit -m FINALIZATION_STEP_1_SCRIPT_COMMIT -m WILL_BE_AUTOMATICALLY_REVERTED ; repo upload --cbr --no-verify -t -y . ; fi'
}

function finalize_step_1_main() {
    local top="$(dirname "$0")"/../..

    repo selfupdate

    # revert local changes
    revert_local_changes

    # vndk etc finalization
    finalize

    # move all changes to fina-step1 branch and commit with a robot message
    create_branches_and_cls
}

finalize_step_1_main
