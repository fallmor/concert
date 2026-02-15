{
  description = "Concert Booking System";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go_1_24 or pkgs.go_1_23 or pkgs.go;
        nodejs = pkgs.nodejs_22 or pkgs.nodejs;
        
    
        server = pkgs.buildGoModule {
          pname = "concert-server";
          version = "1.0.0";
          src = ./.;
          vendorHash = null;
          buildFlags = [ "-ldflags=-s -w" ];
          env = {
            CGO_ENABLED = "0";
          };
          subPackages = [ "cmd/server" ];
        };
        
        worker = pkgs.buildGoModule {
          pname = "concert-worker";
          version = "1.0.0";
          src = ./.;
          vendorHash = null;
          buildFlags = [ "-ldflags=-s -w" ];
          env = {
            CGO_ENABLED = "0";
          };
          # Build worker.go as main package
          buildPhase = ''
            runHook preBuild
            go build -o worker -ldflags="-s -w" ./temporal/worker.go
            runHook postBuild
          '';
          installPhase = ''
            runHook preInstall
            mkdir -p $out/bin
            cp worker $out/bin/worker
            runHook postInstall
          '';
        };
        
        frontend = pkgs.stdenv.mkDerivation {
          pname = "concert-frontend";
          version = "1.0.0";
          # Use static directory directly as source
          src = ./static;
          
          buildPhase = ''
            echo "Using pre-built frontend from static/ directory"
            mkdir -p dist
            cp -r ./* dist/ 2>/dev/null || true
          '';
          
          installPhase = ''
            mkdir -p $out
            if [ -d "dist" ] && [ -n "$(ls -A dist 2>/dev/null)" ]; then
              cp -r dist/* $out/
            else
              # copy directly if dist doesn't exist
              cp -r ./* $out/ 2>/dev/null || true
            fi
          '';
        };
        
        dockerImage = pkgs.dockerTools.buildLayeredImage {
          name = "concert-server";
          tag = "latest";
          
          contents = [
            pkgs.cacert
            pkgs.bash
            pkgs.coreutils
          ];
          
          extraCommands = ''
            mkdir -p root/static
            cp ${server}/bin/server root/server
            chmod +x root/server
            if [ -d "${frontend}" ]; then
              cp -r ${frontend}/* root/static/ || true
            fi
          '';
          
          config = {
            Cmd = [ "/root/server" ];
            WorkingDir = "/root";
            ExposedPorts = {
              "8080/tcp" = {};
            };
            Env = [
              "PATH=/usr/bin:/bin"
            ];
          };
        };
        
        workerDockerImage = pkgs.dockerTools.buildLayeredImage {
          name = "concert-worker";
          tag = "latest";
          
          contents = [
            pkgs.cacert
          ];
          
          extraCommands = ''
            mkdir -p root
            cp ${worker}/bin/worker root/worker
            chmod +x root/worker
          '';
          
          config = {
            Cmd = [ "/root/worker" ];
            WorkingDir = "/root";
            Env = [
              "PATH=/usr/bin:/bin"
            ];
          };
        };
        
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            nodejs
            npm
            postgresql
          ];
          shellHook = ''
            echo "Concert Booking System Development Environment"
            echo "Go: $(go version)"
            echo "Node: $(node --version)"
            echo "NPM: $(npm --version)"
          '';
        };
        
        packages.server = server;
        packages.worker = worker;
        packages.frontend = frontend;
        packages.docker = dockerImage;
        packages.dockerWorker = workerDockerImage;
        packages.default = dockerImage;
      });
}