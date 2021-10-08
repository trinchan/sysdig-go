# How to contribute #

I appreciate any and all help in developing the client -- patches welcome!

## Reporting issues ##
Bugs, feature requests, and development-related questions should be directed to
the [GitHub issue tracker](https://github.com/trinchan/sysdig-go/issues).

If reporting a bug, please try and provide as much context as possible such as
your operating system, Go version, and anything else that might be relevant to
the bug. Including (redacted) client logs with `Debug` enabled can help for parse or request errors.

For feature requests, please explain what you're trying to do, and
how the requested feature would help you do that.

## Submitting a patch ##

1. It's generally best to start by opening a new issue describing the bug or
   feature you're intending to fix. Even if you think it's relatively minor,
   it's helpful to know what people are working on. Mention in the initial
   issue that you are planning to work on that bug or feature so that it can
   be assigned to you.

1. Follow the normal process of [forking](https://help.github.com/articles/fork-a-repo) the project, and setup a new
   branch to work in. It's important that each group of changes be done in
   separate branches in order to ensure that a pull request only includes the
   commits related to that bug or feature.

1. Go makes it very simple to ensure properly formatted code, so always run
   `go fmt` on your code before committing it.

1. Any significant changes should almost always be accompanied by tests. The
   project already has good test coverage, so look at some of the existing
   tests if you're unsure how to go about it. [gocov](https://github.com/axw/gocov) and [gocov-html](https://github.com/matm/gocov-html)
   are invaluable tools for seeing which parts of your code aren't being
   exercised by your tests.

1. Please run:`
    * `go test github.com/trinchan/sysdig-go/...`
    * `go vet github.com/trinchan/sysdig-go/...`

1. Do your best to have [well-formed commit messages](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html) for each change.
   This provides consistency throughout the project, and ensures that commit
   messages are formatted properly by various git tools.

1. Finally, push the commits to your fork and submit a [pull request](https://help.github.com/articles/creating-a-pull-request).
   **NOTE:** Please do not use force-push on PRs in this repo, as it makes
   it more difficult for reviewers to see what has changed since the last
   code review.
