# Contributing to RnR Core

Thank you for your interest in contributing to RnR Core! We welcome contributions from developers, researchers, and blockchain enthusiasts. As an **educational testbed** and **research platform**, we value code quality, clear documentation, and honest trade-offs.

## 🤝 How to Contribute

### 1. Reporting Bugs
- **Check Exising Issues**: Before submitting, please check if the issue has already been reported.
- **Use the Template**: When opening an issue, please use our default template providing:
  - System specs (OS, CPU, RAM)
  - Go version
  - Exact command run
  - Logs/Panic output

### 2. Suggesting Features
- We use a **Proposal Process** for major features.
- Open an Issue with the tag `[PROPOSAL]`
- Describe the feature, the problem it solves, and the technical approach.
- Wait for feedback from maintainers before writing code.

### 3. Submitting Pull Requests
1. **Fork** the repository.
2. **Clone** your fork locally.
3. Create a **Feature Branch** (`git checkout -b feature/amazing-feature`).
4. **Commit** your changes with clear messages (`git commit -m 'feat: implement new sorting algo'`).
5. **Push** to your fork.
6. Open a **Pull Request**.

---

## 🛠️ Development Setup

### Prerequisites
- **Go**: Version 1.21 or higher
- **Git**: Latest version
- **Make** (Optional): For running build scripts

### Local Build
```bash
# clone repo
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE

# install dependencies
go mod download

# build binary
go build -o ./bin/rnr-node ./cmd/rnr-node

# run tests
go test ./...
```

---

## 📐 Coding Standards

### Go Style
- We follow effective Go guidelines.
- ALWAYS run `go fmt ./...` before committing.
- Variable names should be descriptive (e.g., `blockHeight` not `h`).
- Exported functions MUST have comments.

### Architecture Rules
- `internal/` package is for private implementation details.
- `pkg/` is for public APIs (if any).
- **No Circular Dependencies**: Be careful with imports.
- **Interfaces**: Use interfaces to decouple components (e.g., Consensus Engine should be an interface).

---

## 🧪 Testing Guidelines

- **Unit Tests**: Required for all new logic.
- **Table-Driven Tests**: Preferred for algorithm verification.
- **Simulations**: For network-level changes, update or create a simulation in `simulation/`.

Example Table-Driven Test:
```go
func TestSort(t *testing.T) {
    tests := []struct{
        input    []byte
        expected []byte
    }{
        {[]byte{3,1,2}, []byte{1,2,3}},
        {[]byte{}, []byte{}},
    }
    // ... loop and assert
}
```

---

## ⚖️ License
By contributing, you agree that your contributions will be licensed under the MIT License.
