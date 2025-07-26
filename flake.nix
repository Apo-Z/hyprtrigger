{
  description = "A simple Go package";
  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-25.05";
  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          hyprtrigger = pkgs.buildGoModule rec {
            pname = "hyprtrigger";
            version = "0.0.4";
            srcHash = "sha256-jwHjsyTuTQk0Kk6uI5puBOa9iNEGQt9vSwD/BfaCI2Y=";
            src = pkgs.fetchFromGitHub {
              owner = "Apo-Z";
              repo = "hyprtrigger";
              rev = "v${version}";
              hash = srcHash;
            };
            vendorHash = null;
            ldflags = [
              "-X hyprtrigger/cmd/hyprtrigger.version=v${version}"
              "-X hyprtrigger/cmd/hyprtrigger.commit=${srcHash}"
            ];
            meta = with pkgs.lib; {
              description = "A trigger system for Hyprland";
              homepage = "https://github.com/Apo-Z/hyprtrigger";
              license = licenses.mit;
              maintainers = [ ];
              platforms = platforms.unix;
            };
          };
        });
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              go-tools
              git
              gnumake
            ];
            # Variables d'environnement pour le développement
            shellHook = ''
              echo "Environnement de développement hyprtrigger"
              echo "Commandes disponibles:"
              echo "  go build           - Construire l'application"
              echo "  make build         - Construire avec le Makefile"
              echo "  make dev-run       - Lancer en mode développement"
              echo "  make help          - Afficher l'aide complète du Makefile"
            '';
          };
        });
      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.hyprtrigger);
    };
}
