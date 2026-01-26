### Mathlang Conventions

#### 1. Naming Conventions

**Directory / Module Names**

* Use lowercase or snake_case
* No CamelCase in directory names

Examples:

```
demo1_coord_evval/
demo2_linear_combo_transform/
```

**Go Code Naming**

* Types / Structs / Interfaces: `PascalCase`
* Exported functions: `PascalCase`
* Internal variables / functions: `camelCase`
* Constants / enums: `PascalCase` or `SCREAMING_SNAKE`

**Math Objects**

* Mathematical concepts (Vec, Basis, Transform) are always modeled as types
* Mathematical names should not be embedded in directory names

---

#### 2. Layering Rules

This project follows a strict separation of layers:

```
Input Text
   ↓
Lexer
   ↓
Parser
   ↓
AST  ← semantic boundary
   ↓
Eval / Compute
   ↓
Output
```

Rules:

* Lexer only tokenizes, no semantic knowledge
* Parser builds AST, no math computation
* AST stores *intent*, not results
* Eval operates only on AST

---

#### 3. Matrix Policy (Critical)

> Matrix is a derived representation, not a primitive.

Rules:

* Matrices must NOT appear in AST
* Matrices must NOT be required for semantic understanding
* Matrices may be constructed only as a final or derived representation

Linear combinations and mappings are primary.
Matrices are secondary artifacts.

---

#### 4. Demo Policy

* Each demo may duplicate code intentionally
* Abstraction into shared libraries is deferred until concepts stabilize
* Clarity of mathematical meaning is prioritized over DRY