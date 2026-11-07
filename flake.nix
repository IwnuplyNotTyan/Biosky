{
  description = "Biosky ~ Bluesky bio rotator";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    nixpkgs-darwinx86.url = "github:NixOS/nixpkgs/nixpkgs-26.05-darwin";
  };

  outputs = { self, nixpkgs, nixpkgs-darwinx86, ... }:
    let
      systems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
      pkgsFor = system:
        if system == "x86_64-darwin"
        then nixpkgs-darwinx86.legacyPackages.${system}
        else nixpkgs.legacyPackages.${system};
      version = "0.1.0";
      mkBiosky = pkgs: pkgs.buildGo127Module {
        pname = "biosky";
        inherit version;
        src = self;

        vendorHash = null;

        ldflags = [
          "-s"
          "-w"
          "-X"
          "main.version=${version}"
        ];

        trimpath = true;

        meta = {
          description = "Bluesky bio rotator";
          mainProgram = "biosky";
        };
      };
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = mkBiosky pkgs;
        });
      devShells = forAllSystems (system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go_1_27
              gotools
              golangci-lint
            ];
          };
        });
    };
}
