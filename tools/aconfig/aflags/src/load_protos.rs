use crate::{Flag, FlagPermission, FlagValue, ValuePickedFrom};
use aconfig_protos::ProtoFlagPermission as ProtoPermission;
use aconfig_protos::ProtoFlagState as ProtoState;
use aconfig_protos::ProtoParsedFlag;
use aconfig_protos::ProtoParsedFlags;
use anyhow::Result;
use std::{fs, str};

// TODO(b/329875578): use container field directly instead of inferring.
fn infer_container(path: &str) -> String {
    path.strip_prefix("/apex/")
        .or_else(|| path.strip_prefix('/'))
        .unwrap_or(path)
        .strip_suffix("/etc/aconfig_flags.pb")
        .unwrap_or(path)
        .to_string()
}

fn convert_parsed_flag(path: &str, flag: &ProtoParsedFlag) -> Flag {
    let namespace = flag.namespace().to_string();
    let package = flag.package().to_string();
    let name = flag.name().to_string();

    let value = match flag.state() {
        ProtoState::ENABLED => FlagValue::Enabled,
        ProtoState::DISABLED => FlagValue::Disabled,
    };

    let permission = match flag.permission() {
        ProtoPermission::READ_ONLY => FlagPermission::ReadOnly,
        ProtoPermission::READ_WRITE => FlagPermission::ReadWrite,
    };

    Flag {
        namespace,
        package,
        name,
        container: infer_container(path),
        value,
        staged_value: None,
        permission,
        value_picked_from: ValuePickedFrom::Default,
    }
}

pub(crate) fn load() -> Result<Vec<Flag>> {
    let mut result = Vec::new();

    let paths = load_aconfig_proto_paths::load_paths()?;
    for path in paths {
        let bytes = fs::read(path.clone())?;
        let parsed_flags: ProtoParsedFlags = protobuf::Message::parse_from_bytes(&bytes)?;
        for flag in parsed_flags.parsed_flag {
            // TODO(b/334954748): enforce one-container-per-flag invariant.
            result.push(convert_parsed_flag(&path, &flag));
        }
    }
    Ok(result)
}
