
from common.utils import Bet

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

def parse_batch(expected, lines) -> list[Bet]:

    if len(lines) < 1:
        raise ValueError("Empty batch.")

    bets = []
    for line in lines:
        if line.strip() == "":
            continue
        bet = parse_bet(line.strip())
        if bet is not None:
            bets.append(bet)
    
    return bets
