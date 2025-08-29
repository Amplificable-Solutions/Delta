{
  pkgs ? import ../../nix { },
}:
let
  deltad = (pkgs.callPackage ../../. { });
in
deltad.overrideAttrs (oldAttrs: {
  patches = oldAttrs.patches or [ ] ++ [ ./broken-deltad.patch ];
})
