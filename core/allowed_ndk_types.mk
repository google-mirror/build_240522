# Determines the types of NDK modules the current module is allowed to link to.
# Input variables:
#   ???
# Output variables:
#   my_ndk_stl_family: Family of the NDK STL.
#   my_ndk_stl_link_type: STL link type, static or shared.
#   my_allowed_ndk_types: Types of NDK modules that may be linked.

ifdef LOCAL_SDK_VERSION
    ifeq ($(LOCAL_NDK_STL_VARIANT),)
        my_ndk_stl_family := system
        my_ndk_stl_link_type := shared
    else ifeq ($(LOCAL_NDK_STL_VARIANT),system)
        my_ndk_stl_family := system
        my_ndk_stl_link_type := shared
    else ifeq ($(LOCAL_NDK_STL_VARIANT),c++_shared)
        my_ndk_stl_family := libc++
        my_ndk_stl_link_type := shared
    else ifeq ($(LOCAL_NDK_STL_VARIANT),c++_static)
        my_ndk_stl_family := libc++
        my_ndk_stl_link_type := static
    else ifeq ($(LOCAL_NDK_STL_VARIANT),gnustl_static)
        my_ndk_stl_family := gnustl
        my_ndk_stl_link_type := static
    else ifeq ($(LOCAL_NDK_STL_VARIANT),stlport_shared)
        my_ndk_stl_family := stlport
        my_ndk_stl_link_type := shared
    else ifeq ($(LOCAL_NDK_STL_VARIANT),stlport_static)
        my_ndk_stl_family := stlport
        my_ndk_stl_link_type := static
    else ifeq ($(LOCAL_NDK_STL_VARIANT),none)
        my_ndk_stl_family := none
        my_ndk_stl_link_type := none
    else
        $(error $(LOCAL_MODULE): invalid LOCAL_NDK_STL_VARIANT: $(LOCAL_NDK_STL_VARIANT))
    endif
    my_allowed_ndk_types :=

    # The system STL is only the C++ ABI layer, so it's compatible with any STL.
    my_allowed_ndk_types += native:ndk:system

    # Libaries that don't use the STL can be linked to anything.
    my_allowed_ndk_types += native:ndk:none

    # And it's okay to link your own STL type. Strictly speaking there are more
    # restrictions depending on static vs shared STL, but that will be a follow
    # up patch.
    my_allowed_ndk_types += native:ndk:$(my_ndk_stl_family)

    # For an app package, it's actually okay to depend on any set of STLs. If
    # any of the individual libraries depend on each other they've already been
    # checked for consistency, and if they don't they'll be kept isolated by
    # RTLD_LOCAL anyway.
    ifeq ($(LOCAL_MODULE_CLASS),APPS)
        my_allowed_ndk_types += \
            native:ndk:gnustl native:ndk:libc++ native:ndk:stlport

    ifeq ($(LOCAL_MODULE_CLASS),STATIC_LIBRARIES)
        # Since static libraries are not actually linked, the decision between
        # static/shared isn't actually meaningful here. Only the family of the
        # STL matters.
        my_allowed_ndk_types += native:ndk:$(my_ndk_stl_family):shared
        my_allowed_ndk_types += native:ndk:$(my_ndk_stl_family):static
    else ifeq ($(my_ndk_stl_link_type),shared)
        # Modules linked a shared STL can only use another shared STL.
        my_allowed_ndk_types += native:ndk:$(my_ndk_stl_family):shared
    endif
    # Else we are a non-static library that uses a static STL, and are
    # incompatible with all other shared libraries that use an STL.
else
    my_allowed_ndk_types := native:ndk:none native:ndk:system
endif
