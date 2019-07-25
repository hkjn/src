.DEFAULT_GOAL=update-packages

deploy-i3-config:
	cp config ~/.config/i3/
	cp brightness{inc,dec} ~/bin/

update-packages:
	cp configuration.nix /etc/nixos/
	nixos-rebuild switch
