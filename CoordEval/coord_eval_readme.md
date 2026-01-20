# CoordEval

CoordEval is **Demo1** of the mathlang project.

It implements a single, concrete linear algebra task:

> **Given a vector `v` and a basis `b`, compute the coordinate representation `[v]_b`.**

Nothing more. Nothing less.

---

## Mathematical Meaning

CoordEval evaluates the equation:

```
B · x = v
```

Where:
- `B` is the matrix whose columns are the basis vectors
- `v` is a vector in standard coordinates
- `x` is the coordinate vector of `v` relative to basis `b`

The output is `x`.

This demo intentionally treats matrices as *derived objects*, not primitives.

---

## Supported Statements

### 1. Vector Assignment

```
\vec{b}_1 = \begin{pmatrix}1\\2\end{pmatrix},\quad
\vec{b}_2 = \begin{pmatrix}3\\4\end{pmatrix}
```

- Vectors may optionally have numeric subscripts
- Only integer components are supported
- Dimensionality is inferred from component count

---

### 2. Basis Definition

```
b = \{\vec{b}_1,\vec{b}_2\}
```

- All vectors must be defined before use
- Order matters (defines column order)
- Basis dimension = number of vectors

---

### 3. Evaluation Request

```
[\vec{v}]_b \leftarrow \text{eval}
```

- Requests computation of coordinates of `v` in basis `b`
- Only one evaluation per file is supported

---

## Execution Model

1. Parse input into an AST
2. Resolve symbols (vectors, basis)
3. Construct basis matrix column-by-column
4. Solve the linear system using Gaussian elimination

There is **no symbolic manipulation** at this stage.

---

## Error Handling

CoordEval is intentionally strict:

- Undefined vectors → panic
- Undefined basis → panic
- Empty basis → panic
- Dimension mismatch → runtime failure

This makes semantic errors visible immediately.

---

## What CoordEval Is *Not*

- A matrix calculator
- A symbolic algebra system
- A general linear algebra library

It is a **conceptual bridge** between mathematical notation and executable structure.

---

## Why This Demo Matters

CoordEval exists to demonstrate one idea:

> *Coordinates are not numbers attached to vectors; they are solutions to a linear system defined by a basis.*

Once this is executable, matrices can emerge naturally rather than being assumed.

---

## Status

CoordEval is complete as Demo1.
Future demos will build on its AST and execution model.

