use std::{
    fs, io,
    path::{Path, PathBuf},
    sync::mpsc::channel,
};

use anyhow::Context;
use rayon::prelude::*;

#[derive(Debug)]
pub struct Files {
    pub dirs: Vec<(String, PathBuf)>,
    pub misc: Vec<PathBuf>,
}

pub fn list_files(dir_path: &Path) -> io::Result<Files> {
    let mut dirs = vec![];
    let mut misc = Vec::<PathBuf>::new();

    for path in read_dir_with_sorted(dir_path)? {
        let name = path
            .file_name()
            .unwrap_or_default()
            .to_str()
            .unwrap_or_default();

        match name {
            "codelists" | "metadata" | "schemas" => {
                dirs.push((name.to_string(), path));
            }
            "udx" => {
                for path in read_dir_with_sorted(&path)? {
                    let name = path
                        .file_name()
                        .unwrap_or_default()
                        .to_str()
                        .unwrap_or_default();

                    dirs.push((name.to_string(), path));
                }
            }
            _ => {
                misc.push(path);
            }
        }
    }

    Ok(Files { dirs, misc })
}

const SKIPED_FILES: &[&str] = &[".DS_Store", "Thumbs.db", "__MACOSX"];

pub fn copy_files(files: &Files, output_dir: &Path) -> anyhow::Result<Vec<PathBuf>> {
    let opts: fs_extra::dir::CopyOptions = fs_extra::dir::CopyOptions {
        overwrite: true,
        ..Default::default()
    };

    let (sender, receiver) = channel();

    files
        .dirs
        .par_iter()
        .try_for_each_with(sender, |s, f| -> anyhow::Result<()> {
            if SKIPED_FILES.contains(&f.0.as_str()) {
                return Ok(());
            }

            let files = if f.1.is_dir() {
                read_dir_with_sorted(&f.1).with_context(|| {
                    format!(
                        "{} ディレクトリのファイル一覧の取得に失敗しました。",
                        f.1.display()
                    )
                })?
            } else {
                vec![f.1.clone()]
            };

            let output_dir = output_dir.join(&f.0);
            fs::create_dir_all(&output_dir)?;
            fs_extra::copy_items(&files, &output_dir, &opts).with_context(|| {
                format!(
                    "{} ディレクトリのファイルのコピーに失敗しました。",
                    f.1.display()
                )
            })?;

            s.send(output_dir.clone())?;
            Ok(())
        })?;

    let mut copied = receiver.iter().collect::<Vec<_>>();
    copied.sort();

    let output_dir = output_dir.join("misc");
    fs::create_dir_all(&output_dir)?;
    fs_extra::copy_items(&files.misc, &output_dir, &opts)?;
    copied.push(output_dir);

    Ok(copied)
}

fn read_dir_with_sorted(dir_path: &Path) -> io::Result<Vec<PathBuf>> {
    let mut files = fs::read_dir(dir_path)?.collect::<io::Result<Vec<_>>>()?;

    files.sort_by(|a, b| {
        a.file_name()
            .cmp(&b.file_name())
            .then(a.path().cmp(&b.path()))
    });

    Ok(files.into_iter().map(|entry| entry.path()).collect())
}

#[cfg(test)]
mod tests {
    use tempdir::TempDir;

    use super::*;

    #[test]
    fn test_list_files() -> anyhow::Result<()> {
        let tmpdir = TempDir::new("pvt-test")?;
        let root = tmpdir.path().join("26100_kyoto-shi_city_2022_citygml_3");

        create_dummy_files(&root)?;
        let files = list_files(&root).unwrap();

        assert_eq!(
            files.dirs,
            vec![
                ("codelists".to_string(), root.join("codelists")),
                ("metadata".to_string(), root.join("metadata")),
                ("schemas".to_string(), root.join("schemas")),
                ("bldg".to_string(), root.join("udx").join("bldg")),
                ("tran".to_string(), root.join("udx").join("tran"))
            ],
        );

        assert_eq!(
            files.misc,
            vec![
                root.join("26100_indexmap.pdf"),
                root.join("README.md"),
                root.join("specification"),
            ]
        );

        Ok(())
    }

    #[test]
    fn test_copy_files() -> anyhow::Result<()> {
        let tmpdir = TempDir::new("pvt-test")?;
        let root = tmpdir.path().join("26100_kyoto-shi_city_2022_citygml_3");

        create_dummy_files(&root)?;
        let entries = list_files(&root).unwrap();

        let output_dir = tmpdir.path().join("output");
        copy_files(&entries, &output_dir).unwrap();

        let result = read_dir_with_sorted(&output_dir)?;
        assert_eq!(
            result,
            vec![
                output_dir.join("bldg"),
                output_dir.join("codelists"),
                output_dir.join("metadata"),
                output_dir.join("misc"),
                output_dir.join("schemas"),
                output_dir.join("tran"),
            ]
        );

        let result = read_dir_with_sorted(output_dir.join("codelists").as_path())?;
        assert_eq!(result, vec![output_dir.join("codelists").join("hoge.gml")]);

        let result = read_dir_with_sorted(output_dir.join("schemas").as_path())?;
        assert_eq!(result, vec![output_dir.join("schemas").join("iur")]);

        let result = read_dir_with_sorted(output_dir.join("metadata").as_path())?;
        assert_eq!(result, vec![output_dir.join("metadata").join("foo.gml")]);

        let result = read_dir_with_sorted(output_dir.join("bldg").as_path())?;
        assert_eq!(result, vec![output_dir.join("bldg").join("bar.gml")]);

        let result = read_dir_with_sorted(output_dir.join("tran").as_path())?;
        assert_eq!(result, vec![output_dir.join("tran").join("fuga.gml")]);

        let result = read_dir_with_sorted(output_dir.join("misc").as_path())?;
        assert_eq!(
            result,
            vec![
                output_dir.join("misc").join("26100_indexmap.pdf"),
                output_dir.join("misc").join("README.md"),
                output_dir.join("misc").join("specification"),
            ]
        );

        Ok(())
    }

    fn create_dummy_files(root: &Path) -> io::Result<()> {
        fs::create_dir_all(root.join("codelists"))?;
        fs::create_dir_all(root.join("metadata"))?;
        fs::create_dir_all(root.join("schemas"))?;
        fs::create_dir_all(root.join("specification"))?;
        fs::create_dir_all(root.join("udx").join("bldg"))?;
        fs::create_dir_all(root.join("udx").join("tran"))?;
        fs::write(root.join("codelists").join("hoge.gml"), "dummy")?;
        fs::write(root.join("metadata").join("foo.gml"), "dummy")?;
        fs::write(root.join("schemas").join("iur"), "dummy")?;
        fs::write(root.join("specification").join("iur"), "dummy")?;
        fs::write(root.join("udx").join("bldg").join("bar.gml"), "dummy")?;
        fs::write(root.join("udx").join("tran").join("fuga.gml"), "dummy")?;
        fs::write(root.join("26100_indexmap.pdf"), "dummy")?;
        fs::write(root.join("README.md"), "dummy")?;
        Ok(())
    }
}
