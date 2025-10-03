Agent Brief: Build pg_faux — an in-memory Postgres-compatible test DB for Ruby
Objective
Deliver a small, fast, framework-agnostic in-memory database that accepts a subset of Postgres SQL and can stand in for SQLite :memory: during tests. It must be usable by any Ruby app via a PG-like client, and optionally by ActiveRecord through a thin adapter. No real Postgres process, no wire protocol, no Docker.

Core Outcomes
Schema ingestion (dynamic): Load db/structure.sql at connect time using pg_query (libpg_query) to build an in-memory catalog: tables, columns, types, defaults, PK/UNIQUE/FK, sequences, enums, indexes (metadata only).

Executor (subset): Interpret a strict subset of SQL against in-memory tables:

Supported v1: single-table SELECT (columns/expressions pass-through, simple WHERE predicates with =, AND), ORDER BY, LIMIT/OFFSET; INSERT (with DEFAULT/sequence), UPDATE (set literals/params), DELETE; BEGIN/COMMIT/ROLLBACK.
Not supported v1: joins, CTEs, subqueries, window functions, aggregates, DDL at runtime, COPY. Fail fast with NotImplemented that names the unsupported construct.
Client surface (PG-like): Provide a Ruby object that duck-types PG::Connection for exec, exec_params, prepare, exec_prepared, and transaction methods. No global monkey-patching by default. Enable via an explicit factory or an env flag.

(Optional) ActiveRecord adapter: A thin adapter that maps AR calls to the client so projects can set adapter: pg_faux in database.yml for test. Keep this layer separate; the core must not depend on Rails.

Determinism & safety: Catalog objects are immutable after load; transactions are copy-on-write snapshots; constraints (NOT NULL, UNIQUE, FK) and simple type checks are enforced. Errors are precise and actionable.

Non-Goals (v1)
No Postgres wire protocol or socket server.
No full SQL planner/optimizer.
No binary format/OIDs compatibility layer beyond what AR’s ActiveRecord::Result needs.
No silent fallbacks or regex-based SQL parsing.
Performance & Footprint Targets
Parse & build catalog for a typical structure.sql (≤1 MB) in ≤50 ms per process.
Query execution overhead for supported subset: ≤1 ms for small tables (<1k rows).
Memory: ≤50 MB total for typical test suites; catalog kept frozen and re-used.
High-Level Architecture
Use hexagonal architecture (ports & adapters) and CQS (Command-Query Separation).

Domain (pure):

Catalog: value objects for schema (tables, columns, constraints, enums, sequences). Immutable/frozen. Provides lookups and metadata.
Engine: interprets a parsed SQL AST and executes against in-memory tables with constraints, sequences, and transactional snapshots.
Result Model: minimal rowset abstraction (fields, to_a, counts) that duck-types what typical Ruby code expects from PG::Result.
Adapters:

DDL Loader (input): parses structure.sql with pg_query → builds Catalog. Include a schema cache layer (Msgpack) keyed by content hash to skip parsing in hot paths.
PG-like Client (output): provides exec, exec_params, prepare, exec_prepared, transaction methods; routes to Engine.
AR Adapter (optional): implements execute, exec_query, begin/commit/rollback, minimal quoting/types; turns Engine results into ActiveRecord::Result.
Support:

