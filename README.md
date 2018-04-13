# src

Repo src holds some source code by hkjn.

## Subtree

Some source code have their own repos outside of `src/`, added with:

```
$ git subtree add --prefix lnmon https://github.com/hkjn/lnmon.git master --squash
```

The subtree repos can be updated with:

```
$ git subtree pull --prefix lnmon https://github.com/hkjn/lnmon.git master --squash
```
