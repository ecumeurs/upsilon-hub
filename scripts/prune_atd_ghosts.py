import sqlite3
import os

db_path = "docs/.atd_index.db"
if not os.path.exists(db_path):
    print(f"DB not found at {db_path}")
    exit(1)

conn = sqlite3.connect(db_path)
cursor = conn.cursor()

cursor.execute("SELECT DISTINCT file_path FROM atom_index")
rows = cursor.fetchall()

ghosts = []
for (file_path,) in rows:
    # Normalize path if needed (MCP tool paths might be relative or absolute)
    if not os.path.exists(file_path):
        ghosts.append(file_path)

print(f"Found {len(ghosts)} ghost file paths.")

if ghosts:
    print("Pruning ghosts...")
    for ghost in ghosts:
        cursor.execute("DELETE FROM atom_index WHERE file_path = ?", (ghost,))
    conn.commit()
    print("Pruning complete.")
else:
    print("No ghosts found.")

cursor.execute("SELECT COUNT(*) FROM atom_index")
print(f"Remaining entries in index: {cursor.fetchone()[0]}")

conn.close()
