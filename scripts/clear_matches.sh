#!/usr/bin/env bash
cd "$(dirname "$0")/.."
# Clear ghost matches from the database

echo "Clearing matches..."
cd /workspace/battleui || exit 1
php artisan tinker --execute="DB::statement('TRUNCATE table game_matches CASCADE'); DB::statement('TRUNCATE table match_participants CASCADE'); echo 'Matches and participants cleared.'; echo \"\n\";"
echo "Done!"
