use std::path::PathBuf;

use clap::{Parser, Subcommand, ValueEnum};

#[derive(Debug, Parser)]
#[command(author, version, about, long_about = None)]
#[command(propagate_version = true)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Debug, Subcommand)]
pub enum Commands {
    /// CMSにアップロードできるように、指定されたPLATEAUデータのCityGMLファイルを地物単位に分割・圧縮します。
    Prepare {
        /// 圧縮形式を指定します。
        #[clap(short, long, default_value = "auto")]
        format: Format,
        /// 地物単位に分割・圧縮するPLATEAUデータのCityGMLファイルが格納されたフォルダへのパスを指定します。
        targets: Vec<PathBuf>,
        /// 地物単位に分割・圧縮した結果を格納するフォルダへのパスを指定します。無指定の場合は入力元と同じフォルダに作成します。
        #[clap(short, long)]
        output: Option<PathBuf>,
    },
}

#[derive(Debug, Clone, ValueEnum)]
pub enum Format {
    #[clap(name = "auto")]
    Auto,
    #[clap(name = "none")]
    None,
    #[clap(name = "zip")]
    Zip,
    // #[clap(name = "7z")]
    // SevenZip,
}

impl Into<crate::prepare::Format> for Format {
    fn into(self) -> crate::prepare::Format {
        match self {
            Format::Auto => crate::prepare::Format::Auto,
            Format::None => crate::prepare::Format::None,
            Format::Zip => crate::prepare::Format::Zip,
            // Format::SevenZip => crate::prepare::Format::SevenZip,
        }
    }
}
