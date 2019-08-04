.DEFAULT_GOAL=update-packages

deploy-i3-config:
	cp i3_config ~/.config/i3/config
	cp brightness{inc,dec} ~/bin/
	cp i3status_config ~/.config/i3status/config

refresh-packages:
	bash refresh_packages

update-packages:
	cp configuration.nix /etc/nixos/
	nixos-rebuild switch -I nixpkgs=/home/user/src/hkjn.me/nixpkgs
