#!/usr/bin/env bash

echo "Importing public keys.."
for key in /etc/keys/*.asc;
	do gpg --import < $key
done

OWNER_TRUST=$(cat <<EOL
425BF55E014AF99C3BA6A6E8D85FAD19F4971232:6:
D30CCB024F9A403708846F135CE7AD68FA100CAC:6:
FB5C74D329E598A6693D2498E82168C7C4DCEC4B:6:
EOL
)
echo "${OWNER_TRUST}" | gpg --import-ownertrust
