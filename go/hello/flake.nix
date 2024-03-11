{
  description = "Devshell and Building {{ .ProjectName }}";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        module = pkgs.buildGoModule {
          pname = "{{ .ProjectName }}";
          version = self.shortRev or "dirty";
          src = ./.;

          vendorHash = "";
        };
      in
      {
        packages.default = module;
        packages.{{ .ProjectName }} = module;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
      }
    );
}
