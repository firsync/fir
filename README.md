
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


I think there are some ideas that are forming around how `fir` should differ from `git`:

- we should keep a folder inside the current project directory to track our files, for this example lets call it `.ignore`
- instead of using `git add ...` we should just assume all files are tracked, unless they're listed in `.ignore`
- instead of requiring a `git init` we should initialize the `.fir` folder and write the first commit if one doesn't exist. 
- we should use a better hashing algorithm, like a sha3 hash for generating our hashes
- we should use `diff` to form our checkpoints
- we should use `rsync` for our syncing
- we should have a command called `fir history` that combines the behavior of `git log` and `git status` to show the user both the files that have changed since the last checkpoint, but also the list of historic checkpoints and messages
- new checkpoints should have a message attached if the user desires, and we should automatically insert a digest of the file changes at the end like "M files changed, N files added, O files removed" so that if a user just does `fir save` to make a checkpoint, without a message, they'll automatically have the file change summary in the message automatically.
- In the `./.fir/config` we should have spots for `name`, `email`, `remote`, and an optional `pubkey` 
- we should store the initial project data as a copy of itself in `./.fir/checkpoints/base/`
- when the user does `fir save` we should make prompt the user for a message, then create a commit/checkpoint in `./.fir/checkpoints/unix-timestamp.diff`
- when the user does `fir load <unix-timestamp>` we should rebuild the project from the `base` using the diffs up until the unix-timestamped diff the user has specified.



1. user runs `fir save` in the `~/code/thisproject` directory
2. `fir` checks for `~/code/thisproject/.fir/`, if the directory doesn't exist `fir` creates it.
3. `fir` checks for `~/code/thisproject/.fir/config`, if the file doesn't exist `fir` creates it. If it exists, `fir` loads the values into the current application state.
4. `fir` checks for `~/code/thisproject/.fir/checkpoints/base/`, if the folder doesn't exist, `fir` creates it and stores a copy of the existing project files in `~/code/thisproject/.fir/checkpoints/base/`
5. `fir` checks for `~/code/thisproject/.fir/checkpoints/<unix-timestamp>.diff` files and looks for the most recent timestamp among them.
6. `fir` sha3 hashes each file in the current project, and compares it to the list in `~/code/thisproject/.fir/checkpoints/<unix-timestamp>.list` where each file in that checkpoint is listed one per line like `<sha3 hash> <relative file path>`
7. if the sha3 hashes from the most recent `~/code/thisproject/.fir/checkpoints/<unix-timestamp>.list` and the current list of hashes don't match, create a new diff from the most recent diff. If no diffs exist, create a diff between `base` and the current state.
8. Once the new hash list and diff have been created, write them to this project's `.fir` folder.
9. Prompt the user for a message to go with their save. 
10. Take the user's input for the message and then append to it a string composed of: <# of files changed since last save>, <# of files added/removed since last save>, current `date` string, then write that info to `~/code/thisproject/.fir/checkpoints/<unix-timestamp>.meta`.  