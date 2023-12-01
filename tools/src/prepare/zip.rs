use std::{
    fs::File,
    io::{copy, BufWriter, Write},
    path::Path,
};

use anyhow::Context as _;
use walkdir::WalkDir;
use zip::ZipWriter;

pub fn zip(paths: &[impl AsRef<Path>], zip_path: &Path) -> anyhow::Result<()> {
    let bufw = BufWriter::new(File::create(zip_path)?);
    let mut zw = ZipWriter::new(bufw);

    for p in paths {
        if p.as_ref().is_dir() {
            for entry in WalkDir::new(p) {
                let entry = entry.context("ファイルを取得できませんでした。")?;
                let path = entry.path();
                let path_str = path.to_str().with_context(|| {
                    format!("{} は正しいパスではありません。", entry.path().display())
                })?;

                if path.is_dir() {
                    zw.add_directory(path_str, Default::default())
                        .with_context(|| {
                            format!(
                                "{} を圧縮ファイルに追加できませんでした。",
                                entry.path().display()
                            )
                        })?
                } else {
                    zw.start_file(path_str, Default::default())
                        .with_context(|| {
                            format!(
                                "{} を圧縮ファイルに追加できませんでした。",
                                entry.path().display()
                            )
                        })?;
                    copy(
                        &mut File::open(path).context("ファイルを開くことができませんでした。")?,
                        &mut zw,
                    )
                    .with_context(|| {
                        format!(
                            "{} を圧縮ファイルに追加できませんでした。",
                            entry.path().display()
                        )
                    })?;
                }
            }
        } else {
            zw.start_file(
                p.as_ref().to_str().with_context(|| {
                    format!("{} は正しいパスではありません。", p.as_ref().display())
                })?,
                Default::default(),
            )
            .with_context(|| {
                format!(
                    "{} を圧縮ファイルに追加できませんでした。",
                    p.as_ref().display()
                )
            })?;

            copy(
                &mut File::open(p.as_ref()).context("ファイルを開くことができませんでした。")?,
                &mut zw,
            )
            .with_context(|| {
                format!(
                    "{} を圧縮ファイルに追加できませんでした。",
                    p.as_ref().display()
                )
            })?;
        }
    }

    zw.flush()?;
    zw.finish()?;
    Ok(())
}
