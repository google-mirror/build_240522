_all_versions = ["OPR1", "OPD1", "OPD2", "OPM1", "OPM2", "PPR1", "PPD1", "PPD2", "PPM1", "PPM2", "QPR1"] + [
    version + subversion
    for version in ["Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"]
    for subversion in ["P1A", "P1B", "P2A", "P2B", "D1A", "D1B", "D2A", "D2B", "Q1A", "Q1B", "Q2A", "Q2B", "Q3A", "Q3B"]
]

variables_to_export_to_make = {
    "ALL_VERSIONS": _all_versions,
}