Stub Layer (optional): allow exact/regex SQL stubs to return deterministic results for cases outside supported SQL. Stubs are opt-in and explicitly registered in tests.
Error system: typed errors (ValidationError, ConstraintViolation, NotImplemented, TypeMismatch) with context (table, column, sql fragment).
Clean Code Standards (what to follow)
Immutability for shared state: Catalog and type descriptors are frozen; Engine’s transactional state uses snapshotting. Avoid global singletons; pass dependencies explicitly.
Small, cohesive modules: Split by responsibility: Catalog, Validator, Executor, Transaction, Codec, Client, ARAdapter.
Explicitness over magic: For unsupported features, raise immediately with a one-line reason and, when relevant, suggest: “use a stub or a real PG for this test”.
Pure functions where possible: Validators and predicate evaluators are pure and easy to fuzz.
Guard rails: “Strict mode” default—unknown columns, implicit casts, or ambiguous operations error out; “Permissive mode” only by explicit flag for legacy tests.
Dependency boundaries: The core depends only on pg_query and standard library; schema cache may add msgpack. The AR adapter is an optional gem or sub-package.
Testable seams: The PG-like client is a thin facade; the Engine accepts AST + params; the DDL loader returns a plain Catalog.
No reflection-heavy metaprogramming; no monkey-patches outside the optional “enable via ENV” test harness.
Naming: use domain names (Catalog, Table, Column, Constraint, Sequence) and avoid overloading “Model/Record”.
Things to Avoid (foot-guns)
Regex SQL parsing. Always use pg_query.
Generating many Ruby files from schema. Prefer a single Msgpack cache or one frozen hash literal if you must.
Silent coercions (e.g., comparing string “1” to integer 1). Enforce strict typing where catalog says so.
Implicit joins or partial join emulation. Either implement correctly later or fail fast.
Global mutable state (e.g., class variables holding connection state). Keep connections instance-scoped.
Hidden monkey-patching of PG::Connection. If provided, it must be behind an explicit flag and only in test.
SQL Support Matrix (v1)
SELECT: column list (* allowed), single table; WHERE with = on columns/params/literals; AND chains; ORDER BY simple columns; LIMIT/OFFSET. No functions/operators beyond basic comparisons in v1.
INSERT: column list + values (params/literals); DEFAULT handling; RETURNING * or column list allowed if the row exists in memory. Sequences auto-increment on PK if configured.
UPDATE/DELETE: single table + WHERE with =; enforce constraints post-update.
Transactions: BEGIN, COMMIT, ROLLBACK using copy-on-write snapshots.
Constraints: NOT NULL, UNIQUE, FK (same-process referential checks), PK uniqueness, basic default expressions (e.g., nextval, constants).
Types: map common Postgres types to Ruby scalars (int4/int8/numeric/boolean/text/uuid/date/timestamp/timestamptz/jsonb). Jsonb stored as Ruby Hash/Array without operators in v1 (validate presence, not behavior).
Error & Validation Rules
Validate table/column existence, arity of INSERT, unknown RETURNING fields, NOT NULL on missing values, UNIQUE before commit, FK on insert/update/delete.
Errors must include: operation, table, column(s) involved, and a compact slice of the offending SQL.
Concurrency & Isolation
Single-process, thread-safe by coarse-grained mutex around Engine state. Transaction snapshots are per-client. No cross-process sharing.
Configuration
Connection factory accepts:

schema_path: (string) or schema_sql: (string)
use_cache: (bool, default true)
mode: (:strict | :permissive)
seed: initial table data for tests (hash of arrays)
Environment flag PG_FAUX=1 may enable a PG monkey-patch only in test (optional).

Testing Strategy (the agent must implement)
Contract tests for the PG-like client surface (exec[_params], prepared statements, tx).
Catalog loader tests: feed representative structure.sql samples (enums, sequences, constraints) → assert catalog shape.
Behavioral engine tests: golden tests for SELECT/INSERT/UPDATE/DELETE/transactions/constraints.
Fuzz tests (property-based): generate small random schemas and query combinations within the supported subset → assert determinism and invariants (uniques, FKs).
Negative tests: unsupported features must raise NotImplemented with clear messages.
(Optional) Cross-check: a nightly job (behind a flag) compares supported-subset results against a real Postgres container on a tiny schema to catch semantic drift. Not required for local dev.
Deliverables
Gem: pg_faux (core) with zero Rails deps.

Optional gem: activerecord-pg_faux-adapter (separate package).

Docs:

README with quick start, support matrix, limitations, flags.
“When to use stubs vs engine” guide.
“Migrating from SQLite :memory:” guide.
CI: GitHub Actions running lint + tests + fuzz suite; optional nightly cross-check.

Versioning: Semantic versioning; start at 0.1.0. Changelog required.

Implementation Order (strict)
Catalog loader via pg_query + schema cache (Msgpack, content-hash keyed).
Engine data model (in-memory tables) + sequences + constraints.
Executor for INSERT + SELECT (no ORDER/LIMIT yet) + transactions.
Add UPDATE/DELETE, ORDER BY/LIMIT, and RETURNING.
PG-like client, prepared statements.
Error taxonomy + strict/permissive modes.
Optional AR adapter.
Stub layer.
Benchmarks + docs.
Quality Bars (acceptance)
Can run a sample non-Rails Ruby repo’s repository tests with pg_faux instead of SQLite/PG.
Optional: a small Rails app can run basic model CRUD and transactions in test with adapter: pg_faux.
Parsing structure.sql ≤50 ms; simple SELECT on 1k rows ≤1 ms median.
Test suite covers: positive, negative, fuzz; code coverage for core ≥85% (line), with meaningful branch checks on validator.
call it instead rubock
