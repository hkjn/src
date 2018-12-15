# bitcoin

This doc has some notes on interesting things to study / contribute to within Bitcoin.

## core support for hardware wallets

The following issue and [achow's work on HWI](https://gist.github.com/achow101/a9cf757d45df56753fae9d65db4d6e1d)
intend to make it easier to run Bitcoin Core with hardware wallets:

- https://github.com/bitcoin/bitcoin/issues/14145

## output descriptors

Study the new output descriptor DSL added in Bitcoin Core v0.17.0:

- https://github.com/bitcoin/bitcoin/blob/master/doc/descriptors.md

This is a more general way to specify what types of output scripts are used for certain keys.

## hardware wallet clients

- Coldcard: https://github.com/Coldcard/ckcc-protocol
- Keepkey: https://github.com/keepkey/python-keepkey

Via achow101's HWI:

```
pip3 install hidapi # HID API needed in general
pip3 install trezor[hidapi] # Trezor One
pip3 install btchip-python # Ledger Nano S
pip3 install keepkey # KeepKey
pip3 install ckcc-protocol # Coldcard
pip3 install pyaes # For digitalbitbox
```
