# mathlang

mathlang is a **minimal experimental language** for exploring linear algebra concepts through executable structure rather than symbolic manipulation alone.

mathlang demo files are note-first.
Only lines matching supported statement patterns are interpreted; all others are ignored.

Its core design goal is **learning-first**:
- Use LaTeX-like mathematical notation as input
- Strip away presentation details early
- Build a small but semantically meaningful AST
- Execute linear algebra ideas directly (not via pre-baked matrix formulas)

This project is intentionally *not* a full LaTeX parser, nor a general CAS. It is a scaffold for understanding how linear algebra objects become computable structures.

---

## Design Philosophy

1. **Math objects first, syntax second**  
   Vectors, bases, and operations exist as AST nodes. Strings are only a transport layer.

2. **Ignore presentation noise**  
   Display math wrappers such as `\\[` and `\\]` are stripped in preprocessing. They carry no semantic meaning.

3. **Small executable demos**  
   Each demo implements exactly one linear algebra idea end-to-end.

4. **Fail loudly**  
   Undefined symbols, invalid constructions, or dimension errors trigger immediate panics. This is a learning language, not a forgiving one.

---

## Current Capabilities (Demo1)

- Vector assignment using LaTeX-like syntax
- Basis definition from previously defined vectors
- Explicit evaluation request (eval)
- Execution via numeric linear algebra (Gaussian elimination)

Example input fragment:

```
\[
\vec{b}_1 = \begin{pmatrix}1\\2\end{pmatrix},\quad
\vec{b}_2 = \begin{pmatrix}3\\4\end{pmatrix}
\]

\[
\vec{v} = \begin{pmatrix}1\\1\end{pmatrix}
\]

\[
b = \{\vec{b}_1,\vec{b}_2\}
\]

\[
[\vec{v}]_b \leftarrow \text{eval}
\]
```

---

## Architecture Overview

```
Input (LaTeX-like lines)
   ↓
Lexer (line-based, semantic classification)
   ↓
Parser (AST construction)
   ↓
Evaluator (numeric execution)
```

Key point: **there is no full grammar**. Lines are classified by intent, not parsed by precedence rules.

---

## Non-Goals (by design)

- Full LaTeX compliance
- Symbolic algebra
- Automatic simplification
- Implicit operations

These may appear in later demos, but are intentionally excluded now.

---

## Roadmap (Conceptual)

- Demo2: Linear transformations as composable operations
- Demo3: Emergence of matrices from repeated transformations
- Demo4: Dimension checks, linear independence, symbolic reasoning

---

## Status

This project is **experimental and educational**. Expect rough edges. That is the point.

