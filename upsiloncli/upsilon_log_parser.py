#!/usr/bin/env python3
import sys
import json
import re

def clean_line(line):
    # Remove prefix like [{2026-04-14T06:16:18Z}] [Bot-01] 
    line = re.sub(r'^\[\{[^}]+\}\]\s+\[[^\]]+\]\s+', '', line)
    # Remove ANSI escape codes
    line = re.sub(r'\x1b\[[0-9;]*m', '', line)
    return line.strip()

def parse_log(filepath):
    bots = {}
    current_bot = None
    
    # Regex patterns
    bot_pattern = re.compile(r'\[(Bot-\d+)\]')
    reply_pattern = re.compile(r'\[REPLY (\d+)\]')
    ws_pattern = re.compile(r'\[WS\].*board\.updated event received\.')
    cli_error_pattern = re.compile(r'(Event loop interrupted|Execution failed|timeout waiting for event):? (.*)')
    
    with open(filepath, 'r') as f:
        lines = f.readlines()
        
    i = 0
    while i < len(lines):
        line = lines[i]
        line_no = i + 1
        
        # Identify bot context from prefix
        bot_match = bot_pattern.search(line)
        if bot_match:
            current_bot = bot_match.group(1)
            if current_bot not in bots:
                bots[current_bot] = {
                    'entities': {},
                    'deaths': [],
                    'errors': [],
                    'winner': None,
                    'winner_team': None,
                    'last_board': None
                }
        
        if not current_bot:
           i += 1
           continue

        # Detect REPLYs (Both Errors and Success)
        reply_match = reply_pattern.search(line)
        if reply_match:
            status_code = int(reply_match.group(1))
            
            # --- Handle JSON Body for both Errors and Success ---
            try:
                json_lines = []
                j = i + 1
                bracket_count = 0
                started = False
                while j < i + 10000 and j < len(lines):
                    cleaned = clean_line(lines[j])
                    if not cleaned:
                        j += 1
                        continue
                    
                    json_lines.append(cleaned)
                    bracket_count += cleaned.count('{')
                    bracket_count -= cleaned.count('}')
                    if '{' in cleaned: started = True
                    if started and bracket_count == 0: break
                    j += 1
                
                body_str = "".join(json_lines)
                full_body = json.loads(body_str)
                
                if status_code >= 400:
                    error_msg = full_body.get('message', 'Unknown Error')
                    bots[current_bot]['errors'].append((error_msg, line_no))
                else:
                    # Success reply - extract game state if present
                    data = full_body.get('data', {})
                    # Standardized for API: Board state is in data.game_state
                    if data and 'game_state' in data:
                        gs = data['game_state']
                        players = gs.get('players', [])
                        
                        entities = []
                        nick_map = {}
                        for p in players:
                            t = p.get('team')
                            nick = p.get('nickname', 'Unknown')
                            nick_map[t] = nick
                            if p.get('is_self'):
                                nick_map['self'] = nick
                            
                            for e in p.get('entities', []):
                                e['nickname'] = nick
                                entities.append(e)
                        
                        bots[current_bot]['nickname_map'] = nick_map
                        bots[current_bot]['game_finished'] = gs.get('game_finished', False)
                        bots[current_bot]['winner_is_self'] = gs.get('winner_is_self', False)
                        
                        winner_team = gs.get('winner_team_id')
                        if winner_team is not None:
                            bots[current_bot]['winner_team'] = winner_team
                        
                        new_entity_ids = {e['id']: e for e in entities}
                        bots[current_bot]['last_board'] = new_entity_ids
                        for e in entities:
                            bots[current_bot]['entities'][e['id']] = e
                i = j # Skip processed JSON
            except:
                pass

        # Detect CLI-level Errors
        cli_err_match = cli_error_pattern.search(line)
        if cli_err_match:
            error_msg = f"System Error: {cli_err_match.group(2)}"
            bots[current_bot]['errors'].append((error_msg, line_no))

        # Detect Board Updates (Deaths and Winner)
        if ws_pattern.search(line):
            try:
                json_lines = []
                j = i + 1
                bracket_count = 0
                started = False
                while j < len(lines):
                    cleaned = clean_line(lines[j])
                    if not cleaned:
                        j += 1
                        continue
                        
                    json_lines.append(cleaned)
                    bracket_count += cleaned.count('{')
                    bracket_count -= cleaned.count('}')
                    if '{' in cleaned: started = True
                    if started and bracket_count == 0: break
                    j += 1
                
                board_data = json.loads("".join(json_lines))
                # Standardized for WebSocket: board state is in the 'data' field of the envelope
                gs = board_data.get('data', {})
                players = gs.get('players', [])
                
                # Flatten
                entities = []
                for p in players:
                    nick = p.get('nickname', 'Unknown')
                    for e in p.get('entities', []):
                        e['nickname'] = nick
                        entities.append(e)
                
                # Check for winner
                bots[current_bot]['game_finished'] = gs.get('game_finished', False)
                bots[current_bot]['winner_is_self'] = gs.get('winner_is_self', False)
                
                winner_team = gs.get('winner_team_id')
                if winner_team is not None:
                    bots[current_bot]['winner_team'] = winner_team
                
                # Check for deaths
                new_entity_ids = {e['id']: e for e in entities}
                
                if bots[current_bot]['last_board']:
                    for old_id, old_entity in bots[current_bot]['last_board'].items():
                        new_entity = new_entity_ids.get(old_id)
                        # Detect death either by transition to 'dead' flag or HP reaching 0
                        is_now_dead = new_entity and (new_entity.get('dead') or new_entity.get('hp') == 0)
                        was_alive = not old_entity.get('dead') and old_entity.get('hp', 0) > 0
                        
                        if is_now_dead and was_alive:
                            bots[current_bot]['deaths'].append((old_entity['name'], line_no, old_entity.get('nickname', 'System/AI')))
                        
                        # Fallback for old behavior (removal)
                        if old_id not in new_entity_ids:
                            bots[current_bot]['deaths'].append((old_entity['name'], line_no, old_entity.get('nickname', 'System/AI')))
                
                bots[current_bot]['last_board'] = new_entity_ids
                # Update current entity status
                for e in entities:
                    bots[current_bot]['entities'][e['id']] = e
                
                i = j # Skip processed JSON
            except:
                pass
        
        i += 1
        
    return bots

