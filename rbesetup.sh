#!/bin/bash

# This function prefixes the given command with appropriate variables needed
# for the build to be executed with RBE.
function use_rbe() {
  RBE_LOG_DIR="/tmp"
  DOCKER_IMAGE="gcr.io/androidbuild-re-dockerimage/android-build-remoteexec-image@sha256:582efb38f0c229ea39952fff9e132ccbe183e14869b39888010dacf56b360d62"
  RBE_BINARIES_DIR="prebuilts/remoteexecution-client/latest/"

  # Do not set an invocation-ID and let reproxy auto-generate one.
  USE_RBE="true" \
  FLAG_use_application_default_credentials="true" \
  FLAG_server_address="unix:///tmp/reproxy_$RANDOM.sock" \
  NINJA_REMOTE_NUM_JOBS="500" \
  FLAG_exec_root="$PWD" \
  FLAG_log_dir="${RBE_LOG_DIR}" \
  FLAG_reproxy_wait_seconds="20" \
  FLAG_output_dir="${RBE_LOG_DIR}" \
  FLAG_shutdown_proxy="true" \
  FLAG_log_path="text://${RBE_LOG_DIR}/reproxy_log.txt" \
  FLAG_exec_strategy="remote_local_fallback" \
  FLAG_platform="container-image=docker://${DOCKER_IMAGE},jdk-version=10" \
  FLAG_cpp_dependency_scanner_plugin="${RBE_BINARIES_DIR}/dependency_scanner_go_plugin.so" \
  RBE_DIR=${RBE_BINARIES_DIR} \
  FLAG_re_proxy="${RBE_BINARIES_DIR}/reproxy" \
  $@
}
