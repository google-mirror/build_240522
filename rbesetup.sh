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
CUR_DIR=$PWD
  \cd "${TOP_DIR}"
  source "${FULL_PATH_ENV_SETUP_SCRIPT}"
"${FULL_PATH_ENV_SETUP_SCRIPT}"


/reproxy_$RANDOM.sock" \LAG_exec_root="$(gettop)" \
platform="container-image=docker://${DOCKER_IMAGE}" \
R_IMAGE}" \
\
eproxy_wait_seconds="20" \


text://${RBE_LOG_DIR}/reproxy_log.txt" \
EGY="remote_local_fallback" \
DIR=${RBE_BINAE_BINARIES_DIR}/reproxy" \
  $@
}

