# Data Contracts

This document describes the JSON data contracts for Jobs and Results used in the File Processing Pipeline Library.

## Job

The `Job` contract is the primary input to any Unit of Work (UoW). It contains all the necessary information for the UoW to perform its task.

```json
{
  "version": "1.0",
  "job_id": "j_abc123",
  "uow": "ocr_pdf",
  "file": {
    "id": "f_123",
    "tenant_id": "t_1",
    "owner_id": "u_42",
    "name": "scan.pdf",
    "blob": { "location": "s3://bucket/key", "size": 1234567 },
    "attributes": { "mime": "application/pdf", "sha256": "..." }
  },
  "presigned_get": "https://s3/...sig...",
  "return": { "type": "http", "url": "https://engine/uow/callback" },
  "idem_key": "t_1:u_42:<sha256>",
  "hints": { "language": "en" }
}
```

## Result

The `Result` contract is the output of a UoW. It contains any patched attributes and a list of generated artifacts.

```json
{
  "job_id": "j_abc123",
  "uow": "ocr_pdf",
  "file_id": "f_123",
  "attributes_patch": { "pages": 12 },
  "artifacts": [
    { "kind": "transcript", "mime": "text/plain", "bytes": 54231,
      "location": "s3://bucket/artifacts/f_123/ocr.txt" }
  ]
}
```

## CloudEvents Envelope

When jobs travel over external transports (e.g., the NATS bus), they are wrapped in a minimal [CloudEvents 1.0](https://cloudevents.io) envelope before delivery. The event header adds routing metadata while the `data` field carries the JSON job payload described above.

```json
{
  "specversion": "1.0",
  "type": "simpleprocess.job",
  "source": "simple-process/nats",
  "id": "j_abc123",
  "time": "2024-05-06T12:34:56Z",
  "datacontenttype": "application/json",
  "data": {
    "job_id": "j_abc123",
    "uow": "ocr_pdf",
    "file": { "id": "f_123", "blob": { "location": "s3://bucket/key" } }
  }
}
```

Consumers written in other languages can rely on the CloudEvents SDK for parsing, then hand the decoded job to their local UoW implementation.
