#
# Print LN nodes that are responsive on specified addr / port.
#
# Usage:
# $ sh connect_nodes.sh > connect.sh # generate list of lightning-cli connect commands
# $ sh connect.sh                    # run the lightning-cli connect commands
#
set -eu

lightning-cli listnodes > nodes.json
python parse_nodes.py
