def switch_txid_endianness(txid):
    """Switch endianness (big- vs little-endian) of txid."""

    return ''.join(tx[i:i+2] for i in range(0, len(tx), 2)[::-1])
