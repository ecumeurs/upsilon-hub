#!/usr/bin/env python3
import sys
import json
import re
import argparse

def clean_line(line):
    # Remove prefix like [{2026-04-14T06:16:18Z}] [Bot-01] 
    line = re.sub(r'^\[\{[^}]+\}\]\s+\[[^\]]+\]\s+', '', line)
    # Remove ANSI escape codes
    line = re.sub(r'\x1b\[[0-9;]*m', '', line)
    return line.strip()

def render_ascii_board(gs):
    if not gs or 'grid' not in gs:
        return "  [No Grid Data Available]"
    
    grid = gs['grid']
    width = grid.get('width', 10)
    height = grid.get('height', 10)
    cells = grid.get('cells', [])
    
    # Map entities to positions for easy lookup
    entity_map = {}
    for p in gs.get('players', []):
        team_id = p.get('team', '?')
        nick = p.get('nickname', 'Unknown')
        for e in p.get('entities', []):
            if e.get('dead') or e.get('hp', 0) <= 0:
                continue
            pos = e.get('position')
            if pos:
                entity_map[(pos['x'], pos['y'])] = {
                    'team': team_id,
                    'name': e['name'],
                    'hp': e['hp'],
                    'max_hp': e['max_hp'],
                    'active': e.get('id') == gs.get('current_entity_id')
                }

    output = []
    header = "    " + "".join([f"{x:2}" for x in range(width)])
    output.append(header)
    output.append("    " + "--" * width)
    
    for y in range(height):
        row = [f"{y:2} |"]
        for x in range(width):
            char = " ."
            # Check for barriers/obstacles
            try:
                if y < len(cells) and x < len(cells[y]):
                    if cells[y][x].get('obstacle'):
                        char = " #"
            except: pass

            # Check for entities
            if (x, y) in entity_map:
                ent = entity_map[(x, y)]
                symbol = str(ent['team'])
                if ent['active']:
                    char = f"!{symbol}"
                else:
                    char = f" {symbol}"
            
            row.append(char)
        output.append("".join(row))
    
    output.append("\n  Legend: . Empty, # Obstacle, 1-4 Team ID, ! Active")
    return "\n".join(output)

def print_bot_summary(bot_id, data, tactical=False):
    print(f"\n>>>> STATUS FOR {bot_id} <<<<")
    
    is_finished = data.get('game_finished', False)
    winner_is_self = data.get('winner_is_self', False)
    winner_team = data.get('winner_team')

    if is_finished:
        if winner_is_self:
            print("Outcome: CONCLUDED (Winner: Self)")
        elif winner_team is not None:
            print(f"Outcome: CONCLUDED (Winner: Team {winner_team})")
    else:
        has_seen_board = data['last_board'] is not None
        survivors = [e for e in data['entities'].values() if e.get('id') in (data.get('last_board') or {}) and not e.get('dead')]
        if not survivors and has_seen_board:
            print("Outcome: CONCLUDED (ANNIHILATION)")
        elif not has_seen_board:
            print("Outcome: INCOMPLETE / NO DATA RECEIVED")
        else:
            print("Outcome: INCOMPLETE / MATCH TERMINATED")
        
    print("\nCasualties:")
    if not data['deaths']:
        print("  None recorded.")
    else:
        for name, line, owner, team in data['deaths']:
            print(f"  [L{line:05d}] {name} (Owner: {owner}, Team {team}) was eliminated.")
            if tactical:
                # Find pre-death board
                for victim_name, board in data['pre_death_boards']:
                    if victim_name == name:
                        print("\n  Board state just before elimination:")
                        print(render_ascii_board(board))
                        print("-" * 20)
                        break
            
    print("\nSurvivors:")
    last_board = data.get('last_board') or {}
    survivors = [e for e in data['entities'].values() if e.get('id') in last_board and not e.get('dead') and e.get('hp', 0) > 0]
    if not survivors:
         print("  None.")
    else:
        for s in survivors:
            print(f"  - {s['name']} (Owner: {s.get('nickname', 'AI')}, Team {s.get('team', '?')}): {s['hp']}/{s['max_hp']} HP")

    if tactical:
        print("\nRecent Tactical Actions:")
        for msg, line in data['tactical'][-15:]:
            print(f"  [L{line:05d}] {msg}")
        
        if data['last_living_board']:
            print("\nLast Known Board State:")
            print(render_ascii_board(data['last_living_board']))

    print("\nError Report:")
    if not data['errors']:
        print("  CLEAN (0 errors)")
    else:
        for msg, line in data['errors']:
            print(f"  [L{line:05d}] ERROR: {msg}")
    print("-" * 30)
    sys.stdout.flush()

