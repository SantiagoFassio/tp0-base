class Result:
    """
    Class to represent the result of processing a client signal. It contains:
    - ok: a boolean indicating if the processing was successful or not
    - message: a string with an error message in case ok is False, or an empty string
    - payload: any additional data to be returned with the result
    Used to return the result of processing a client signal in a structured way, 
    allowing to include both the success status and any relevant data or error messages.
    """
    def __init__(self, ok: bool, message: str = "", payload=None):
        self.ok = ok
        self.message = message
        self.payload = payload
        