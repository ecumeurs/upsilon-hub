Setting up a robust Continuous Integration (CI) pipeline for a multi-service architecture like Upsilon Battle is a massive milestone. It is completely normal to feel like you are staring at a jungle right now—especially when dealing with real-time WebSockets, Go engines, and Laravel gateways. 

The good news is that your team has already built a remarkably strong foundation. You have the `upsiloncli`, the Goja JS engine, and even a log parser. You correctly identified a gap, though, which is actually documented in your open issues as **ISS-044: Unified Scripting Lifecycle and CI Testing Framework**.

Here is a comprehensive, structured plan to formalize your testing, wrangle your Docker containers, and keep your GitHub Actions clean.

### 1. Assessing Your Toolkit (What You Have vs. What You Need)

You actually have almost all the moving parts required to execute these tests. 

* **Execution:** You have `upsiloncli --farm <scripts>` to run the bots.
* **Assertions:** You have `upsilon.assert(condition, message)` inside your JS scripts to halt execution on failure.
* **Validation:** You have `upsilon_log_parser.py` and `tests/run_all_battles.sh` which are designed to fail if protocol violations are detected.
* **What you are missing:** A formal test manifest (as you mentioned) and a CI-specific Docker Compose file. Your `docker-compose.prod.yaml` is a great starting point, but production files are optimized for persistence and security, whereas CI files must be optimized for speed and ephemeral data.

---

### 2. Formalizing the Test Specification

To prevent your CI from becoming a mess of unreadable bash scripts, we need a declarative way to define what a test is. I recommend creating a YAML-based "Test Manifest" system. You can write a lightweight Python or Bash script (e.g., `tests/runner.py`) that reads these YAML files and executes the CLI accordingly.

Here is a proposed structure for a test manifest (e.g., `tests/manifests/2v2_pvp_basic.yaml`):

```yaml
name: "2v2 PVP Standard Engagement"
intent: "Verify that 4 players can successfully matchmake, enter an arena, and execute basic moves without protocol errors."
timeout_seconds: 120

execution:
  command: "./bin/upsiloncli"
  flags: 
    - "--farm"
    - "--logs=./ci_logs/2v2_pvp"
  scripts:
    - "tests/bots/2v2_player_alpha.js"
    - "tests/bots/2v2_player_beta.js"
    - "tests/bots/2v2_player_gamma.js"
    - "tests/bots/2v2_player_delta.js"

sanction:
  # How we measure success/failure
  require_exit_code: 0
  post_run_checks:
    - command: "python3 upsiloncli/upsilon_log_parser.py ci_logs/2v2_pvp/ --filter"
      must_not_contain: "ERROR"
```

---

### 3. Adapting Docker Compose for CI

GitHub Actions handles multi-container setups beautifully using Docker Compose. However, we need to make a `docker-compose.ci.yaml` tailored for testing:

* **Database:** Use a temporary in-memory database (`tmpfs`) for PostgreSQL to drastically speed up database seeding and migrations. We don't need to persist volume data (`db_data`) between CI runs.
* **Health Checks:** This is critical. GitHub Actions needs to know exactly when Laravel, Reverb, and the Go Engine are fully booted and ready to accept API calls before it fires the `upsiloncli` scripts. 

Your workflow will essentially do this:
1. Checkout code.
2. `docker compose -f docker-compose.ci.yaml up -d`
3. Wait for containers to report "healthy".
4. Run migrations/seeders.
5. Execute the test manifests.

---

### 4. Organizing the CI to Avoid the "Jungle"

Do not put everything into a single `main.yml` file. Break your CI into distinct, purposeful workflows located in `.github/workflows/`:

1.  **`lint-and-build.yml`:** Runs on every push. Checks Go syntax, PHP formatting, builds the UI, and builds `upsiloncli` to ensure it compiles.
2.  **`unit-tests.yml`:** Runs your standard Go tests and PHPUnit tests.
3.  **`e2e-battles.yml`:** This is your heavy lifter. It spins up the Docker Compose stack, runs the `run_all_battles.sh` suite, and processes your new Test Manifests. 

By keeping the heavy integration tests in a separate file, you can easily configure them to only run on Pull Requests to the `main` branch, saving you CI minutes.

---

### 5. Actionable Checklist for Implementation

* [ ] **Create `docker-compose.ci.yaml`:** Strip out persistent volumes from the prod file and add `healthcheck` blocks to the `app`, `ws`, and `engine` services.
* [ ] **Integrate `check_services.sh`:** You already have this utility. Use it in your GitHub Action step right after `docker compose up -d` to block the pipeline until all endpoints return HTTP 200.
* [ ] **Standardize the Teardown:** Ensure every JS bot script strictly implements `upsilon.onTeardown()` to delete accounts and forfeit matches. CI environments get messy if previous test runs leave ghost data in the database.
* [ ] **Draft the YAML Test Runner:** Write a simple script (Python is usually easiest for YAML parsing) that reads the manifests we designed above and executes the `upsiloncli` commands.


