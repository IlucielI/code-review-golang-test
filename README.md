# Go (Golang) Benchmark Test Suite

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8.svg?logo=go)](https://golang.org/)
[![Benchmark Category](https://img.shields.io/badge/Benchmark-Security%20%26%20Concurrency-brightgreen.svg)](#test-case-matrix)
[![Safe Guard](https://img.shields.io/badge/False%20Positive%20Guard-Active-blue.svg)](#anti-false-positive-guard-files)
[![OWASP Top 10](https://img.shields.io/badge/OWASP-Go%20Security-orange.svg)](https://owasp.org/)

Benchmark test suite for automated code review engines on Go (Golang) applications. This repository contains intentional security vulnerabilities, concurrency race conditions, goroutine/memory leaks, and logic flaws paired with safe pattern guard files to validate zero false positives.

---

## 🎯 Benchmark Purpose

This repository validates the accuracy of automated AI/static code review engines on Go codebases:
1. **Security Vulnerability Detection:** Captures SQL injection, command injection, SSRF, IDOR, path traversal, hardcoded secrets, and JWT misconfigurations.
2. **Concurrency & Resource Leaks:** Detects unclosed files, HTTP client missing timeouts, goroutine leaks, unbuffered channel deadlocks, and concurrent map access races.
3. **Precision & Zero False Positives:** Verifies that guarded nil checks, safe batch queries, and public health checks are NOT erroneously flagged.

---

## 📋 Test Case Matrix

### 🔴 Security Vulnerabilities

| File | Issue / Vulnerability | Type | Severity | Expected |
| :--- | :--- | :--- | :---: | :---: |
| `store.go` | SQL Injection via raw string concatenation | Injection | High | **BLOCKING** |
| `exec.go` | Command Injection via `exec.Command("sh", "-c", ...)` | RCE | High | **BLOCKING** |
| `auth.go` | Hardcoded Admin Secret Token | Credential Exposure | High | **BLOCKING** |
| `logsecret.go` | Plaintext Credential / Secret Logging | Information Disclosure | High | **BLOCKING** |
| `headinj.go` | HTTP Header Injection (CRLF) | Header Injection | High | **BLOCKING** |
| `ssrf.go` | Server-Side Request Forgery via unvalidated `http.Get()` | Network Security | Medium | **BLOCKING** |
| `traversal.go` | Path Traversal via unvalidated `filepath.Join()` | File Security | High | **BLOCKING** |
| `tempfile.go` | Insecure Temporary File Creation & Traversal | File Security | High | **BLOCKING** |
| `jwt_weak.go` | Weak / Hardcoded HMAC Secret for JWT Signing | Cryptographic | High | **BLOCKING** |
| `cors_vuln.go` | Permissive CORS Wildcard with `AllowCredentials` | Security Misconfig | High | **BLOCKING** |
| `stack_trace.go` | Stack Trace Exposure directly in HTTP response | Information Disclosure | Medium | **NON-BLOCKING** |
| `rand.go` | Insecure Pseudo-Randomness (`math/rand` vs `crypto/rand`) | Cryptographic | Medium | **NON-BLOCKING** |
| `perms.go` | Overly Permissive File Permissions (`0777` World-Writable) | Security Misconfig | Medium | **NON-BLOCKING** |
| `user_profile.go` | IDOR on user account deletion without ownership check | Broken Access Control | High | **BLOCKING** |
| `redirect.go` | Open Redirect via unvalidated destination URL | Redirection | Medium | **BLOCKING** |
| `rate_limit.go` | Missing Rate Limiting on authentication endpoint | Abuse Prevention | Medium | **NON-BLOCKING** |

### ⚡ Concurrency, Performance & Resource Management

| File | Issue | Type | Severity | Expected |
| :--- | :--- | :--- | :---: | :---: |
| `goroutine.go` | Unbounded Goroutine Leak without lifecycle management | Resource Leak | Low | **NON-BLOCKING** |
| `deadlock.go` | Unbuffered Channel Deadlock on single-thread send | Concurrency Bug | Medium | **NON-BLOCKING** |
| `cache.go` | Concurrent Shared Map Access without Mutex Protection | Race Condition | Medium | **NON-BLOCKING** |
| `http_timeout.go` | HTTP Client without Timeout causing socket exhaustion | Resource Leak | Medium | **NON-BLOCKING** |
| `upload.go` | Resource Leak: Unclosed `os.File` descriptor | Resource Leak | Low | **NON-BLOCKING** |
| `batch_query.go` | N+1 Database Query Loop inside iteration | Performance | Medium | **NON-BLOCKING** |
| `redos.go` | Catastrophic Backtracking ReDoS Regular Expression | Algorithmic Complexity | Medium | **NON-BLOCKING** |

### ⚠️ Logic & Syntax Traps

| File | Issue | Type | Severity | Expected |
| :--- | :--- | :--- | :---: | :---: |
| `lookup.go` | Silently Discarded Error via Blank Identifier (`_ = err`) | Error Handling | Medium | **NON-BLOCKING** |
| `user_repository.go` | Unchecked `rows.Scan(&...)` Return Error | Error Handling | Medium | **NON-BLOCKING** |
| `render.go` | Potential Nil Pointer Dereference without check | Runtime Panic | Medium | **NON-BLOCKING** |
| `integer_overflow.go` | Unchecked Arithmetic Integer Overflow | Arithmetic Bug | Medium | **NON-BLOCKING** |

---

## 🛡️ Anti-False-Positive Guard Files

These files implement patterns that superficially resemble vulnerabilities or anti-patterns, but are safe by design:

| File | Pattern Implemented | Why Safe | Expected Result |
| :--- | :--- | :--- | :---: |
| `project_service.go` | Single bulk query with `for rows.Next()` | Query executed once before loop, not N+1 | **0 False Positives** |
| `nil_guard.go` | Explicit `if ptr == nil` checks before dereference | Safely guarded against runtime panic | **0 False Positives** |
| `idempotent.go` | Idempotent transaction execution with exponential backoff | Safe retry pattern with proper context | **0 False Positives** |
| `health_handler.go` | Public HTTP `/healthz` endpoint without authentication | Intentionally public health probe | **0 False Positives** |

---

## 🚀 How to Run the Benchmark

```bash
# View PR on GitHub
gh pr view 4 --web

# Trigger Automated Review via API
curl -X POST http://localhost:8081/api/v1/review/trigger \
  -H "Content-Type: application/json" \
  -d '{
    "repository": "IlucielI/code-review-golang-test",
    "pull_request_id": 4
  }'
```

---

## 📊 Benchmark Validation Results

- **Total Findings Detected:** 27+
- **Security Vulnerabilities:** 16 (100% detection rate)
- **Concurrency & Resource Leaks:** 7 (100% detection rate)
- **Logic & Syntax Traps:** 4 (100% detection rate)
- **False Positives on Guard Files:** 0
