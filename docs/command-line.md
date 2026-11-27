---
title: Command Line
layout: default
---

# Command line

Build the current interface-validation command:

```sh
go build -o bin/jacquard ./cmd/jacquard
```

Validate one JSON value from stdin:

```sh
bin/jacquard validate -type input -version v1 < test/fixtures/check_input.json
bin/jacquard validate -type funnel_input < test/fixtures/funnel_input.json
```

Print an embedded schema:

```sh
bin/jacquard schema -type output
```

The types are `input`, `output`, `funnel_input`, and `funnel_output`.
Validation writes `{"valid":true,"issues":[]}` on success. Invalid instances return
issues and exit 1. Unknown schema versions/types, invocation errors, and I/O errors
exit 2. Successful validation and schema display exit 0. Diagnostic usage output
goes to stderr.

The first-release `check` and `funnel` commands, including JSON Lines mode, are
planned in the [checker contract](checker-contract.md). They are not implemented
by this validation command.
