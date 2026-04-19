import os
import re

atoms_dir = "docs"
code_dirs = ["upsilonapi", "upsiloncli", "upsilonbattle", "battleui", "scripts"]

ids = []
atom_data = {}

# 1. Get all STABLE orphans from disk
for filename in os.listdir(atoms_dir):
    if filename.endswith(".atom.md"):
        path = os.path.join(atoms_dir, filename)
        with open(path, "r") as f:
            content = f.read()
            match = re.search(r"---(.*?)---", content, re.DOTALL)
            if match:
                fm = match.group(1)
                id_match = re.search(r"id:\s*(.*)", fm)
                status_match = re.search(r"status:\s*(.*)", fm)
                
                aid = id_match.group(1).strip() if id_match else None
                astatus = status_match.group(1).strip() if status_match else "UNKNOWN"
                
                if aid and astatus == "STABLE":
                    ids.append(aid)

# 2. Search for tags in code
implementations = {}
findings = {} # aid -> list of files where found

for aid in ids:
    implementations[aid] = False
    findings[aid] = []

for d in code_dirs:
    if not os.path.exists(d): continue
    for root, _, files in os.walk(d):
        for file in files:
            if file.endswith((".go", ".js", ".vue", ".php", ".py")):
                path = os.path.join(root, file)
                try:
                    with open(path, "r") as f:
                        c = f.read()
                        for aid in ids:
                            if f"[[{aid}]]" in c:
                                implementations[aid] = True
                                findings[aid].append(path)
                except:
                    pass

# 3. Separate into False Orphans and True Orphans
false_orphans = [aid for aid in ids if implementations[aid]]
true_orphans = [aid for aid in ids if not implementations[aid]]

print("--- FALSE ORPHANS (Tagged but reported as orphan) ---")
for aid in false_orphans:
    print(f"- {aid}")
    for path in findings[aid]:
        print(f"  FOUND IN: {path}")

print("\n--- TRUE ORPHANS (Missing tags) ---")
for aid in true_orphans:
    print(f"- {aid}")
