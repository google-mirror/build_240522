#!/bin/bash -e

cd $(dirname $BASH_SOURCE)
source $(pwd)/../../shell_utils.sh
require_top

# Ensure cogsetup (out/ will be symlink outside the repo)
. ${TOP}/build/make/cogsetup.sh

case $(uname -s) in
    Linux)
      export PREBUILTS_CLANG_TOOLS_ROOT="${TOP}/prebuilts/clang-tools/linux-x86/"
      PREBUILTS_GO_ROOT="${TOP}/prebuilts/go/linux-x86/"
      ;;
    *)
      echo "Only supported for linux hosts" >&2
      exit 1
      ;;
esac

export ANDROID_BUILD_TOP=$TOP
export OUT_DIR=${OUT_DIR}
exec "${PREBUILTS_GO_ROOT}/bin/go" "run" "ide_query" "$@"
