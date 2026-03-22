
from common.utils import Bet

def parse_bet(msg: str) -> Bet:
    parts = msg.strip().split("|")

    # TODO
    # Comprobar la bet
    if len(parts) != 5:
        return None
    if parts[3] == "":
        return None
    
    return Bet(
        agency = parts[0],
        first_name = parts[1],
        last_name = parts[2],
        document = parts[3],
        birthdate = parts[4],
        number = parts[5]
    )

def parse_batch(expected, lines) -> list[Bet]:

    if len(lines) < 1:
        raise ValueError("Empty batch.")

    bets = []
    for line in lines[1:]:
        if line.strip() == "":
            continue
        bet = parse_bet(line.strip())
        if bet is not None:
            bets.append(bet)
    
    if len(bets) != expected:
        raise ValueError("Batch header does not match the number of bets.")
    
    return bets
