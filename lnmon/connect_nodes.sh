set -euo pipefail

lightning-cli listnodes > nodes.json
while read -r line; do
	echo "$line"
	$line
done <<< "$(python .lightning/parse_nodes.py)"
