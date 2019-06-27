: red thread of practical examples:
:   - ncoleman tools to generate p2pkh, p2sh-wrapped-p2wpkh and p2wpkh addresses
:     - or generating seed phrase using dice and coin flips: https://github.com/taelfrinn/Bip39-diceware
:   - EC math to from WIF-encoded privkey generate pubkey and hashing it to p2wpkh
:   - computing seed words from entropy
:   - signing tx with privkey
:   - looking up transactions with bitcoin-cli getrawtransaction <txid>
:   - looking up blocks with getblockhash <height> + getblock <blockhash>
:   - types of wallets: software, hardware, full nodes, "light clients" vs remote controls
:   - downloading bitcoin core, verifying signature, syncing
:       blockchain (shared server / tmux)
:   - 21.hkjn.me access from client: download script from here: <short link / github raw link>
:   - 21.hkjn.me steps on server: download script from here: <short link / github raw link>
:     - install bitcoin core, start syncing
:     - check progress with getblockchaininfo
:     - check connection info and peers with getpeers, getnetworkinfo
:     - after IBD, configure Tor, and listen on .onion service
