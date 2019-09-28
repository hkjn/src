{
  packageOverrides = pkgs: with pkgs; {
    tools = python37.withPackages (ps: with ps; [ btchip ckcc-protocol keepkey trezor keepkey_agent ledger_agent trezor_agent]);
  };
}
