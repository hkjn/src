.DEFAULT_GOAL=update-system-packages

deploy-i3-config:
	cp i3_config ~/.config/i3/config
	cp brightness{inc,dec} ~/bin/
	cp i3status_config ~/.config/i3status/config
	cp terminalrc ~/.config/xfce4/terminal/

refresh-packages:
	bash refresh_packages

update-packages:
	cp config.nix /home/user/.config/nixpkgs/
	nix-env -iA nixpkgs.tools

update-system-packages:
	cp configuration.nix /etc/nixos/
	nixos-rebuild switch -I nixpkgs=/home/user/src/hkjn.me/nixpkgs --upgrade -p test0
