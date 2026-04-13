#!/usr/bin/env python3
import sys
import json
import re

def parse_log(filepath):
    bots = {}
    current_bot = None
    
    # Regex patterns
    bot_pattern = re.compile(r'\[(Bot-\d+)\]')
    reply_pattern = re.compile(r'\[REPLY (\d+)\]')
    ws_pattern = re.compile(r'\[WS\] board\.updated event received\.')
    
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
                    'last_board': None
                }
        
        if not current_bot:
            i += 1
            continue

        # Detect Errors
        reply_match = reply_pattern.search(line)
        if reply_match:
            status_code = int(reply_match.group(1))
            if status_code >= 400:
                # Try to parse JSON message
                error_msg = "Unknown Error"
                try:
                    # Look for message in the next few lines (JSON body)
                    json_lines = []
                    j = i + 1
                    while j < i + 10 and j < len(lines):
                        json_lines.append(lines[j].split(']', 1)[-1].strip() if ']' in lines[j] else lines[j].strip())
                        if '}' in lines[j]: break
                        j += 1
                    
                    body_str = "".join(json_lines)
                    # Clean up ANSI codes if present
                    body_str = re.sub(r'\x1b\[[0-9;]*m', '', body_str)
                    body = json.loads(body_str)
                    error_msg = body.get('message', error_msg)
                except:
                    pass
                bots[current_bot]['errors'].append((error_msg, line_no))

        # Detect Board Updates (Deaths and Winner)
        if ws_pattern.search(line):
            try:
                json_lines = []
                j = i + 1
                bracket_count = 0
                started = False
                while j < len(lines):
                    # Multi-level prefix stripping (e.g., "[Bot-01] [WS]   {")
                    content = lines[j]
                    while ']' in content and (content.strip().startswith('[') or content.strip()[0:1] == '\x1b'):
                        content = content.split(']', 1)[-1]
                    
                    content = content.strip()
                    content = re.sub(r'\x1b\[[0-9;]*m', '', content) # Remove colors
                    json_lines.append(content)
                    
                    bracket_count += content.count('{')
                    bracket_count -= content.count('}')
                    if '{' in content: started = True
                    if started and bracket_count == 0: break
                    j += 1
                
                board_data = json.loads("".join(json_lines))
                # board.updated event structure is { match_id: ..., data: { entities: [], ... } }
                grid_data = board_data.get('data', {})
                entities = grid_data.get('entities', [])
                
                # Check for winner
                winner = grid_data.get('winner_id')
                if winner:
                    bots[current_bot]['winner'] = winner
                
                # Check for deaths
                new_entity_ids = {e['id']: e for e in entities}
                if bots[current_bot]['last_board']:
                    for old_id, old_entity in bots[current_bot]['last_board'].items():
                        if old_id not in new_entity_ids:
                            bots[current_bot]['deaths'].append((old_entity['name'], line_no, old_entity.get('nickname', 'System/AI')))
                
                bots[current_bot]['last_board'] = new_entity_ids
                # Update current entity status
                for e in entities:
                    bots[current_bot]['entities'][e['id']] = e
                
                i = j # Skip processed JSON
            except Exception as e:
                # print(f"Error parsing board: {e}")
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
        return 0

    for bot_id, data in bots.items():
        print(f"\n>>>> STATUS FOR {bot_id} <<<<")
        
        # Winner
        if data['winner']:
            print(f"Outcome: CONCLUDED (Winner: {data['winner']})")
        else:
            print("Outcome: INCOMPLETE / MATCH TERMINATED")
            
        # Deaths
        print("\nCasualties:")
        if not data['deaths']:
            print("  None recorded.")
        else:
            for name, line, owner in data['deaths']:
                print(f"  [L{line:05d}] {name} ({owner}) was eliminated.")
                
        # Survivors
        print("\nSurvivors:")
        survivors = [e for e in data['entities'].values() if e.get('id') in data.get('last_board', {})]
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
