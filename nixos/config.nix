{
  packageOverrides = pkgs: with pkgs; {
    tools = python37.withPackages (ps: with ps; [ btchip keepkey ckcc-protocol trezor ]);
# xx: gives attribute 'withPackages' missing
#    wasabi = dotnetPackages.withPackages (ps: with ps; [ wasabiwallet ]);
  };
}
