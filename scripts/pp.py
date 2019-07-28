#! /usr/bin/env nix-shell
#! nix-shell -i python -p python pythonPackages.prettytable
#
# Sample script showing how to use nix-shell to install python
# runtime and parser.
#
import prettytable
t = prettytable.PrettyTable(["N", "N^2"])
for n in range(1, 10): t.add_row([n, n * n])
print t

