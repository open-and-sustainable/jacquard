"""Synthetic smoke test for the C interface; uses only the Python stdlib."""

import ctypes
import json
from pathlib import Path
import sys
import unittest


LIBRARY = Path(sys.argv.pop(1)).resolve() if len(sys.argv) > 1 else None
FIXTURES = Path(__file__).parent / "fixtures"


class SharedLibraryTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        if LIBRARY is None:
            raise RuntimeError("usage: python3 test/shared_library_test.py path/to/jacquard.so")
        cls.lib = ctypes.CDLL(str(LIBRARY))
        cls.lib.FreeCString.argtypes = [ctypes.c_void_p]
        cls.lib.FreeCString.restype = None
        for name in (
            "RunValidationInput",
            "RunValidationOutput",
            "RunValidationFunnelInput",
            "RunValidationFunnelOutput",
        ):
            function = getattr(cls.lib, name)
            function.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
            function.restype = ctypes.c_void_p

    def call(self, name, data, version=b"v1"):
        pointer = getattr(self.lib, name)(data, version)
        self.assertTrue(pointer)
        try:
            return json.loads(ctypes.string_at(pointer))
        finally:
            self.lib.FreeCString(pointer)

    def test_all_synthetic_interfaces(self):
        for function, fixture in (
            ("RunValidationInput", "check_input"),
            ("RunValidationOutput", "check_output"),
            ("RunValidationFunnelInput", "funnel_input"),
            ("RunValidationFunnelOutput", "funnel_output"),
        ):
            with self.subTest(function=function):
                result = self.call(function, (FIXTURES / (fixture + ".json")).read_bytes())
                self.assertEqual(result, {"valid": True, "issues": []})

    def test_errors_are_json(self):
        self.assertFalse(self.call("RunValidationInput", b"{")["valid"])
        self.assertIn("error", self.call("RunValidationInput", b"{}", b"v99"))
        self.assertIn("error", self.call("RunValidationInput", None))
        self.assertIn("error", self.call("RunValidationInput", b"{}", None))
        self.lib.FreeCString(None)

    def test_repeated_allocations(self):
        payload = (FIXTURES / "check_input.json").read_bytes()
        for _ in range(40):
            self.assertTrue(self.call("RunValidationInput", payload)["valid"])


if __name__ == "__main__":
    unittest.main()
