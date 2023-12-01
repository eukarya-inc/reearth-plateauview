use std::{fs, path::PathBuf};

use anyhow::{ensure, Context};

pub use self::compress::Format;

mod check;
mod compress;
mod list;
mod sevenzip;
mod zip;

pub struct Config {
    pub input: Vec<PathBuf>,
    pub output: PathBuf,
    pub format: Format,
}

pub fn prepare(config: Config) -> anyhow::Result<()> {
    if !config.output.exists() {
        fs::create_dir_all(&config.output)?;
    }

    ensure!(
        config.output.is_dir(),
        "{} はディレクトリではありません。",
        config.output.display(),
    );

    for input in config.input {
        eprintln!("{} を処理しています。", input.display());

        ensure!(
            input.is_dir(),
            "{} はディレクトリではありません。",
            input.display(),
        );

        ensure!(
            check::check_dir_name(
                input
                    .file_name()
                    .unwrap_or_default()
                    .to_str()
                    .unwrap_or_default()
            ),
            "フォルダ {} は正しい命名規則に従っていません。 26100_kyoto-shi_city_2022_citygml_3 のような名前にする必要があります。",
            input.file_name().unwrap_or_default().to_str().unwrap_or_default(),
        );

        compress::compress_files(&input, &config.output, &config.format)?;
    }
    Ok(())
}
