import uuid
from typing import List, Optional
 
class Error:
    def __init__(
        self,
        code: str,
        message: str,
        error_code_id: Optional[str] = None
    ):
        self.code = code
        self.message = message
        self.error_code_id = error_code_id
 
    def to_dict(self):
        return {
            "code": self.code,
            "message": self.message,
            "error_code_id": self.error_code_id
        }
 
 
class APIResponse:
    def __init__(
        self,
        data=None,
        errors: Optional[List[Error]] = None,
        code: int = 200,
        message: Optional[str] = None
    ):
        self.data = data
        self.errors = errors or []
        self.code = code
        self.request_id = str(uuid.uuid4())
 
        if message:
            self.message = message
        elif not self.errors:
            self.message = "Text extracted successfully"
        elif len(self.errors) == 1:
            self.message = self.errors[0].message
        else:
            self.message = "Something went wrong."
 
        self.status = "success" if not self.errors else "error"
 
    def to_dict(self):
        response = {
            "code": self.code,
            "status": self.status,
            "message": self.message
        }
 
        if self.status == "success":
            response["data"] = self.data
        else:
            response["errors"] = [e.to_dict() for e in self.errors]
 
        return response
 