'''
from simple_process_sdk.uow import uow
from simple_process_sdk.utils import pypdf, tesseract

@uow("ocr_pdf")
def ocr_pdf(job):
    # TODO: Download file from presigned URL
    # For now, assume the file is available locally
    file_path = job["file"]["blob"]["location"]

    text = ""
    with open(file_path, "rb") as f:
        pdf = pypdf.PdfReader(f)
        for page in pdf.pages:
            # This is a simplified example. A real implementation would handle
            # images within the PDF and use Tesseract for OCR.
            text += page.extract_text()

    # TODO: Upload transcript as an artifact
    return {
        "job_id": job["job_id"],
        "uow": job["uow"],
        "file_id": job["file"]["id"],
        "attributes_patch": {
            "pages": len(pdf.pages),
        },
        "artifacts": [
            {
                "kind": "transcript",
                "mime": "text/plain",
                "bytes": len(text),
                "location": "s3://bucket/artifacts/{}/ocr.txt".format(job["file"]["id"])
            }
        ]
    }
'''
