use std::{fs::create_dir_all, path::PathBuf};

pub struct Config {
    pub list_path: PathBuf,
    pub output: Option<PathBuf>,
}

pub fn migrateuc(config: Config) -> anyhow::Result<()> {
    let output = config.output.unwrap_or_else(|| PathBuf::from("uc"));
    create_dir_all(output)?;

    todo!()
}
