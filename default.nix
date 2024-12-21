let
  pkgs = import <nixpkgs> {};
in pkgs.buildGoModule.override { go = pkgs.go_1_22; } {
  pname = "mygo-app";
  version = "1.0.0";
  src = ./.;
  vendorHash = null;
  
   buildInputs = with pkgs; [
    postgresql
  ];
  
  buildPhase = ''
  cd cmd/server
    export CGO_ENABLED=0
    go build -ldflags="-s -w -extldflags=-static" 
  '';
 installPhase = ''
    mkdir -p $out/bin
    cp server $out/bin/
  '';
  meta = with pkgs.lib; {
    description = "Simple api rest written in Go";
    homepage = "https://github.com/fallmor";
    license = licenses.mit;
    maintainers = [ "Mor FALL" ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}