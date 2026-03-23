class Result:
    def __init__(self, ok: bool, message: str = "", payload=None):
        self.ok = ok
        self.message = message
        self.payload = payload
        