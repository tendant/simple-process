import json
import unittest

from simple_process_sdk.uow import deserialize_job, serialize_result, uow


class UoWDecoratorTests(unittest.TestCase):
    def test_roundtrip(self) -> None:
        payload = json.dumps({"job_id": "1", "uow": "echo"})

        @uow("echo")
        def echo(job):
            job["echoed"] = True
            return job

        result = echo(payload)
        parsed = json.loads(result)
        self.assertTrue(parsed["echoed"])
        self.assertEqual("echo", getattr(echo, "uow_name"))

    def test_helpers_validate_inputs(self) -> None:
        with self.assertRaises(ValueError):
            deserialize_job("")
        with self.assertRaises(ValueError):
            serialize_result(None)


if __name__ == "__main__":
    unittest.main()
