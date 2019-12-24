# Edit this configuration file to define what should be installed on
# your system.  Help is available in the configuration.nix(5) man page
# and in the NixOS manual (accessible by running ‘nixos-help’).

{ config, pkgs, ... }:

{
  imports =
    [ # Include the results of the hardware scan.
      ./hardware-configuration.nix
    ];

  # Use the systemd-boot EFI boot loader.
  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;

  # Minimal list of modules to use the EFI system partition and the YubiKey
  boot.initrd.kernelModules = [ "vfat" "nls_cp437" "nls_iso8859-1" "usbhid" ];

  # Crypto setup, set modules accordingly
  boot.initrd.luks.cryptoModules = [ "aes" "xts" "sha512" ];

  # Enable support for the YubiKey PBA
  # boot.initrd.luks.yubikeySupport = true;

  # Mount USB key before trying to decrypt root filesystem
  boot.initrd.postDeviceCommands = pkgs.lib.mkBefore ''
    echo "waiting for usb disk..";
    mkdir -m 0755 -p /key;
    sleep 5; # To make sure the usb key has been loaded
    echo "mounting apollo (expecting id usb-Generic_Flash_Disk_18082009002113-0:0-part1): $(ls -hsal /dev/disk/by-id)";
    if ! mount -n -t vfat -o ro /dev/disk/by-id/usb-Generic_Flash_Disk_18082009002113-0:0-part1 /key; then
      echo "apollo not found, waiting for sangus..";
      sleep 5;
      echo "mounting sangus (expecting id usb-SMI_USB_DISK_AA00000000014172-0:0-part1): $(ls -hsal /dev/disk/by-id)";
      mount -n -t vfat -o ro /dev/disk/by-id/usb-SMI_USB_DISK_AA00000000014172-0:0-part1 /key || echo "could not find sangus either, giving up."
    fi;
  '';

  # Configuration to use Luks device.
  boot.initrd.luks.devices = [ {
    name = "nixos-enc";
    device = "/dev/disk/by-id/nvme-SAMSUNG_MZVLB1T0HALR-000L7_S3TPNX0M115074-part1";
    keyFile = "/key/keys/.nix.key";
    preLVM = false;
  } ];

  # Add kvm kernel modules.
  boot.kernelModules = [ "kvm-intel" ];

  services.udev.packages = [
    pkgs.ledger-udev-rules
    pkgs.yubikey-personalization
    pkgs.libu2f-host
  ];
  services.udev.extraRules = ''
    # coldcard udev rules
    SUBSYSTEM=="usb", ENV{ID_VENDOR}=="d13e", ENV{ID_PRODUCT}=="cc10", MODE="0660", GROUP="users", TAG+="uaccess", TAG+="udev-acl", SYMLINK+="+coldcard%n"
  '';

  # Set hostname.
  networking.hostName = "velletri";
  # Provide networkmanager (control with nmtui) for easy wireless configuration.
  networking.networkmanager.enable = true;

  # Enable power management, powertop, and cap max frequency.
  powerManagement.enable = true;
  powerManagement.cpuFreqGovernor = "powersave";
  powerManagement.cpufreq.max = 1600000;
  powerManagement.powertop.enable = true;

  # Set time zone:
  # https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
  time.timeZone = "America/Tijuana";

  # List packages installed in system profile. To search, run:
  # $ nix search wget
  environment.systemPackages = with pkgs; [
     arandr
     bind
     bitcoin
     clightning
     electrum
     exfat
     firefox
     git
     gnupg
     gnumake
     go
     file
     imagemagick
     jhead
     libosinfo
     htop
     glxinfo
     gptfdisk
     mosh
     mkpasswd
     mplayer
     mtpaint
     ncdu
     nmap
     prometheus
     pwgen
     python3
     redshift
     scrot
     srm
     steghide
     tightvnc
     tmux
     tor
     tor-browser-bundle-bin
     vim
     wget
     qemu
     unzip
     urbit
     usbutils
     virtmanager
     virtviewer
     wasabiwallet
     whois
     youtube-dl-light
     xautolock
     xbrightness
     xclip
     xfce.terminal
  ];

  # Set environment to support gpg-agent as ssh-agent.
  environment.shellInit = ''
    export GPG_TTY="$(tty)"
    gpg-connect-agent /bye
  '';

  services.xserver = {
    desktopManager = {
      xterm.enable = false;
    };

    windowManager.i3 = {
      enable = true;
      extraPackages = with pkgs; [
        dmenu
        i3status
        i3lock
        i3blocks
      ];
    };
  };

  # Enable the X11 windowing system.
  services.xserver.enable = true;
  services.xserver.layout = "us";

  # Enable touchpad support.
  services.xserver.libinput.enable = true;

  # Some programs need SUID wrappers, can be configured further or are
  # started in user sessions.
  # programs.mtr.enable = true;
  programs = {
    ssh.startAgent = false;
    ssh.extraConfig = "Host admin1.hkjn.me\nPort 6200\nUser zero\nHost admin2.hkjn.me\nPort 6200\nUser zero";
    gnupg.agent = {
      enable = true;
      enableSSHSupport = true;
    };
    bash.shellInit = ''
      export SSH_AUTH_SOCK="/run/user/$UID/gnupg/S.gpg-agent.ssh"
      alias xcl="xclip -selection c"
      alias pp="git pull && git push"
      alias vms="virt-manager -c qemu:///session"
    '';

  };

  services.prometheus.enable = true;
  services.prometheus.listenAddress = "127.0.0.1:9090";
  services.prometheus.extraFlags = [
    "--storage.tsdb.retention.size=5GB"
  ];

  services.prometheus.scrapeConfigs = [
    {
      job_name = "node";
      static_configs = [
        {
          targets = [
            "localhost:9100"
          ];
        }
      ];
    }
    {
      job_name = "custom";
      static_configs = [
        {
          targets = [
            "localhost:8335"
          ];
        }
      ];
    }
  ];
  services.prometheus.exporters = {
    node = {
      listenAddress = "127.0.0.1";
      enable = true;
      enabledCollectors = [
        "conntrack"
        "diskstats"
        "entropy"
        "filefd"
        "filesystem"
        "loadavg"
        "mdadm"
        "meminfo"
        "netdev"
        "netstat"
        "stat"
        "time"
        "vmstat"
        "systemd"
        "logind"
        "interrupts"
        "ksmd"
      ];
    };
    # custom.enable = true;
    tor.enable = true;
  };

  # Enable grafana service.
  services.grafana = {
    enable   = true;
    port     = 3000;
    domain   = "localhost";
    protocol = "http";
    dataDir  = "/var/lib/grafana";
    provision.datasources = [
      {
        type = "prometheus";
        url = "http://localhost:9090";
      }
    ];
  };

  # Enable tor and tor client.
  services.tor.enable = true;
  services.tor.client.enable = true;
  services.tor.controlPort = 9051;
  services.tor.hiddenServices.lightningd = {
    map = [{
      port = 9735;
    }];
    version = 3;
  };

  # Open ports in the firewall.
  networking.firewall.allowedTCPPorts = [ 8333 ];
  # networking.firewall.allowedUDPPorts = [ ... ];

  # Enable bitcoind.
  systemd.user.services.bitcoin = {
     description = "Bitcoin Core daemon";
       serviceConfig = {
         ExecStart = "${pkgs.bitcoin}/bin/bitcoind";
         Restart = "on-failure";
     };
     wantedBy = [ "default.target" ];
  };
  systemd.user.services.bitcoin.enable = true;
  # Enable clightning.
  systemd.user.services.clightning = {
     description = "c-lightning daemon";
       serviceConfig = {
         ExecStart = "${pkgs.clightning}/bin/lightningd";
         Restart = "on-failure";
     };
     wantedBy = [ "default.target" ];
  };
  systemd.user.services.clightning.enable = true;
  # Enable electrs.
  # systemd.user.services.electrs = {
  #   description = "electrs daemon";
  #     serviceConfig = {
  #       ExecStart = "xx-pkgs.electrs-goes-here/bin/electrs";
  #       Restart = "on-failure";
  #   };
  #   wantedBy = [ "default.target" ];
  # };
  # systemd.user.services.electrs.enable = true;

  # Enable sound.
  sound.enable = true;
  hardware.pulseaudio.enable = true;

  # Enable docker.
  virtualisation.docker.enable = true;

  # Enable libvirtd for KVM.
  virtualisation.libvirtd.enable = true;

  # Define a user account. Don't forget to set a password with ‘passwd’.
  users.users.user = {
    isNormalUser = true;
    # Enable sudo and networkmanager for the user.
    extraGroups = [ "wheel" "networkmanager" "hardwarewallets" "libvirtd" ];
    # Set default password.
    initialHashedPassword = "$6$gQ/dMey1PH$aKVUdM1EybW2iFGC80cOby/S2nQNpn3SlCzl3mk7IU39A5b4ew22cAxvpOx8N7yZZ..IOB4vWdnp8ZPrmJvHT0";
  };

  # Enable automatic nix garbage collection of unreferenced packages:
  nix.gc.automatic = true;
  nix.gc.dates = "03:15";

  # This value determines the NixOS release with which your system is to be
  # compatible, in order to avoid breaking some software such as database
  # servers. You should change this only after NixOS release notes say you
  # should.
  system.stateVersion = "19.03"; # Did you read the comment?
}
