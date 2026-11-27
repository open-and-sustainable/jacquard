---
title: C-Shared Library
layout: default
---

# C-shared library

The current C exports validate interface shapes using the root `jacquard` Go
package, as does the CLI. Evidence checking and audit computation are planned separately.

## Build

From the repository root, with Go and a C compiler:

```sh
mkdir -p dist
go build -buildmode=c-shared -o dist/jacquard.so ./utils/sharedlib
```

Use `.dylib` on macOS or `.dll` on Windows. Go generates the corresponding
`jacquard.h`; distribute the header with the library. Release CI is configured to
build native libraries and headers when a release is published.

## Current functions

- `RunValidationInput(input, version)`
- `RunValidationOutput(input, version)`
- `RunValidationFunnelInput(input, version)`
- `RunValidationFunnelOutput(input, version)`
- `FreeCString(result)`

Each validation function accepts two NUL-terminated UTF-8 strings and returns an
owned NUL-terminated JSON string. `version` is `v1`. Results contain `valid` and
`issues`, plus `error` for an unavailable schema or invalid invocation. NULL inputs
produce a JSON error. No diagnostic text is printed by the library.

**Free each result exactly once with `FreeCString`.** Input buffers remain owned by
the caller. `FreeCString(NULL)` is safe. Keep the library loaded until its returned
strings have been freed. The future `RunCheck` and `RunAuditFunnel` exports will use
the same ownership rule and the envelopes specified in the checker contract.

## C

```c
#include <stdio.h>
#include "jacquard.h"

int main(void) {
    char *result = RunValidationInput("{}", "v1");
    puts(result);
    FreeCString(result);
    return 0;
}
```

## Python

Preserve the returned pointer with `c_void_p`; automatic conversion to Python bytes
would lose the pointer that must be freed.

```python
import ctypes
from pathlib import Path

lib = ctypes.CDLL(str(Path("dist/jacquard.so").resolve()))
lib.RunValidationInput.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
lib.RunValidationInput.restype = ctypes.c_void_p
lib.FreeCString.argtypes = [ctypes.c_void_p]
lib.FreeCString.restype = None

data = Path("test/fixtures/check_input.json").read_bytes()
ptr = lib.RunValidationInput(data, b"v1")
try:
    print(ctypes.string_at(ptr).decode("utf-8"))
finally:
    lib.FreeCString(ptr)
```

## Synthetic FFI smoke test

```sh
python3 test/shared_library_test.py dist/jacquard.so
```