def print_summary(bots):
    total_errors = 0
    print("\n" + "="*60)
    print(" UPSILON BATTLE ENGINE DIAGNOSTIC SUMMARY")
    print("="*60)
    
    if not bots:
        print("No bot data found in log.")
        return 1

    for bot_id, data in bots.items():
        print(f"\n>>>> STATUS FOR {bot_id} <<<<")
        
        is_finished = data.get('game_finished', False)
        winner_is_self = data.get('winner_is_self', False)
        winner_team = data.get('winner_team')

        # Outcome
        if is_finished:
            if winner_is_self:
                print("Outcome: CONCLUDED (Winner: Self)")
            elif winner_team is not None:
                print(f"Outcome: CONCLUDED (Winner: Team {winner_team})")
            else:
                print("Outcome: CONCLUDED (Winner: Unknown)")
        else:
            has_seen_board = data['last_board'] is not None
            survivors = [e for e in data['entities'].values() if e.get('id') in (data.get('last_board') or {})]
            if not survivors and has_seen_board:
                print("Outcome: CONCLUDED (ANNIHILATION)")
            elif not has_seen_board:
                print("Outcome: INCOMPLETE / NO DATA RECEIVED")
                total_errors += 1
            else:
                print("Outcome: INCOMPLETE / MATCH TERMINATED")
                total_errors += 1
            
        # Character Status (Deaths and Survivors)
        print("\nCasualties:")
        if not data['deaths']:
            print("  None recorded.")
            print("  [CRITICAL] Error: No deaths occurred (symptomatic of failure).")
            total_errors += 1
        else:
            for name, line, owner in data['deaths']:
                print(f"  [L{line:05d}] {name} ({owner}) was eliminated.")
                
        # Survivors
        print("\nSurvivors:")
        last_board = data.get('last_board') or {}
        survivors = [e for e in data['entities'].values() if e.get('id') in last_board and not e.get('dead') and e.get('hp', 0) > 0]
        if not survivors:
             print("  None.")
        else:
            for s in survivors:
                print(f"  - {s['name']} (Owner: {s.get('nickname', 'AI')}): {s['hp']}/{s['max_hp']} HP")
        
        # Errors
        print("\nError Report:")
        if not data['errors']:
            print("  CLEAN (0 errors)")
        else:
            for msg, line in data['errors']:
                print(f"  [L{line:05d}] ERROR: {msg}")
                total_errors += 1
        print("-" * 30)

    print(f"\nOVERALL RESULT: {'PASS' if total_errors == 0 else 'FAIL'} ({total_errors} total errors detected)")
    print("="*60 + "\n")
    return total_errors

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 upsilon_log_parser.py <path_to_log>")
        sys.exit(1)
    
    results = parse_log(sys.argv[1])
    err_count = print_summary(results)
    sys.exit(0 if err_count == 0 else 1)
