from common.utils import Bet

def process_header(header: str):
    """
    Process the header of the client message and return the action and n.
    Format example: "B|2\n"
    n can be the number of bets sent or the agency number.
    If the header format is invalid, a ValueError is thrown.
    """
    messages = header.strip().split("|")
    if len(messages) != 2 or messages[0] not in ("B", "E", "G"):
        raise ValueError("Invalid header format")
    return (messages[0], int(messages[1]))

def parse_bet(msg: str) -> Bet:
    """
    Parse a bet from a string and return a Bet object.
    The expected format of the bet string is:
    "agency|first_name|last_name|document|birthdate|number"
    If the format is invalid, returns None
    """
    parts = msg.strip().split("|")
    if len(parts) != 6:
        return None
    agency, first_name, last_name, document, birthdate, number = parts
    
    if (not agency.isdigit() or not document.isdigit() or not number.isdigit()):
        return None
    if len(birthdate.split("-")) != 3:
        return None
    if first_name == "" or last_name == "":
        return None
    
    return Bet(
        agency,
        first_name,
        last_name,
        document,
        birthdate,
        number,
    )

def parse_winners(winners: list[Bet]):
    """
    Creates and encodes the message to send the client with the winners.
    """
    msg = f"{len(winners)}\n"
    for winner in winners:
        msg += f"{winner.document}\n"
    return msg.encode("utf-8")
