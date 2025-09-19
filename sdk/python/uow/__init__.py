import functools
import json

def uow(name):
    """Decorator to register a Python function as a UoW."""
    def decorator(func):
        @functools.wraps(func)
        def wrapper(job_str):
            job = json.loads(job_str)
            # TODO: Add helper functions for downloading input file via presigned URL
            # and uploading artifacts.
            result = func(job)
            return json.dumps(result)
        wrapper._uow_name = name
        return wrapper
    return decorator

