from common.utils import Bet

def process_header(header: str):
    messages = header.strip().split("|")
    if len(messages) != 2 or messages[0] not in ("B", "E", "G"):
        raise ValueError("Invalid header format")
    return (messages[0], int(messages[1]))

def parse_bet(msg: str) -> Bet:
    
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
    msg = f"{len(winners)}\n"
    for winner in winners:
        msg += f"{winner.document}\n"
    return msg.encode("utf-8")