def parse_log(filepath, tactical=False):
    bots = {}
    completed_bots = set()
    total_errors = 0
    
    # Regex patterns
    bot_pattern = re.compile(r'\[(Bot-\d+)\]')
    reply_pattern = re.compile(r'\[REPLY (\d+)\]')
    ws_pattern = re.compile(r'\[WS\].*board\.updated event received\.')
    cli_error_pattern = re.compile(r'(Event loop interrupted|Execution failed|timeout waiting for event):? (.*)')
    delete_marker = re.compile(r'Deleting temporary account')
    
    # Tactical Patterns for upsilon.log (Simplified one-liners)
    tactical_patterns = [
        (re.compile(r'--- (My Turn! Acting with entity: .*) ---'), r'\1'),
        (re.compile(r'Moving \d+ cells along path: (.*)'), r'Moving: \1'),
        (re.compile(r'Target in range! Attacking!'), r'Action: Attack'),
        (re.compile(r'Targeting nearest enemy: (\w+)'), r'Target: \1'),
        (re.compile(r'Ending turn with pass\.'), r'Action: Pass'),
        (re.compile(r'No enemies left\. Passing\.'), r'Action: Pass (No Enemies)'),
        (re.compile(r'Thinking\.\.\. \((\d+\.\d+)s\)'), r'Thinking: \1s'),
        (re.compile(r'Starting turn shot clock'), r'CLOCK: Started'),
        (re.compile(r'Turn timeout detected!'), r'CLOCK: TIMEOUT'),
        (re.compile(r'\[ERROR\] (.*)'), r'ERROR: \1'),
        (re.compile(r'Winner: (.*)'), r'RESULT: Winner \1'),
        (re.compile(r'VICTORY IS MINE!'), r'RESULT: VICTORY'),
        (re.compile(r'Defeated\.\.\. perishing with honor\.'), r'RESULT: DEFEAT')
    ]

    print("\n" + "="*60)
    print(" UPSILON BATTLE ENGINE DIAGNOSTIC SUMMARY")
    print("="*60)
    sys.stdout.flush()

    with open(filepath, 'r') as f:
        line_no = 0
        json_buffer = []
        in_json = False
        bracket_count = 0
        current_bot = None
        current_status_code = None
        is_ws_update = False
        last_log_was_delete = {} # bot_id -> bool
        
        for line in f:
            line_no += 1
            
            # Identify bot context from prefix
            bot_match = bot_pattern.search(line)
            if bot_match:
                bot_id = bot_match.group(1)
                
                # If we've already printed this bot, skip (unlikely unless log repeats)
                if bot_id in completed_bots:
                    continue
                    
                if bot_id != current_bot:
                    current_bot = bot_id
                    if current_bot not in bots:
                        bots[current_bot] = {
                            'entities': {},
                            'deaths': [],
                            'errors': [],
                            'tactical': [],
                            'pre_death_boards': [],
                            'game_finished': False,
                            'last_board': None,
                            'last_living_board': None
                        }
                        last_log_was_delete[current_bot] = False

            if not current_bot:
                continue

            # Check for tactical logs if enabled
            cleaned = clean_line(line)
            if tactical:
                for ptrn, fmt in tactical_patterns:
                    match = ptrn.search(cleaned)
                    if match:
                        msg = match.expand(fmt)
                        bots[current_bot]['tactical'].append((msg, line_no))
                        if getattr(args, 'filter', False):
                            print(f"[{current_bot}] {msg}")
                            sys.stdout.flush()
                        break
            
            # Detect teardown for "Finished" heuristic
            if delete_marker.search(cleaned):
                last_log_was_delete[current_bot] = True

            # Detected start of JSON (REPLY or WS)
            reply_match = reply_pattern.search(line)
            ws_match = ws_pattern.search(line)
            
            if (reply_match or ws_match) and not in_json:
                in_json = True
                json_buffer = []
                bracket_count = 0
                is_ws_update = bool(ws_match)
                if reply_match:
                    current_status_code = int(reply_match.group(1))
                    
                    # If this is a REPLY 200 following a DELETE account log, the bot is done
                    if current_status_code == 200 and last_log_was_delete.get(current_bot):
                        # Add a small note to tactical logs
                        if tactical:
                            bots[current_bot]['tactical'].append(("Agent disconnected (Cleanup successful).", line_no))
                        
                        # Process existing buffer if any (usually empty here)
                        # Then print summary
                        print_bot_summary(current_bot, bots[current_bot], tactical)
                        total_errors += count_bot_errors(bots[current_bot])
                        completed_bots.add(current_bot)
                        
                continue

            if in_json:
                if not cleaned: continue
                
                json_buffer.append(cleaned)
                bracket_count += cleaned.count('{')
                bracket_count -= cleaned.count('}')
                
                if bracket_count == 0 and '{' in "".join(json_buffer):
                    # Process JSON
                    try:
                        full_body = json.loads("".join(json_buffer))
                        
                        if is_ws_update:
                            gs = full_body.get('data', {})
                        else:
                            gs = full_body.get('data', {}).get('game_state') if full_body.get('data') else None
                            if current_status_code and current_status_code >= 400:
                                error_msg = full_body.get('message', 'Unknown Error')
                                bots[current_bot]['errors'].append((error_msg, line_no))

                        if gs:
                            process_game_state(bots[current_bot], gs, line_no, tactical)
                            
                    except Exception as e:
                        pass # JSON error or malformed
                        
                    in_json = False
                    json_buffer = []
                    continue

            # Detect CLI-level Errors
            cli_err_match = cli_error_pattern.search(line)
            if cli_err_match:
                error_msg = f"System Error: {cli_err_match.group(2)}"
                bots[current_bot]['errors'].append((error_msg, line_no))

    # Print remaining bots that didn't have a clear cleanup
    for bot_id, data in bots.items():
        if bot_id not in completed_bots:
            print_bot_summary(bot_id, data, tactical)
            total_errors += count_bot_errors(data)
            completed_bots.add(bot_id)

    print(f"\nOVERALL RESULT: {'PASS' if total_errors == 0 else 'FAIL'} ({total_errors} total errors detected)")
    print("="*60 + "\n")
    return total_errors

