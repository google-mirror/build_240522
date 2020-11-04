load(":product_config.star", "prodconf")

base_product = prodconf(
	"base_product",
    [],
    PRODUCT_PACKAGES=[
        "fs_config_dirs_product",
        "fs_config_files_product",
        "group_product",
        "ModuleMetadata",
        "passwd_product",
        "product_compatibility_matrix.xml",
        "product_manifest.xml",
        "selinux_policy_product",
    ]
)
