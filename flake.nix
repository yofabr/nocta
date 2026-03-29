{
  description = "Nocta - GTK port manager with Fyne UI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in
    {
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gtk3
          glib
          cairo
          pango
          gdk-pixbuf
          libGL
          xorg.libX11
          xorg.libXcursor
          xorg.libXrandr
          libxkbcommon
        ];

        shellHook = ''
          echo "Nocta dev shell - run: go run ./cmd/nocta"
        '';
      };
    };
}