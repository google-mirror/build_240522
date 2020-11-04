load(":product_config.star", "prodconf")

base_system_ext = prodconf(
    "base_system_ext",
    [],
    PRODUCT_PACKAGES=[
        "fs_config_dirs_system_ext",
        "fs_config_files_system_ext",
        "group_system_ext",
        "passwd_system_ext",
        "selinux_policy_system_ext",
        "system_ext_manifest.xml",
    ],
)
