load("//build/bazel/rules/apis:cc_api_contribution.bzl", "cc_api_contribution", "cc_api_headers")
load("//build/bazel/rules/apis:api_surface.bzl", "api_surface")

# api contributions from system
cc_api_contribution(
  name="libhwui",
  api=":libhwui.map.txt",
  hdrs=[":libhwui_headers"],
)

cc_api_headers(
  name="libhwui_headers",
  root="build/bazel/tests/apis/include",
  hdrs=glob(["include/*.h"]),
)

api_surface(
  name="public_api_system_v1",
  surface_name="public_api",
  version=1,
  api_domain="system",
  cc_libraries=[":libhwui"]
)

# api contributions from bionic
cc_api_contribution(
  name="libc",
  api=":libc.map.txt",
  hdrs=[], # TODO: Use headers
)

api_surface(
  name="public_api_bionic_v1",
  surface_name="public_api",
  version=1,
  api_domain="com.android.bionic",
  cc_libraries=[":libc"],
)

# top-level phony target for cdk
# contains surfaces from different api domains
# contains select per api_domain
filegroup(
  name="cdk",
  srcs=select({
    ":system": [":public_api_system_v1"],
    ":bionic": [":public_api_bionic_v1"],
  }),
)

config_setting(
  name="system",
  values = {"define": "api_domain=system"}
)

config_setting(
  name="bionic",
  values = {"define": "api_domain=bionic"}
)

