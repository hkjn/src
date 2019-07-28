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
    sleep 5;
    echo "creating /key"
    mkdir -m 0755 -p /key;
    sleep 2; # To make sure the usb key has been loaded
    echo "mounting disk: $(ls -hsal /dev/disk/by-id)";
    mount -n -t vfat -o ro /dev/disk/by-id/usb-SMI_USB_DISK_AA00000000014172-0:0-part1 /key;
  '';

  # Configuration to use Luks device.
  boot.initrd.luks.devices = [ {
    name = "nixos-enc";
    device = "/dev/disk/by-id/nvme-SAMSUNG_MZVLB1T0HALR-000L7_S3TPNX0M115074-part1";
    keyFile = "/key/keys/.nix.key";
    preLVM = false;
  } ];

  # Add udev packages for yubikey device.
  services.udev.packages = [
    pkgs.yubikey-personalization
    pkgs.libu2f-host
  ];

  # Set hostname.
  networking.hostName = "velletri";
  # Provide networkmanager (control with nmtui) for easy wireless configuration.
  networking.networkmanager.enable = true;

  # Enable power management, powertop, and cap max frequency.
  powerManagement.enable = true;
  powerManagement.cpuFreqGovernor = "powersave";
  powerManagement.cpufreq.max = 2600000;
  powerManagement.powertop.enable = true;

  # Set your time zone.
  # time.timeZone = "Europe/Amsterdam";

  # List packages installed in system profile. To search, run:
  # $ nix search wget
  environment.systemPackages = with pkgs; [
     electrum
     firefox
     git
     gnupg
     gnumake
     go
     file
     htop
     mosh
     mkpasswd
     mplayer
     ncdu
     nmap
     pwgen
     redshift
     srm
     tmux
     vim
     wget
     usbutils
     xbrightness
     xclip
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
  # programs.gnupg.agent = { enable = true; enableSSHSupport = true; };
  programs = {
    ssh.startAgent = false;
    ssh.extraConfig = "Host admin1.hkjn.me\nPort 6200\nUser zero\nHost admin2.hkjn.me\nPort 6200\nUser zero";
    gnupg.agent = {
      enable = true;
      enableSSHSupport = true;
    };
    bash.shellInit = ''
      export SSH_AUTH_SOCK="/run/user/$UID/gnupg/S.gpg-agent.ssh"
      alias elec="electrum --oneserver --server=127.0.0.1:50001:t"
      alias xcl="xclip -selection c"
    '';

  };

  # List services that you want to enable:

  # Enable the OpenSSH daemon.
  # services.openssh.enable = true;

  # Open ports in the firewall.
  # networking.firewall.allowedTCPPorts = [ ... ];
  # networking.firewall.allowedUDPPorts = [ ... ];
  # Or disable the firewall altogether.
  # networking.firewall.enable = false;

  # Enable sound.
  sound.enable = true;
  hardware.pulseaudio.enable = true;

  virtualisation.docker.enable = true;

  # Define a user account. Don't forget to set a password with ‘passwd’.
  users.users.user = {
    isNormalUser = true;
    extraGroups = [ "wheel" "networkmanager" ]; # Enable sudo and networkmanager for the user.
    initialHashedPassword = "$6$gQ/dMey1PH$aKVUdM1EybW2iFGC80cOby/S2nQNpn3SlCzl3mk7IU39A5b4ew22cAxvpOx8N7yZZ..IOB4vWdnp8ZPrmJvHT0";
  };

  # This value determines the NixOS release with which your system is to be
  # compatible, in order to avoid breaking some software such as database
  # servers. You should change this only after NixOS release notes say you
  # should.
  system.stateVersion = "19.03"; # Did you read the comment?
}
