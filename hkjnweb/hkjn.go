// Package hkjnweb is the personal website hkjn.me.
//
// See https://github.com/hkjn/autosite for the framework that enables
// this site.
package hkjnweb

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"hkjn.me/autosite"
)

var (
	baseTemplates = []string{
		"tmpl/base.tmpl",
		"tmpl/base_header.tmpl",
		"tmpl/head.tmpl",
		"tmpl/style.tmpl",
		"tmpl/fonts.tmpl",
		"tmpl/js.tmpl",
	}
	isProd = false
	Logger = autosite.Glogger{}
)

// getLogger returns the autosite.Logger to use.
func getLogger(r *http.Request) autosite.Logger {
	return Logger
}

var (
	webDomain = "www.hkjn.me"
	web       = autosite.New(
		"Henrik Jonsson",
		"pages/*.tmpl", // glob for pages
		webDomain,      // live domain
		append(baseTemplates, "tmpl/page.tmpl"),
		getLogger,
		isProd,
		template.FuncMap{},
	)

	goImportTmpl = `<head>
    <meta http-equiv="refresh" content="0; URL='%s'">
    <meta name="go-import" content="%s git %s">
  </head>`

	redirects = map[string]string{}
)

// Register registers the handlers.
func Register(prod bool) {
	isProd = prod
	http.HandleFunc("/keybase.txt", keybaseHandler)
	if isProd {
		log.Println("We're in prod, remapping some paths")
		http.HandleFunc("/", nakedIndexHandler)
	} else {
		log.Println("We're not in prod, remapping some paths")
		http.HandleFunc("/nakedindex", nakedIndexHandler)
	}
	for uri, newUri := range redirects {
		web.AddRedirect(uri, newUri)
	}

	web.Register()
}

// nakedIndexHandler serves requests to hkjn.me/
func nakedIndexHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)
	l.Infof("nakedIndexHandler for URI %s\n", r.RequestURI)
	if r.URL.Path == "/" {
		url := "/index"
		l.Debugf("visitor to / of naked domain, redirecting to %q..\n", url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// Our response tells the `go get` tool where to find
		// `hkjn.me/[package]`.
		parts := strings.Split(r.URL.Path, "/")
		godocUrl := fmt.Sprintf("https://godoc.org/hkjn.me%s", r.URL.Path)
		repoRoot := fmt.Sprintf("https://github.com/hkjn/%s", parts[1])
		importPrefix := fmt.Sprintf("hkjn.me/%s", parts[1])
		fmt.Fprintf(w, goImportTmpl, godocUrl, importPrefix, repoRoot)
	}
}

