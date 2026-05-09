#!/usr/bin/env python3
# @spec-link [[rule_code_health_monitoring]]
# This script enforces the project's "Zero Error" code health standards, including
# file length, function complexity, and intent-based documentation requirements.
# Note: this is highly specific for our project: it will only work for the files
# and frameworks we use (Go, PHP, JS, Vue, Python).
import os
import re
import sys

# Constants
EXTENSIONS = {'.go', '.py', '.php', '.js', '.vue'}
IGNORE_DIRS = {'vendor', 'node_modules', '.git', 'dist', 'build'}

# Thresholds
LOC_WARN = 300
LOC_ERROR = 500
NESTING_MAX = 3
# Documentation policy: Preceding comments only, ATD tags excluded.
ATD_MIN = 2
ATD_WARN_MAX = 5
ATD_ERROR_MAX = 10

class HealthCheck:
    def __init__(self):
        self.errors = 0
        self.warnings = 0
        self.valid_atds = self._load_atd_ids()

    def _load_atd_ids(self):
        atd_ids = set()
        docs_dir = 'docs'
        if not os.path.exists(docs_dir):
            return atd_ids
        for root, _, files in os.walk('.'):
            if 'docs' in root:
                for file in files:
                    if file.endswith('.atom.md'):
                        path = os.path.join(root, file)
                        with open(path, 'r', encoding='utf-8', errors='ignore') as f:
                            content = f.read()
                            match = re.search(r'^id:\s*([a-zA-Z0-9_-]+)', content, re.MULTILINE)
                            if match:
                                atd_ids.add(match.group(1))
        return atd_ids

    def check_file(self, filepath):
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            lines = f.readlines()
        content = "".join(lines)
        ignore_bloating = '@lint-ignore-file-bloating' in content
        ignore_complexity = '@lint-ignore-complexity' in content
        ignore_docs = '@lint-ignore-documentation' in content
        ignore_atd = '@lint-ignore-atd' in content

        print(f"Checking {filepath}...")
        if not ignore_bloating:
            loc = len(lines)
            if loc > LOC_ERROR:
                self.report(filepath, "ERROR", f"File too long: {loc} LOC (limit {LOC_ERROR})")
            elif loc > LOC_WARN:
                self.report(filepath, "WARN", f"File long: {loc} LOC (limit {LOC_WARN})")

        is_test_file = any(t in filepath for t in ['_test.go', '.spec.', '.test.']) or '/tests/' in filepath.lower()
        atd_links_matches = re.findall(r'@(spec|test)-link\s+\[\[([a-zA-Z0-9_:-]+)\]\]', content)
        atd_links = [m[1] for m in atd_links_matches]
        
        if not ignore_atd:
            atd_count = len(atd_links)
            if atd_count < ATD_MIN:
                self.report(filepath, "ERROR", f"Too few ATD links: {atd_count} (min {ATD_MIN})")
            elif atd_count > ATD_ERROR_MAX:
                self.report(filepath, "ERROR", f"Too many ATD links: {atd_count} (max {ATD_ERROR_MAX})")
            elif atd_count > ATD_WARN_MAX:
                self.report(filepath, "WARN", f"Many ATD links: {atd_count} (limit {ATD_WARN_MAX})")
            
            # Enforce @test-link in test environment
            if is_test_file:
                spec_in_test = [m[1] for m in atd_links_matches if m[0] == 'spec']
                if spec_in_test:
                    self.report(filepath, "ERROR", f"Test file uses @spec-link instead of @test-link for: {spec_in_test}")

            for atd_id in atd_links:
                # If it's a cross-project link (contains ':'), we only check the local part if it's the current project
                # or we skip it if we don't have full workspace context.
                # For simplicity in this script, we skip phantom check for cross-project links.
                if ':' in atd_id:
                    continue
                if atd_id not in self.valid_atds:
                    self.report(filepath, "ERROR", f"Phantom ATD link: [[{atd_id}]] does not exist")

        if not ignore_complexity or not ignore_docs:
            self.check_functions(filepath, lines, ignore_complexity, ignore_docs)

    def check_functions(self, filepath, lines, ignore_complexity, ignore_docs):
        func_start_re = re.compile(r'^\s*(?:func|def|function)\s+(?:\([^*)]*(?:\*?[a-zA-Z0-9_]+)?\)\s+)?([a-zA-Z0-9_]+)|^\s*([a-zA-Z0-9_]+)\s*[:=]\s*(?:\(.*\)|[a-zA-Z0-9_]+)?\s*=>')
        is_python = filepath.endswith('.py')
        
        current_func = None
        func_preceding = []
        func_body = []
        depth = 0
        max_depth = 0
        preceding_comments = []
        
        for line in lines:
            clean_line = line.strip()

            if not clean_line:
                if current_func:
                    func_body.append(line)
                continue

            if clean_line.startswith('//') or clean_line.startswith('#') or clean_line.startswith('/*') or clean_line.startswith('*'):
                if current_func:
                    func_body.append(line)
                else:
                    preceding_comments.append(line)
                continue

            match = func_start_re.search(line)
            if match:
                if current_func:
                    self.analyze_func(filepath, current_func, func_preceding, func_body, max_depth, ignore_complexity, ignore_docs)
                current_func = match.group(1) or match.group(2)
                func_preceding = preceding_comments
                func_body = [line]
                preceding_comments = []
                if is_python:
                    depth = 1
                else:
                    depth = line.count('{') - line.count('}')
                    if depth < 1: depth = 1 
                max_depth = depth
                continue

            if current_func:
                func_body.append(line)
                if is_python:
                    if clean_line.endswith(':'):
                        depth += 1
                else:
                    depth += clean_line.count('{')
                    depth -= clean_line.count('}')
                
                nest_keywords = ['if ', 'if(', 'for ', 'for(', 'while ', 'while(', 'switch ', 'switch(', 'case ', 'select ', 'catch ', 'try ']
                current_line_nesting = 0
                for kw in nest_keywords:
                    if kw in clean_line:
                        current_line_nesting = 1
                        break
                
                if depth + current_line_nesting > max_depth:
                    max_depth = depth + current_line_nesting
                
                if not is_python and depth <= 0:
                    self.analyze_func(filepath, current_func, func_preceding, func_body, max_depth, ignore_complexity, ignore_docs)
                    current_func = None
                    func_body = []
                    func_preceding = []
            else:
                preceding_comments = []

        if current_func:
            self.analyze_func(filepath, current_func, func_preceding, func_body, max_depth, ignore_complexity, ignore_docs)

    def analyze_func(self, filepath, name, preceding, body, max_depth, ignore_complexity, ignore_docs):
        if not ignore_complexity:
            if max_depth > NESTING_MAX + 1:
                self.report(filepath, "ERROR", f"Function '{name}' too complex: nesting {max_depth-1} (limit {NESTING_MAX})")

        if not ignore_docs:
            # Filter out ATD tags from preceding comments
            real_comments = []
            for l in preceding:
                l_strip = l.strip()
                if (l_strip.startswith('//') or l_strip.startswith('#') or l_strip.startswith('/*') or l_strip.startswith('*')) and \
                   not re.search(r'@(spec|test)-link', l_strip):
                    # Remove comment markers to check if there is actual text
                    clean_comment = re.sub(r'^(\/\/|#|\/\*|\*|\s)+', '', l_strip).strip()
                    if clean_comment:
                        real_comments.append(clean_comment)
            
            # Determine if exposed
            is_exposed = False
            if filepath.endswith('.go'):
                is_exposed = name[0].isupper()
            else:
                # Basic check for non-private or explicit keywords
                signature = body[0] if body else ""
                if 'public' in signature or 'export' in signature or not name.startswith('_'):
                    is_exposed = True

            if not real_comments:
                level = "ERROR" if is_exposed else "WARN"
                msg = f"Function '{name}' missing preceding documentation."
                if is_exposed:
                    msg += " Public/Exposed functions require clear documentation of intent, inputs, and outputs."
                self.report(filepath, level, msg)

    def report(self, filepath, level, message):
        if level == "ERROR":
            self.errors += 1
            print(f"\033[91m[ERROR]\033[0m {filepath}: {message}")
        else:
            self.warnings += 1
            print(f"\033[93m[WARN]\033[0m {filepath}: {message}")

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="Upsilon Code Health Check")
    parser.add_argument("path", nargs="?", default=".", help="File or directory to check (default: current directory)")
    args = parser.parse_args()
    check = HealthCheck()
    if os.path.isfile(args.path):
        check.check_file(args.path)
    else:
        for root, dirs, files in os.walk(args.path):
            dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
            for file in files:
                ext = os.path.splitext(file)[1]
                if ext in EXTENSIONS:
                    check.check_file(os.path.join(root, file))
    print("\n" + "="*40)
    print(f"Health Check Summary:")
    print(f"Errors: {check.errors}")
    print(f"Warnings: {check.warnings}")
    print("="*40)
    print("Doc Policy: Preceding comments required. ATD tags excluded. Body docs discouraged.")
    print("="*40)
    if check.errors > 0:
        sys.exit(1)
    sys.exit(0)
