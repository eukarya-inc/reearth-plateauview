use std::path::Path;

use anyhow::Context;
use rayon::iter::IntoParallelIterator;

use crate::prepare::{list::list_files, zip::zip};

use super::{list::copy_files, sevenzip::sevenzip};
use rayon::prelude::*;

pub enum Format {
    Auto,
    None,
    Zip,
    #[allow(dead_code)]
    SevenZip,
}

pub fn compress_files(
    input: &Path,
    output: Option<impl AsRef<Path>>,
    format: &Format,
) -> anyhow::Result<()> {
    let output = if let Some(output) = output {
        output.as_ref().to_path_buf()
    } else {
        input
            .parent()
            .with_context(|| {
                format!(
                    "入力ファイルの親ディレクトリが取得できませんでした。: {}",
                    input.display()
                )
            })?
            .to_path_buf()
    };

    let format = if let Format::Auto = format {
        // TODO
        &Format::Zip
    } else {
        format
    };

    let prefix = input
        .file_name()
        .unwrap_or_default()
        .to_str()
        .unwrap_or_default();
    let files = list_files(input)?;
    let output = output.join(format!("{}_files", prefix));

    println!("{} のコピーを開始します。", prefix);
    let copied_files = copy_files(&files, &output)?;
    println!("{} のコピーが完了しました。", prefix);

    if let Format::None = format {
        return Ok(());
    }

    println!("{} の圧縮を開始します。", prefix);

    copied_files
        .into_par_iter()
        .try_for_each(|path| -> anyhow::Result<()> {
            let name = path
                .file_name()
                .unwrap_or_default()
                .to_str()
                .unwrap_or_default();

            println!("{}/{} を圧縮しています。", prefix, name);

            match format {
                Format::Zip => {
                    let zip_path = output.join(format!("{}_{}.zip", prefix, name));
                    zip(&path, &zip_path)?;
                }
                Format::SevenZip => {
                    let zip_path = input.join(format!("{}_{}.7z", prefix, name));
                    sevenzip(&path, &zip_path)?;
                }
                _ => return Ok(()),
            }

            println!("{}/{} の圧縮が完了しました。", prefix, name);
            Ok(())
        })?;

    println!("{} の圧縮が完了しました。", prefix);
    Ok(())
}
