# Contributing to Dodream

Thank you for your interest in contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork
3. Run `go mod download` to install dependencies
4. Make your changes
5. Run `go build ./...` to ensure everything compiles
6. Submit a pull request

## Style Guide

### Layer Separation

This project follows a layered architecture with clear boundaries:

```
internal/
├── core/        # Domain models - no external dependencies
├── store/       # Data access layer - depends on core + ent
├── db/          # Database utilities - transaction helpers
├── application/ # Business logic - orchestrates stores
├── cli/         # CLI interface
└── engine/      # Background processing
```

**Key principles:**

- **core** must not import any other internal package
- **store** may import core and ent, but not application or cli
- **application** may import core, store, and db
- **No cyclic dependencies** between any layers

### Package Design

#### Consumer-Defined Interfaces

Interfaces are defined by the consumer, not the producer. Store packages provide concrete types only.

```go
// Good: Handler defines its own minimal interface
type UserReader interface {
    GetByID(ctx context.Context, id core.UserID) (*core.User, error)
}

// Good: Store provides concrete types
func NewUserStore(client *ent.Client) *UserStore

// Bad: Store defines and exports interfaces
// type UserRepository interface { ... }
```

#### Domain Model Conventions

All domain types in `core` follow these rules:

1. **Value objects for IDs** - Use typed string aliases, not raw strings:
   ```go
   type UserID string
   func (id UserID) String() string { return string(id) }
   ```

2. **Immutable fields** - No setters for `id` and `creator`:
   ```go
   type User struct {
       id         UserID    // immutable
       nickname   string    // mutable
       providerID string    // mutable
       createdAt  time.Time // immutable
       updatedAt  time.Time // auto-updated
   }
   ```

3. **Constructor pattern** - Use `NewXxx` functions that set timestamps:
   ```go
   func NewUser(id UserID, nickname string, providerID string) *User {
       now := time.Now()
       return &User{
           id: id, nickname: nickname, providerID: providerID,
           createdAt: now, updatedAt: now,
       }
   }
   ```

4. **Setters update timestamps** - Mutating a field updates `updatedAt`:
   ```go
   func (u *User) SetNickname(nickname string) {
       u.nickname = nickname
       u.updatedAt = time.Now()
   }
   ```

### Store Conventions

1. **No mapper files** - Convert ent types inline using `core.NewXxx` constructors
2. **Context first** - All store methods accept `context.Context` as the first parameter
3. **Error wrapping** - Use `fmt.Errorf` with `%w` to preserve context:
   ```go
   return nil, fmt.Errorf("create user: %w", err)
   ```
4. **Return domain types** - Never return `*ent.Xxx` from store methods
5. **Independent stores** - Each store is created independently with `NewXxxStore(client)`

### Transaction Handling

Use `db.WithTx` for transactions that span multiple stores:

```go
err := db.WithTx(ctx, client, func(tx *ent.Client) error {
    userStore := store.NewUserStore(tx)
    cardStore := store.NewCardStore(tx)
    // ... operations
    return nil
})
```

### Code Formatting

- Use `gofmt` for formatting
- Use `goimports` for import management
- Run `go vet ./...` before committing
- Write godoc comments for all exported types and functions

### Testing

- Use table-driven tests
- Name test functions descriptively: `TestUserStore_Create`
- Use `t.Parallel()` for independent tests
- Mock at the store level, not the ent level

## Questions?

Open an issue for discussion before making significant changes.
