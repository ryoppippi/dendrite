{
  description = "Generate Nix default.nix files from YAML tool definitions";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      perSystem =
        { pkgs, system, ... }:
        let
          lib = pkgs.lib;

          goMod = builtins.readFile ./go.mod;
          goLine = lib.findFirst (line: lib.hasPrefix "go " line) (
            throw "go.mod does not contain a go directive"
          ) (lib.splitString "\n" goMod);
          goVersion = lib.removePrefix "go " goLine;
          goVersionParts = lib.splitVersion goVersion;
          goAttr = "go_${builtins.elemAt goVersionParts 0}_${builtins.elemAt goVersionParts 1}";
          go =
            if pkgs ? ${goAttr} then
              pkgs.${goAttr}
            else
              throw "nixpkgs does not provide ${goAttr} for go ${goVersion}";

          versionLine = lib.findFirst (line: lib.hasInfix "const Version = " line) (
            throw "version.go does not contain Version"
          ) (lib.splitString "\n" (builtins.readFile ./version.go));
          version = builtins.elemAt (builtins.match ''.*"([^"]+)".*'' versionLine) 0;

          nixPath = lib.makeBinPath [ pkgs.nix ];
        in
        {
          packages.default = (pkgs.buildGoModule.override { inherit go; }) {
            pname = "dendrite";
            inherit version;

            src = self;
            subPackages = [ "cmd/dendrite" ];
            vendorHash = "sha256-g+yaVIx4jxpAQ/+WrGKxhVeliYx7nLQe/zsGpxV4Fn4=";

            nativeBuildInputs = [ pkgs.makeWrapper ];

            postFixup = ''
              wrapProgram $out/bin/dendrite --prefix PATH : ${nixPath}
            '';
          };

          apps.default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/dendrite";
          };

          devShells.default = pkgs.mkShell {
            packages = [
              go
              pkgs.golangci-lint
              pkgs.nix
            ];
          };
        };
    };
}
