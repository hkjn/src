// ci steps, +nightly needed for aarch64 fmt:
//   cargo build
//   cargo test
//   cargo +nightly fmt --all -- --check
//
// faster sanity checking:
//   cargo check
//
// build with optimizations:
//   cargo build --release
//
// execute program:
//   cargo run
use rand::Rng;
use std::io;

fn main() {
    println!("Please input your guess.");

    let secret_number = rand::thread_rng().gen_range(1, 101);

    println!("The secret number is: {}", secret_number);
    let mut guess = String::new();

    io::stdin()
        .read_line(&mut guess)
        .expect("Failed to read line");

    println!("You guessed: {}", guess);
}
