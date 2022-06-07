# src

Repo src holds some source code by hkjn.

## Submodules

Some directories hold source code from separate repos outside of `src/`, 
which can be initialized with:

```
git submodule init
git submodule update
```

These repos were added with:

```
git submodule add https://github.com/hkjn/lnmon
git submodule add https://github.com/hkjn/probes
git submodule add https://github.com/hkjn/prober
git submodule add https://github.com/hkjn/dashboard
```

The submodule repos can be updated like:

```
git submodule update --remote probes
```
