#!/bin/bash

# Creates bare git repo. Meant to be run by 'git' user.
REPO=$1
mkdir /home/git/gitrepos/${REPO}
cd /home/git/gitrepos/${REPO}
git init --bare

