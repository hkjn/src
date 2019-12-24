.DEFAULT_GOAL=update-system-packages

deploy-i3-config:
	cp i3_config ~/.config/i3/config
	cp brightness{inc,dec} lockscreen ~/bin/
	cp i3status_config ~/.config/i3status/config
	cp terminalrc ~/.config/xfce4/terminal/
	# note that we install these scripts ourselves in /usr/local/bin,
	# to be able to reference them in i3_config, but there's likely
	# a more "pure" nixos way to receive the same results
	sudo install brightnessdec brightnessinc lockscreen /usr/local/bin/

refresh-packages:
	bash refresh_packages

install-user-config:
	cp config.nix /home/user/.config/nixpkgs/
	nix-env -iA nixpkgs.tools

update-system-packages:
	cp configuration.nix /etc/nixos/
	nixos-rebuild switch -I nixpkgs=/home/user/src/hkjn.me/nixpkgs --upgrade --profile-name test0
