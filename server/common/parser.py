from common.utils import Bet

def process_header(header: str):
    messages = header.strip().split("|")
    if len(messages) != 2 or messages[0] not in ("B", "E", "G"):
        raise ValueError("Invalid header format")
    return (messages[0], int(messages[1]))

def parse_bet(msg: str) -> Bet:
    parts = msg.strip().split("|")

    # TODO
    # Comprobar la bet
    if len(parts) != 6:
        return None
    
    return Bet(
        agency = parts[0],
        first_name = parts[1],
        last_name = parts[2],
        document = parts[3],
        birthdate = parts[4],
        number = parts[5]
    )

def parse_winners(winners: list[Bet]):
    msg = f"{len(winners)}\n"
    for winner in winners:
        msg += f"{winner.document}\n"
    return msg.encode("utf-8")
