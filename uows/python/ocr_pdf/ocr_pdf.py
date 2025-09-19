"""Reference OCR UoW implemented in Python.

This example keeps dependencies to the standard library so it can run anywhere. It
pretends the incoming file is a UTF-8 encoded PDF transcript and produces a simple
artifact containing the extracted text while updating metadata with word counts.
"""

from __future__ import annotations

from pathlib import Path
from typing import Any, Dict

from simple_process_sdk.uow import uow


def _read_text(path: Path) -> str:
    if not path.exists():
        raise FileNotFoundError(f"input file {path} not found")
    return path.read_text(encoding="utf-8", errors="ignore")


@uow("ocr_pdf")
def ocr_pdf(job: Dict[str, Any]) -> Dict[str, Any]:
    file_blob = job.get("file", {}).get("blob", {})
    file_id = job.get("file", {}).get("id")
    file_path = Path(file_blob.get("location", ""))

    text = _read_text(file_path)
    words = len(text.split())
    artifact_location = f"{file_path}.transcript.txt"

    return {
        "job_id": job.get("job_id"),
        "uow": job.get("uow"),
        "file_id": file_id,
        "attributes_patch": {
            "ocr_words": words,
        },
        "artifacts": [
            {
                "kind": "transcript",
                "mime": "text/plain",
                "bytes": len(text.encode("utf-8")),
                "location": artifact_location,
            }
        ],
    }
