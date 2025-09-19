"""Utilities for authoring Units of Work in Python."""

from __future__ import annotations

import functools
import json
from typing import Any, Callable, Dict

__all__ = ["uow", "deserialize_job", "serialize_result"]

JobDict = Dict[str, Any]
ResultDict = Dict[str, Any]


def deserialize_job(job_payload: str) -> JobDict:
    """Convert a JSON job payload into a dictionary."""
    if not job_payload:
        raise ValueError("job payload must be a non-empty string")
    return json.loads(job_payload)


def serialize_result(result: ResultDict) -> str:
    """Serialize a result dictionary into JSON."""
    if result is None:
        raise ValueError("result cannot be None")
    return json.dumps(result)


def uow(name: str) -> Callable[[Callable[[JobDict], ResultDict]], Callable[[str], str]]:
    """Decorator that turns a function into a UoW entrypoint.

    The wrapped function receives a job dictionary and must return a result dictionary
    adhering to the library's contract. The decorator captures the UoW name so runners
    can auto-register the handler.
    """

    if not name:
        raise ValueError("uow name must be provided")

    def decorator(func: Callable[[JobDict], ResultDict]) -> Callable[[str], str]:
        @functools.wraps(func)
        def wrapper(job_payload: str) -> str:
            job = deserialize_job(job_payload)
            result = func(job)
            return serialize_result(result)

        setattr(wrapper, "uow_name", name)
        return wrapper

    return decorator