def count_bot_errors(data):
    errs = len(data['errors'])
    if not data.get('game_finished'):
        # Check if incomplete
        has_seen_board = data['last_board'] is not None
        survivors = [e for e in data['entities'].values() if e.get('id') in (data.get('last_board') or {}) and not e.get('dead')]
        if not (not survivors and has_seen_board) and has_seen_board:
            errs += 1
        elif not has_seen_board:
            errs += 1
    return errs

def process_game_state(bot_data, gs, line_no, tactical):
    players = gs.get('players', [])
    entities = []
    nick_map = {}
    
    for p in players:
        nick = p.get('nickname', 'Unknown')
        team = p.get('team')
        nick_map[team] = nick
        if p.get('is_self'):
            nick_map['self'] = nick
        
        for e in p.get('entities', []):
            e['nickname'] = nick
            entities.append(e)

    bot_data['nickname_map'] = nick_map
    bot_data['game_finished'] = gs.get('game_finished', False)
    bot_data['winner_is_self'] = gs.get('winner_is_self', False)
    winner_team = gs.get('winner_team_id')
    if winner_team is not None:
        bot_data['winner_team'] = winner_team

    new_entity_ids = {e['id']: e for e in entities}
    
    # Check for deaths and HP changes (Damage)
    if bot_data['last_board']:
        for old_id, old_entity in bot_data['last_board'].items():
            new_entity = new_entity_ids.get(old_id)
            
            # HP Change detection
            if new_entity and not new_entity.get('dead') and not old_entity.get('dead'):
                old_hp = old_entity.get('hp', 0)
                new_hp = new_entity.get('hp', 0)
                if new_hp < old_hp:
                    damage = old_hp - new_hp
                    if tactical:
                        bot_data['tactical'].append((f"{new_entity['name']} ({new_entity.get('nickname')}) took {damage} damage! ({new_hp} HP left)", line_no))

            # Death detection
            is_now_dead = new_entity and (new_entity.get('dead') or new_entity.get('hp') == 0)
            was_alive = not old_entity.get('dead') and old_entity.get('hp', 0) > 0
            
            if (is_now_dead and was_alive) or (old_id not in new_entity_ids and was_alive):
                bot_data['deaths'].append((old_entity['name'], line_no, old_entity.get('nickname', 'System/AI'), old_entity.get('team', '?')))
                if tactical and bot_data['last_living_board']:
                    bot_data['pre_death_boards'].append((old_entity['name'], bot_data['last_living_board']))

    bot_data['last_board'] = new_entity_ids
    bot_data['last_living_board'] = gs # Full object since we need grid
    
    for e in entities:
        bot_data['entities'][e['id']] = e

