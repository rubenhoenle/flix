{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, treefmt-nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        treefmtEval = treefmt-nix.lib.evalModule pkgs {
          projectRootFile = "flake.nix";
          programs = {
            nixpkgs-fmt.enable = true;
            prettier = {
              enable = true;
              includes = [ "*.md" "*.yaml" "*.yml" ];
            };
            gofmt.enable = true;
          };
        };

        flix-backend = pkgs.buildGoModule {
          name = "flix-backend";
          version = "0.0.1";
          vendorHash = "sha256-Ef2XLxGq8TO3WVh9EvLE30Is2CBwH4pqXxkq1tcuR0Q=";
          src = ./backend;
        };

        containerImage = pkgs.dockerTools.buildLayeredImage {
          name = "ghcr.io/rubenhoenle/flix";
          tag = "unstable";
          config = {
            Entrypoint = [ "${flix-backend}/bin/flix-backend" ];
          };
        };
      in
      {
        formatter = treefmtEval.config.build.wrapper;
        checks.formatter = treefmtEval.config.build.check self;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];
        };

        packages = flake-utils.lib.flattenTree {
          default = flix-backend;
          containerimage = containerImage;
        };
      }
    );
}
