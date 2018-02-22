set -euo pipefail

echo "----------------------------"
echo "Listing all known outputs:"
echo "----------------------------"
sqlite3 .lightning/lightningd.sqlite3 <<EOF
.headers on
.mode csv
SELECT
	HEX(prev_out_tx) AS prev_out_tx_hex,
	prev_out_index,
	value,
	-- Via wallet_output_type defined at https://github.com/ElementsProject/lightning/blob/master/wallet/wallet.h#L47
	CASE type
		WHEN '0' THEN 'p2sh_wpkh'
		WHEN '1' THEN 'to_local'
		WHEN '3' THEN 'htlc_offer'
		WHEN '4' THEN 'htlc_offer'
		WHEN '4' THEN 'htlc_recv'
		WHEN '5' THEN 'our_change'
		WHEN '6' THEN 'p2wpkh'
		ELSE 'unknown_state'
	END AS type_desc,
	-- Via output_status defined at https://github.com/ElementsProject/lightning/blob/master/wallet/wallet.h#L34
       CASE status
		WHEN '0' THEN 'output_status_available'
		WHEN '1' THEN 'output_state_reserved'
		WHEN '2' THEN 'output_state_spent'
		WHEN '255' THEN 'output_state_any'
		ELSE 'unknown_status'
	END AS status_desc
FROM outputs;
EOF

echo "----------------------------"
echo "Listing available outputs:"
echo "----------------------------"
sqlite3 .lightning/lightningd.sqlite3 <<EOF
.headers on
.mode csv
SELECT
	HEX(prev_out_tx) AS prev_out_tx_hex,
	prev_out_index,
	value,
	-- Via wallet_output_type defined at https://github.com/ElementsProject/lightning/blob/master/wallet/wallet.h#L47
	CASE type
		WHEN '0' THEN 'p2sh_wpkh'
		WHEN '1' THEN 'to_local'
		WHEN '3' THEN 'htlc_offer'
		WHEN '4' THEN 'htlc_offer'
		WHEN '4' THEN 'htlc_recv'
		WHEN '5' THEN 'our_change'
		WHEN '6' THEN 'p2wpkh'
		ELSE 'unknown_state'
	END AS type_desc
FROM outputs
WHERE
  -- output_status_available
  status=0;
EOF

echo "----------------------------"
echo "Summing up spendable UTXOs.."
echo "----------------------------"
RESULT=$(sqlite3 .lightning/lightningd.sqlite3 <<EOF
.headers on
.mode csv
SELECT
	SUM(value) AS sum,
	-- Via wallet_output_type defined at https://github.com/ElementsProject/lightning/blob/master/wallet/wallet.h#L47
	CASE type
		WHEN '0' THEN 'p2sh_wpkh'
		WHEN '1' THEN 'to_local'
		WHEN '3' THEN 'htlc_offer'
		WHEN '4' THEN 'htlc_offer'
		WHEN '4' THEN 'htlc_recv'
		WHEN '5' THEN 'our_change'
		WHEN '6' THEN 'p2wpkh'
		ELSE 'unknown_state'
	END AS type_desc
FROM outputs
WHERE
  -- output_status_available
  status=0
GROUP BY type_desc;
EOF
)
echo "${RESULT}"
# TODO: could reverse to get little-endian version of txids, i.e. compatible with bitcoin-cli:
# d='${RESULT}'
# txids=[x.split(',')[0] for x in d.split()[1:]]
# for txid in txids:
#   print(''.join(txid[i:i+2] for i in range(0, len(txid), 2)[::-1]))
