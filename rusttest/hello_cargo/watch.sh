while true; do
	inotifywait src/main.rs
	cargo +nightly fmt --all
	cargo run
	sleep 0.1
	echo
done
