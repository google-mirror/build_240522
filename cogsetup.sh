function _source_env_setup_script() {
  local -r ENV_SETUP_SCRIPT="build/make/envsetup.sh"
  local -r TOP_DIR=$(
    while [[ ! -f "${ENV_SETUP_SCRIPT}" ]] && [[ "${PWD}" != "/" ]]; do
      \cd ..
    done
    if [[ -f "${ENV_SETUP_SCRIPT}" ]]; then
      echo "$(PWD= /bin/pwd -P)"
    fi
  )

  local -r FULL_PATH_ENV_SETUP_SCRIPT="${TOP_DIR}/${ENV_SETUP_SCRIPT}"
  if [[ ! -f "${FULL_PATH_ENV_SETUP_SCRIPT}" ]]; then
    echo "ERROR: Unable to source ${ENV_SETUP_SCRIPT}"
    return 1
  fi

  # Need to change directory to the repo root so vendor scripts can be sourced
  # as well.
  local -r CUR_DIR=$PWD
  \cd "${TOP_DIR}"
  source "${FULL_PATH_ENV_SETUP_SCRIPT}"
  \cd "${CUR_DIR}"
}

# This function needs to run first as the remaining defining functions may be
# using the envsetup.sh defined functions. Skip this part if this script is already
# being invoked from envsetup.sh.
if [[ "$1" != "--skip-envsetup" ]]; then
  _source_env_setup_script || return
fi

# This function detects if the uploader is available and sets the path of it to
# ANDROID_ENABLE_METRICS_UPLOAD.
function _export_metrics_uploader() {
  local uploader_path="$(gettop)/vendor/google/misc/metrics_uploader_prebuilt/metrics_uploader.sh"
  if [[ -x "${uploader_path}" ]]; then
    export ANDROID_ENABLE_METRICS_UPLOAD="${uploader_path}"
  fi
}

# This function sets RBE specific environment variables needed for the build to
# executed by RBE on Cog. This file should be sourced once per checkout of
# Android code.
function _set_rbe_vars() {
  export USE_RBE="true"
  export RBE_CXX_EXEC_STRATEGY="remote_local_fallback"
  export RBE_JAVAC_EXEC_STRATEGY="remote_local_fallback"
  export RBE_R8_EXEC_STRATEGY="remote_local_fallback"
  export RBE_D8_EXEC_STRATEGY="remote_local_fallback"
  export RBE_enable_deps_cache="true"
  export RBE_cache_dir="/tmp"
  export RBE_JAVAC=1
  export RBE_R8=1
  export RBE_D8=1
}

_export_metrics_uploader
_set_rbe_vars
