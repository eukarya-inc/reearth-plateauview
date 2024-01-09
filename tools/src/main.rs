mod args;
mod migrateuc;
mod prepare;
use args::{Cli, Commands};

use clap::Parser;

fn main() {
    let cli = Cli::parse();

    if let Err(err) = match cli.command {
        Commands::Prepare {
            format,
            targets,
            output,
        } => prepare::prepare(prepare::Config {
            input: targets,
            output,
            format: format.into(),
        }),
        Commands::MigrateUC { list_path, output } => {
            migrateuc::migrateuc(migrateuc::Config { list_path, output })
        }
    } {
        eprintln!("{}", err);
        std::process::exit(1);
    }
}