def parse_log_stream(stream, args):
    # Minimal version of parse_log for streaming
    bot_pattern = re.compile(r'\[(Bot-\d+)\]')
    # Re-use patterns from parse_log? No, let's just use the ones defined locally
    tactical_patterns = [
        (re.compile(r'--- (My Turn! Acting with entity: .*) ---'), r'\1'),
        (re.compile(r'Moving \d+ cells along path: (.*)'), r'Moving: \1'),
        (re.compile(r'Target in range! Attacking!'), r'Action: Attack'),
        (re.compile(r'Targeting nearest enemy: (\w+)'), r'Target: \1'),
        (re.compile(r'Ending turn with pass\.'), r'Action: Pass'),
        (re.compile(r'No enemies left\. Passing\.'), r'Action: Pass (No Enemies)'),
        (re.compile(r'Thinking\.\.\. \((\d+\.\d+)s\)'), r'Thinking: \1s'),
        (re.compile(r'Starting turn shot clock'), r'CLOCK: Started'),
        (re.compile(r'Turn timeout detected!'), r'CLOCK: TIMEOUT'),
        (re.compile(r'\[ERROR\] (.*)'), r'ERROR: \1'),
        (re.compile(r'Winner: (.*)'), r'RESULT: Winner \1'),
        (re.compile(r'VICTORY IS MINE!'), r'RESULT: VICTORY'),
        (re.compile(r'Defeated\.\.\. perishing with honor\.'), r'RESULT: DEFEAT')
    ]
    
    current_bot = None
    for line in stream:
        bot_match = bot_pattern.search(line)
        if bot_match:
            current_bot = bot_match.group(1)
        
        if not current_bot: continue
        
        cleaned = clean_line(line)
        for ptrn, fmt in tactical_patterns:
            match = ptrn.search(cleaned)
            if match:
                msg = match.expand(fmt)
                print(f"[{current_bot}] {msg}")
                sys.stdout.flush()
                break

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Upsilon Log Parser and Tactical Analyzer")
    parser.add_argument("logfile", nargs='?', help="Path to the .log file")
    parser.add_argument("--tactical", action="store_true", help="Enable tactical analysis and board visualization")
    parser.add_argument("--filter", action="store_true", help="Live filter mode: output tactical actions only from stdin or file")
    
    args = parser.parse_args()
    
    if args.filter:
        args.tactical = True
        input_stream = sys.stdin if not args.logfile else open(args.logfile, 'r')
        # We need to hack parse_log to accept a stream
        try:
            # Simple wrapper to use existing logic
            # Since parse_log expects a filepath, we'll re-implement a minimal version for stream
            parse_log_stream(input_stream, args)
        except KeyboardInterrupt:
            sys.exit(0)
    else:
        if not args.logfile:
            parser.print_help()
            sys.exit(1)
        try:
            parse_log(args.logfile, tactical=args.tactical)
        except FileNotFoundError:
            print(f"Error: File {args.logfile} not found.")
            sys.exit(1)
        except KeyboardInterrupt:
            print("\nInterrupted by user.")
            sys.exit(0)