---
---
---

Good call on addressing the `.env` strategy early. Managing environment variables is usually where CI pipelines break first, especially when services are trying to talk to each other across Docker networks.

Since this is an automated testing environment, **you do not need real secrets**. Your CI database is ephemeral, and your API keys only need to be mathematically valid, not cryptographically secure. 

Here is the exact strategy for handling your environment variables and the `docker-compose.ci.yaml` file to pull it all together.

### 1. The CI `.env` Strategy

Instead of managing GitHub Secrets for basic test variables, commit a dedicated **`.env.ci`** file directly into your repository. This guarantees that anyone (or any bot) running the tests has the exact same configuration. 

In your GitHub Action, your first step will simply be to copy this file:
`cp .env.ci .env`

Here is what your `.env.ci` should look like. Notice that we are pointing the services to the Docker hostnames (like `db` and `engine`), but we are also setting the `UPSILON_BASE_URL` to `localhost` so your `upsiloncli` running on the GitHub runner can reach inside the containers:

```env
# .env.ci
APP_ENV=testing
APP_DEBUG=true
# A dummy key for Laravel Sanctum
APP_KEY=base64:c29tZV9yYW5kb21fc3RyaW5nX2Zvcl90ZXN0aW5nX2tleQ== 

# Ephemeral DB settings
DB_CONNECTION=pgsql
DB_HOST=db
DB_PORT=5432
DB_DATABASE=upsilon_test
DB_USERNAME=postgres
DB_PASSWORD=postgres

# Upsilon Engine
UPSILON_API_URL=http://engine:8081/internal

# Reverb / WebSockets
REVERB_APP_ID=ci_app_id
REVERB_APP_KEY=ci_app_key
REVERB_APP_SECRET=ci_app_secret
REVERB_SERVER_HOST=0.0.0.0
REVERB_SERVER_PORT=8080

# CLI Configuration (Running on the host, targeting the exposed ports)
UPSILON_BASE_URL=http://localhost:8000
REVERB_HOST=127.0.0.1:8080
```

---

### 2. The `docker-compose.ci.yaml`

This file is optimized purely for speed and reliability. We are dropping persistent volumes in favor of `tmpfs` (RAM disks), and we are adding `healthcheck` blocks so GitHub Actions knows exactly when the stack is ready.

```yaml
# docker-compose.ci.yaml
services:
  db:
    image: postgres:18-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: upsilon_test
    ports:
      - "5432:5432"
    # Use tmpfs instead of volumes for lightning-fast I/O during tests
    tmpfs:
      - /var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d upsilon_test"]
      interval: 2s
      timeout: 5s
      retries: 10

  app:
    build:
      context: ./battleui
    ports:
      - "8000:80"
    depends_on:
      db:
        condition: service_healthy
    env_file:
      - .env
    healthcheck:
      # Waits for Laravel to respond
      test: ["CMD", "curl", "-f", "http://localhost/up"]
      interval: 5s
      timeout: 5s
      retries: 5

  ws:
    build:
      context: ./battleui
    command: php artisan reverb:start --host=0.0.0.0 --port=8080
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    env_file:
      - .env
    healthcheck:
      # Reverb responds to a simple HTTP GET on its port
      test: ["CMD", "curl", "-f", "http://localhost:8080"]
      interval: 5s
      timeout: 5s
      retries: 5

  engine:
    build:
      context: .
      dockerfile: ./upsilonapi/Dockerfile
    ports:
      - "8081:8081"
    depends_on:
      db:
        condition: service_healthy
    env_file:
      - .env
    healthcheck:
      # Assuming you have a health endpoint, otherwise check the port
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 5s
      timeout: 5s
      retries: 5
```

---

### 3. How This Looks in GitHub Actions

With those health checks in place, your GitHub Action syntax becomes incredibly clean. You use the `--wait` flag, which forces the runner to pause until all health checks report as `healthy`.

```yaml
    steps:
      - uses: actions/checkout@v4
      
      - name: Prepare Environment
        run: cp .env.ci .env

      - name: Boot Upsilon Services
        run: docker compose -f docker-compose.ci.yaml up -d --wait
        
      - name: Run Laravel Migrations & Seeders
        run: docker compose -f docker-compose.ci.yaml exec app php artisan migrate:fresh --seed
        
      - name: Build Upsilon CLI
        run: go build -o bin/upsiloncli ./cmd/upsiloncli

      # ... execute the test runner here
```

This ensures we aren't relying on arbitrary `sleep 10` commands in our CI pipeline—it only moves forward when the Go engine, Reverb websocket server, and Laravel gateway are truly ready.

