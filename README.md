# src

Repo src holds some source code by hkjn.

## Submodules

Some directories hold source code from separate repos outside of `src/`, 
which can be initialized with:

```
git submodule init
git submodule update
```

These repos were added like:

```
git submodule add https://github.com/hkjn/probes
```

The submodule repos can be updated like:

```
git submodule update --remote probes
```
