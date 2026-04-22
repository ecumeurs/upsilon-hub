Based on my investigation of issue ISS-062, I've identified critical silent failure defaults across configuration and communication layers that violate the crash-early principle and break strict API enforcement.

CRITICAL FINDINGS

1. Configuration Layer Violations

battleui/config/services.php:38-41
'upsilon' => [ 
'url' => env('UPSILON_API_URL'),// ❌ Returns null if missing
'webhook_url' => env('UPSILON_WEBHOOK_URL'), // ❌ Returns null if missing
],
- Impact: UpsilonApiService receives null baseUrl, causing HTTP requests to fail with cryptic errors
- Fix: Add validation or throw exception on missing mandatory URLs
 
upsiloncli/cmd/upsiloncli/main.go:43-46
defaultBaseURL := os.Getenv("UPSILON_BASE_URL")
if defaultBaseURL == "" {
defaultBaseURL = "http://localhost:8000" // ❌ Hides missing config
}
- Impact: CLI attempts to connect to localhost in production when config is missing
- Fix: Follow the REVERB_APP_KEY pattern (lines 35-41) - exit with clear error 

battleui/.env:85-86
UPSILON_API_URL=http://localhost:8081# ❌ Hardcoded dev URL 
UPSILON_WEBHOOK_URL=http://127.0.0.1:8000/api/webhook/upsilon
- Impact: Production deployments may accidentally use localhost URLs 

2. API Communication Violations

battleui/app/Services/UpsilonApiService.php:108,113-120
return $response->json() ?? []; // ❌ Line 108: Silent empty array on failure

catch (\Exception $e) {
Log::error("Upsilon API Connection Failed..."); 
return [// ❌ Lines 113-120: Fake envelope instead of exception
'request_id' => $requestId,
'message' => 'Connection to Game Engine failed',
'success' => false,
'data' => [],
'meta' => ['exception' => $e->getMessage()]
];
}
- VIOLATION: This breaks [[api_standard_envelope]] contract by returning fake envelopes
- Impact: Frontend receives success: false but treats it as valid response, hiding real connection failures
- Fix: Throw UpsilonConnectionException that Laravel's global handler can convert to proper error response
 
3. Frontend Ghost State Issues 
 
battleui/resources/js/Pages/BattleArena.vue:146
const grid = computed(() =>
gameState.value?.grid || { width: 10, height: 10, cells: [] } // ❌ Ghost 10x10 board
); 
- Impact: Users see empty 10x10 grid when API fails, thinking the game loaded
- Fix: Show loading/error state, never render ghost board
 
battleui/resources/js/services/tactical.js:81,101,121 
return me ? me.entities || [] : []; // ❌ Silent empty arrays
- Impact: Missing entities silently vanish instead of triggering errors

4. WebSocket Connection Failures

battleui/resources/js/services/game.js:56-65
if (!window.Echo) {
console.error('Laravel Echo is not initialized.');
return; // ❌ Silent return, game continues without WebSocket
}
- Impact: Real-time updates fail silently, game appears frozen 

📋 PRIORITY FIXES

IMMEDIATE (Crash Early violations):
1. Remove silent returns in UpsilonApiService - throw exceptions
2. Remove ghost grid default in BattleArena.vue
3. Add config validation for UPSILON_API_URL on Laravel boot

HIGH (API Contract violations):
4. CLI should exit on missing UPSILON_BASE_URL (follow REVERB_APP_KEY pattern)
5. Remove all || [] defaults in tactical services
6. WebSocket init failures should prevent game start 
 
MEDIUM (Configuration hygiene):
7. Remove hardcoded localhost URLs from .env
8. Add environment-specific validation in Laravel AppServiceProvider

The core issue: Silent defaults mask configuration and communication failures, making debugging impossible in production. The fix must enforce crash-early behavior especially
for API communication that must strictly follow [[api_standard_envelope]].
