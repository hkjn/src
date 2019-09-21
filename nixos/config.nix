{
  packageOverrides = pkgs: with pkgs; {
    tools = python37.withPackages (ps: with ps; [ btchip keepkey ckcc-protocol trezor ledgeragent ]);
  };
}