var keybaseVerifyText = `
==================================================================
https://keybase.io/hkjn
--------------------------------------------------------------------

I hereby claim:

  * I am an admin of https://hkjn.me
  * I am hkjn (https://keybase.io/hkjn) on keybase.
  * I have a public key with fingerprint D618 7A03 A40A 3D56 62F5  4B46 03EF BF83 9A5F DC15

To do so, I am signing this object:

{
    "body": {
        "key": {
            "eldest_kid": "0101c778a84ebd54581e354296aa8171d79a23a3bcce4dbf99f880dfeafc2e0030030a",
            "fingerprint": "d6187a03a40a3d5662f54b4603efbf839a5fdc15",
            "host": "keybase.io",
            "key_id": "03efbf839a5fdc15",
            "kid": "0101c778a84ebd54581e354296aa8171d79a23a3bcce4dbf99f880dfeafc2e0030030a",
            "uid": "d1a36099363d76e027441cb534d6ba19",
            "username": "hkjn"
        },
        "revoke": {
            "sig_ids": [
                "dfa4fe72179e5cfbcee877c1df6c1b604d70f34efa6c59ffb144495a2d8ff4400f"
            ]
        },
        "service": {
            "hostname": "hkjn.me",
            "protocol": "https:"
        },
        "type": "web_service_binding",
        "version": 1
    },
    "ctime": 1452604724,
    "expire_in": 157680000,
    "prev": "6261a0a8692d7bf936641d51bbab2243eb1933769b780a09b59fac9d4b2570a4",
    "seqno": 6,
    "tag": "signature"
}

which yields the signature:

-----BEGIN PGP MESSAGE-----
Version: GnuPG v2

owGtUn1MVWUYvxeU+LRboTZGgWcWpHdwPt7zdcOPKwzLSJqQ6UBu7znnPZcjl3vp
3gtcMMaGtYHJwimK1IxF2ZBwKCBKGaZRGAgzN5CPYUxcmiBwET+qS51Da2vr3969
2/u+z36/3/N7nvepCvPXBer9pzv31lpGl+t77gi67Yd86/dggkMqwkx7sBy0eCCb
hFxuS44iYSYMJ3BCZFkOcgAJEg1ojkAUDUiegZAjWEJieUhSkBJEEQFJkHle5jhc
khGURRLhOKVtiBkxWbFbkTPPqdjdqqzEEBwLcQoCHFISzTCkTAMBMDiFZEHmKB7S
siQStErMdrg0hmpOgC4UpzjUmPqw/G3vv/j/2Xf+opxEQIrBeZ5iKIllEE6yABCi
QFNAYgRI8BrQhZx2mItUdHbObjtWYsScqMCRg7SeuhSratiFmTIwSYZARixJsDyi
RVkQEeJYViQkmREJgcGBxOIyBZAMGZHmZVkgAAA8DUmJk2UAcFzGdqnaarYCRVwU
1zr0r8xx6s2I5TkdbofosGlBtzvPZdIMuYvyNFQhEiwqv1ETsAiKXVL/RqUUIKdL
cdgxE6FCRbeiKRKAJlVPLAmMGPLkKU5kUTQEzTIcri4tESpQNRmSISAOOYYnJVbt
J8UwgJBoQhCgQJKAQgLBUxTL8ALL4RDnBbU2KPISEEiaxSHAtJLesTswE6P6hFZV
Um2aHbrznQgrCa7wi1yi0wfqApb6aUOrCw4y/DPJa3wG3ej6zvjwg4cNI0t/i4ha
+eMPmeazxqi5t/sHQ2pXe8y/xtVElW579sP2Ca9/UtT5zuL+An3l9HTjtacqO8CV
fUPD7312MfVIYuLwt7nO262Bl78MyRyKTu6e3XzX3np67NbeCGXL3b74ufBgh69p
eMXDjVt9XacHvfumzLFKw6V0KvZgfezjluhNZ/2Siud3xmRZz0njXROX3vrk0wlQ
FrSq4taZzvahe9UtI1xybu83Yt/KMcuCEHSospNPP7C5ztJmfDMiMiHpo5obO588
UtHdUd+sX/dn/AL0Hy8PSu0xGBIr94+2bksNs3tW+PoUMFNUnmLOOX7lxTLPfPPA
dk9z+h9fjz2fndxtu78pZs3cqYDx4+cNZmpWXFd69WX2cotuIq/X+jja+0HD9f22
2jde39VQH5cZEvrgJNoa3PVd8VTJjrbQxITCp821Bz4vjUm4+L75aPyFoVlvWnVF
5O3wDfGnMibDTqYMzIx8dayz5kIA+Hn5zOoFrBx9Uby2N73wmVDwhNe27OO5QevN
9mxb0GvX9KvK7/lXDzRcTeufqhIepbQRZyZfiukb0RFZnvGNjRu6m7x66ws9qLHu
5omm6WMen/K9ZdnRhvxHyfd3dxT3ntAVO2I9r/wU2dZW+At4dUnspN+DHU3nnks0
7ru+9qHku1HVeCc/bXQsY7DucJkv693f6+c7tuT8BQ==
=Xxo5
-----END PGP MESSAGE-----

And finally, I am proving ownership of this host by posting or
appending to this document.

View my publicly-auditable identity here: https://keybase.io/hkjn

==================================================================`

func keybaseHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)
	l.Infof("keybaseHandler for URI %s\n", r.RequestURI)
	fmt.Fprintf(w, keybaseVerifyText)
}
