{
  pkgs,
  config,
  delta ? (import ../. { inherit pkgs; }),
}:
rec {
  start-delta = pkgs.writeShellScriptBin "start-delta" ''
    # rely on environment to provide deltad
    export PATH=${pkgs.test-env}/bin:$PATH
    ${../scripts/start-delta} ${config.delta-config} ${config.dotenv} $@
  '';
  start-geth = pkgs.writeShellScriptBin "start-geth" ''
    export PATH=${pkgs.test-env}/bin:${pkgs.go-ethereum}/bin:$PATH
    source ${config.dotenv}
    ${../scripts/start-geth} ${config.geth-genesis} $@
  '';
  start-scripts = pkgs.symlinkJoin {
    name = "start-scripts";
    paths = [
      start-delta
      start-geth
    ];
  };
}
