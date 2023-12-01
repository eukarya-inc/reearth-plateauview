use std::path::Path;

use crate::prepare::{list::list_files, zip::zip};

use super::sevenzip::sevenzip;

pub enum Format {
    Auto,
    Zip,
    SevenZip,
}

pub fn compress_files(input: &Path, output: &Path, format: &Format) -> anyhow::Result<()> {
    let format = if let Format::Auto = format {
        // TODO
        &Format::Zip
    } else {
        format
    };

    let entires = list_files(input)?;

    if !output.exists() {
        std::fs::create_dir_all(&output)?;
    }

    for entry in entires {
        match format {
            Format::Auto => unreachable!(),
            Format::Zip => {
                let zip_path = output.join(format!("{}.zip", entry.name));
                zip(&entry.files, &zip_path)?;
            }
            Format::SevenZip => {
                let zip_path = input.join(format!("{}.7z", entry.name));
                sevenzip(&entry.files, &zip_path)?;
            }
        }
    }

    Ok(())
}
