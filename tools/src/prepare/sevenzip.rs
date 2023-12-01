use std::{fs::File, io::BufWriter, path::Path};

use sevenz_rust::SevenZWriter;

pub fn sevenzip(paths: &[impl AsRef<Path>], zip_path: &Path) -> anyhow::Result<()> {
    let bufw = BufWriter::new(File::create(zip_path)?);
    let mut z = SevenZWriter::new(bufw)?;

    for p in paths {
        z.push_source_path(p, |_| true)?;
    }

    z.finish()?;
    Ok(())
}
