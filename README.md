# codeowners-verifier
Small binary to validate the CODEOWNERS file and check if a path matches a CODEOWNER entry.
Useful for using on CI to ensure CODEOWNERS coverage.


## Docker
Available at: [Docker Hub](https://hub.docker.com/r/tfgco/codeowners-verifier) and [Quay.io](https://quay.io/repository/tfgco/codeowners-verifier).
You can change build arguments for GOOS and GOARCH using `--build-arg=linux --build-arg=arm`.

## Overview

codeowners-verifier verifies the entries inside a CODEOWNERS file. At the time being, it only supports Gitlab Code Owners. Although CodeOwners works the same way for Gitlab and Github (the same validation checks are done by both), Wildlife's implementation also checks for valid Users and Groups inside a CODEOWNERS file. 

## Environment Variables

+ `CODEOWNER_PROVIDER_URL`: The URL to the chosen provider. Each provider will have a default value. 
+ `CODEOWNER_PROVIDER_TOKEN`: Token to authenticate toward the chosen provider. There isn't default.
+ `CODEOWNER_PATH`: Path to the CODEOWNERS file. There isn't a default.

Those environment variables may also be defined by the respective flags: `--base-url`, `--codeowners` and `--token`.

A combination of using both flags and environment variables is possible, but keep in mind that flag values override environment variables values.

## Usage

:warning: Following how [.gitignore files works](https://git-scm.com/docs/gitignore), we didn't implement the "negate pattern". Github doesn't support it too. :warning:

**There are three verbs available: `help`, `verify` and `validate`.**

### Help

Help displays basic help message on the available verbs and global flags. Can also be used to get help on the available verbs:

```
$ ./codeowners-verifier help verify
```

### Verify

Verify must receive a path as argument. It then checks if the given path is covered by any of the existing entries.

Example:

```bash
$ codeowners-verifier verify dir1/
INFO[0000] Found matching rule on line 7: /**/ [@group1]
```

Verify supports the `-i (--ignore)` flag to ignore users/groups. It can be used multiples times and/or by a comma separated list of groups/users.

```bash
$ codeowners-verifier verify dir1/ -i @user1 --ignore @user2
INFO[0000] Found matching rule on line 7: /**/ [@group1]

$ codeowners-verifier verify dir1/ -i @user1 -i @user2,@group1
FATA[0000] Missing CODEOWNER entry, matched rule from line 7 don't have valid owners: /**/ [@group1]. Check your ignore rules.
```

### Validate

Validate validates the entire CODEOWNERS file, checking if the users and/or groups are valid. It does that by checking if the user or group is validy on the provider API.

It must receive the name of the provider. YOu can check for the available providers by executing the help:

```bash
$ codeowners-verifier validate -h
```

Then, execute with the provider name:

```bash
$ codeowners-verifier validate gitlab
INFO[0007] Valid CODEOWNERS file
```

In case something is wrong:

```bash
$ codeowners-verifier/codeowners-verifier validate gitlab
ERRO[0007] Error parsing line 7: user/group @user1 is invalid 
ERRO[0007] Error parsing line 8, path test-dir/ does not exist 
ERRO[0008] Error parsing line 8: user/group @group1 is invalid 
ERRO[0008] Error parsing line 8: user/group @group2 is invalid 
FATA[0008] Invalid CODEOWNERS file
```