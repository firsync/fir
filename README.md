
I call it 'fir'. It's a tool for version control like git but simpler.

some problems with git I want to address are:
- git syntax is overly arcane and complex to new people, and hard to explain in a way that doesn't scare people off.
- reduce verbosity and procedural effort, instead of `git init`, `git add`, `git commit ...`, `git push`, just use `fir save` and `fir sync`.
- github repos, and git, have low portability for project metadata like issues. why not allow the user to just sync code or sync everything, issues included?

some goals/concepts of the project are
- simplify version control as it pertains to `git`, first and foremost
- assume all files are being tracked unless specifically `.ignore`d
- limited keywords for simplicity: `save`, `load`, `history`, `sync`
- push and pull during syncs
- better hashing for collision resistance (uses sha3 Shake256)
- remote repos are created for you on your favorite tracker on-push by default
- add support for pulling project data only, or full project data with issues included

some todos that haven't been accounted for yet:
- branches
- sig/verify, designed but not implemented, uses ed25519 for sig/verify.
- remote tracking software for running a github-like website for collaboration. this will likely be a project in itself, but id like to get the base version control taken care of first

priority:
- right now the first priority is just tracking, observing, and writing the history and diffs of what has changed on each save
- next add support for rewinding via rebuilding with the diffs for each checkpoint from the content in `.fir/checkpoints/base/`
- after rewinding and tracking are done, create syncing so that we can rsync our latest save checkpoints to our remote
- add signing/verifying with a global/local setting for default sig/verify on or off
- document every line, document every concept, document everything.