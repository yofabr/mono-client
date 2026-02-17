{
  description = "Nix development environment for mono-client";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
            air
            redis
            postgresql
          ];

          shellHook = ''
            echo "mono-client dev shell ready"
            echo "Run: go test ./..."
          '';
        };
      });
}
