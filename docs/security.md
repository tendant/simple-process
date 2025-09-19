# Security & Compliance Notes

This library is designed to run inside larger file-processing systems that may handle sensitive or regulated data. Keep the following practices in mind when extending or embedding it:

- **Principle of least privilege:** UoWs should only receive the presigned URLs and metadata they require. Avoid embedding raw credentials or long-lived tokens in jobs or artifacts.
- **Transport hygiene:** When enabling external transports (e.g., NATS, Kafka, HTTP callbacks), enforce TLS and authentication at the broker or gateway. CloudEvents metadata can be inspected without parsing the payload, so avoid leaking secrets through headers.
- **Artifact storage:** Configure `adapters.Storage` implementations (in-memory for tests, S3/MinIO behind the `s3` build tag, or your own) to write to segregated buckets/containers with appropriate retention policies. Document any encryption requirements in repository ADRs or PRs.
- **Telemetry:** If you introduce logging or tracing adapters, scrub PII before emission and label spans/fields so SIEM tooling can filter access patterns.
- **Dependency review:** New third-party SDKs (such as CloudEvents clients or message brokers) should be pinned in `go.mod`/`requirements.txt` equivalents and reviewed for license compatibility. Record major upgrades in the changelog or relevant docs.
- **Testing with external services:** Gate integration tests that rely on live infrastructure behind build tags or environment flags to prevent accidental data egress during CI runs.

Refer back to this document whenever you introduce a new adapter, transport, or data-contract change, and update it with additional guardrails learned in production.
