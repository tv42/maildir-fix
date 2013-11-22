maildir-fix -- Fix Maildirs after git has pruned empty dirs
===========================================================

`maildir-fix` makes sure the subdirectories required for a Maildir are
present: `new`, `cur`, `tmp`.

This is useful if you're the kind of crazy that stores email in Git
repositories, as Git will prune directories that become empty.

Use like this:

    $ ls Spam
	cur/  tmp/
	$ maildir-fix Spam
	$ ls Spam
	cur/  new/  tmp/

It can process a Binc IMAP -style mail depot with the `-depot=PATH`
flag:

	$ maildir-fix -depot=.


Git hooks
---------

You can make `git` run `maildir-fix` automatically, by creating the
[git hooks](https://www.kernel.org/pub/software/scm/git/docs/githooks.html)
`post-checkout`, `post-merge` and `post-rewrite` with the following:

```
#!/bin/sh
set -e

. "$(git --exec-path)/git-sh-setup"

require_work_tree
cd_to_toplevel

maildir-fix -depot=.
```
