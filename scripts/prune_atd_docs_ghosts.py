import sqlite3
import os

db_path = "docs/.atd_docs_index.db"
if not os.path.exists(db_path):
    print(f"DB not found at {db_path}")
    exit(1)

conn = sqlite3.connect(db_path)
cursor = conn.cursor()

cursor.execute("SELECT id, file_path FROM atom_docs_index")
rows = cursor.fetchall()

ghosts = []
for (aid, file_path) in rows:
    if not file_path or not os.path.exists(file_path):
        ghosts.append(aid)

print(f"Found {len(ghosts)} ghost atoms in docs_index.")

if ghosts:
    print("Pruning ghosts...")
    for aid in ghosts:
        cursor.execute("DELETE FROM atom_docs_index WHERE id = ?", (aid,))
    conn.commit()
    print("Pruning complete.")
else:
    print("No ghosts found.")

cursor.execute("SELECT COUNT(*) FROM atom_docs_index")
print(f"Remaining atoms in index: {cursor.fetchone()[0]}")

conn.close()
