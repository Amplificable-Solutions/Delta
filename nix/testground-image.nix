{
  dockerTools,
  runCommandLocal,
  delta-matrix,
  benchmark-testcase,
}:
let
  patched-deltad = delta-matrix.deltad.overrideAttrs (oldAttrs: {
    patches = oldAttrs.patches or [ ] ++ [ ./testground-deltad.patch ];
  });
in
let
  tmpDir = runCommandLocal "tmp" { } ''
    mkdir -p $out/tmp/
  '';
in
dockerTools.buildLayeredImage {
  name = "delta-testground";
  created = "now";
  contents = [
    benchmark-testcase
    patched-deltad
    tmpDir
  ];
  config = {
    Expose = [
      9090
      26657
      26656
      1317
      26658
      26660
      26659
      30000
    ];
    Cmd = [ "/bin/stateless-testcase" ];
    Env = [ "PYTHONUNBUFFERED=1" ];
  };
}
