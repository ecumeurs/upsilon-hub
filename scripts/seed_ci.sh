#!/bin/bash
cd "$(dirname "$0")/.."
# seed_ci.sh - Reset and seed the database for CI/E2E testing

set -e

echo "--- Resetting Database ---"
php battleui/artisan migrate:fresh --force

echo "--- Seeding Database ---"
# We run the seeders. DatabaseSeeder calls ShopItemsSeeder and SkillTemplatesSeeder.
# We pass --force because we might be in "production" environment according to Laravel.
php battleui/artisan db:seed --force

echo "--- Creating CI Admin User ---"
php battleui/artisan tinker --execute="\App\Models\User::updateOrCreate(['account_name' => 'admin'], ['email' => 'admin@admin.com', 'password_hash' => \Illuminate\Support\Facades\Hash::make('AdminPassword123!'), 'role' => 'Admin', 'full_address' => 'SYS_ADMIN_CI', 'birth_date' => '1970-01-01']);" > /dev/null 2>&1

echo "--- Creating CI Test User ---"
# Note: testuser is also seeded via TestAccountSeeder in DatabaseSeeder.php
# We use tinker here to ensure it's created even if db:seed logic changes.
php battleui/artisan tinker --execute="\App\Models\User::updateOrCreate(['account_name' => 'testuser'], ['email' => 'test@example.com', 'password_hash' => \Illuminate\Support\Facades\Hash::make('TestUserPassword123!'), 'role' => 'Player', 'full_address' => 'Test Street 1', 'birth_date' => '1990-01-01', 'credits' => 1000]);" > /dev/null 2>&1

echo "--- Seeding Complete ---"
