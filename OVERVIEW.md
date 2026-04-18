# Distributed Job Queue — Project Overview

## What I'm Building

A distributed job queue system — think a simplified version of Sidekiq, Celery, or AWS SQS + workers. Producers push jobs onto a queue, workers pull them off and execute them, and the system handles failures, retries, scheduling, and scaling workers horizontally.

This is **not** a personal task manager. It's backend infrastructure — the kind of system that sits behind a web app to handle background work (sending emails, processing images, running ML inference, etc.).

## Why I'm Building It

Portfolio project #2, explicitly chosen for **system design learning**. My first project (Field Ops at Synovion) shows I can ship a product end-to-end. This second project is meant to prove I can reason about distributed systems concepts that come up in real backend interviews and engineering work.

## Core Concepts I Want to Learn / Demonstrate

- **Producer / consumer architecture** — decoupling job submission from execution
- **Persistence & durability** — jobs survive crashes (Redis or Postgres as the backing store)
- **At-least-once delivery** — visibility timeouts, acknowledgements, redelivery on worker death
- **Retries with exponential backoff** — transient failure handling
- **Dead letter queue** — jobs that keep failing get parked, not lost
- **Priority queues** — urgent jobs jump the line
- **Scheduled / delayed jobs** — "run this in 10 minutes" or "run this at 3pm"
- **Worker scaling** — multiple workers pulling concurrently without double-processing
- **Observability** — metrics for queue depth, job latency, failure rate
- **Idempotency** — jobs that can be safely retried

## Rough Feature List (MVP → Stretch)

### MVP
- Enqueue a job (JSON payload + job type)
- Workers dequeue and execute
- Job status tracking (queued / running / done / failed)
- Retries with backoff
- Simple CLI or HTTP API to submit jobs

### Phase 2
- Priority levels
- Delayed / scheduled jobs
- Dead letter queue
- Basic dashboard (queue depth, recent jobs, worker status)

### Stretch
- Multi-queue routing
- Rate limiting per queue
- Distributed tracing for job lifecycle
- Horizontal scaling story (how to add workers without downtime)

## Tech Stack (Tentative — Decide in New Session)

Open questions to nail down at the start of the next chat:

- **Language**: Go (great for concurrency, common in infra) vs Python (faster to prototype) vs Java/Spring (matches my resume skills) vs Rust (performance flex)
- **Backing store**: Redis (classic choice, fast, simple) vs Postgres (ACID, SKIP LOCKED pattern) vs both
- **Transport**: Direct DB polling vs pub/sub vs gRPC between components
- **Dashboard**: Simple HTML + htmx, or React, or skip entirely for MVP

My leaning: **Go + Redis** as the main stack, because it's the most "infrastructure-y" combo and will force me to learn goroutines, channels, and Redis data structures (LIST, SORTED SET, streams). But I'm open to Java/Spring Boot since that's what I already know — the learning would then be more about the distributed systems patterns than the language.

## What "Done Enough to Showcase" Looks Like

- Clean README with architecture diagram
- Docker Compose to spin up everything locally (queue service + Redis/Postgres + worker + demo producer)
- A short write-up / blog post explaining the design decisions (why visibility timeouts, why exponential backoff, etc.)
- Live demo deployed somewhere (Fly.io, Railway, or a small VPS) that recruiters can actually poke at
- Load test results — "handles X jobs/sec on a single worker, scales linearly to Y"

## Why This Matters for My Portfolio

Field Ops shows product thinking and full-stack execution. This project is meant to show:

1. I can design a system from scratch, not just assemble a framework's pieces
2. I understand the hard parts of distributed systems (not just CRUD)
3. I can write clean backend code, not just ship features
4. I can explain tradeoffs — which is what system design interviews actually test

## Starting Point for the New Claude Chat

When I open the new chat in this folder, I want to:

1. Pick the stack (language + storage) — have a real discussion about tradeoffs
2. Sketch the architecture (components, data flow, what talks to what)
3. Design the data model (job schema, queue schema, worker registry if needed)
4. Scaffold the project (folder structure, Dockerfile, docker-compose, basic Makefile)
5. Implement the MVP path end-to-end before adding features

Don't let me skip step 2. I want to actually think through the design before writing code.
