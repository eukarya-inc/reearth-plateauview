use std::{
    fs, io,
    path::{Path, PathBuf},
};

#[derive(Debug)]
pub struct Entry {
    pub name: String,
    pub files: Vec<PathBuf>,
}

pub fn list_files(dir_path: &Path) -> io::Result<Vec<Entry>> {
    let mut entries = vec![];
    let mut misc = Vec::<PathBuf>::new();

    let files = fs::read_dir(dir_path)?.collect::<io::Result<Vec<_>>>()?;

    for entry in files {
        let path = entry.path();

        let name = path
            .file_name()
            .unwrap_or_default()
            .to_str()
            .unwrap_or_default();

        match name {
            "codelists" | "metadata" | "schemas" => {
                entries.push(Entry {
                    name: name.to_string(),
                    files: vec![path],
                });
            }
            "udx" => {
                let files = fs::read_dir(&path)?.collect::<io::Result<Vec<_>>>()?;

                for entry in files {
                    let path = entry.path();

                    let name = path
                        .file_name()
                        .unwrap_or_default()
                        .to_str()
                        .unwrap_or_default();

                    entries.push(Entry {
                        name: name.to_string(),
                        files: vec![path],
                    });
                }
            }
            _ => {
                misc.push(path);
            }
        }
    }

    if !misc.is_empty() {
        entries.push(Entry {
            name: "misc".to_string(),
            files: misc,
        });
    }

    Ok(entries)
}

// #[cfg(test)]
// mod tests {
//     use super::*;

//     #[test]
//     fn test_list_files() {
//         let tempdir = tempfile::tempdir().unwrap();
//         let dir_path = Path::new("tests/data/prepare/list_files");

//         let entries = list_files(dir_path).unwrap();

//         assert_eq!(entries.len(), 3);

//         assert_eq!(entries[0].name, "codelists");
//         assert_eq!(entries[0].files.len(), 1);
//         assert_eq!(entries[0].files[0], dir_path.join("codelists"));

//         assert_eq!(entries[1].name, "metadata");
//         assert_eq!(entries[1].files.len(), 1);
//         assert_eq!(entries[1].files[0], dir_path.join("metadata"));

//         assert_eq!(entries[2].name, "misc");
//         assert_eq!(entries[2].files.len(), 2);
//         assert_eq!(entries[2].files[0], dir_path.join("schemas"));
//         assert_eq!(entries[2].files[1], dir_path.join("udx"));
//     }
// }
