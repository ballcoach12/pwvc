# Data Model: JWT Authentication and Authorization System

**Feature**: JWT Authentication and Authorization System  
**Date**: October 28, 2025  
**Branch**: 001-jwt-auth-system

## Entity Definitions

### User Entity (Extended)

**Purpose**: Represents a system user with authentication credentials and profile information

**Attributes**:

- `id` (uint): Unique identifier, primary key, auto-generated
- `username` (string): Unique username for login, required, max 50 characters
- `password_hash` (string): bcrypt hashed password, required, max 255 characters
- `email` (string): User email address, optional, max 100 characters
- `created_at` (timestamp): Account creation time, auto-generated
- `updated_at` (timestamp): Last modification time, auto-updated
- `is_active` (boolean): Account status flag, default true

**Relationships**:

- Many-to-many with Role through UserRole junction table
- One-to-many with audit logs (existing)

**Validation Rules**:

- Username must be unique across system
- Username must be 3-50 characters, alphanumeric plus underscore/dash
- Password must be hashed using bcrypt with cost factor 12
- Email must be valid format if provided
- is_active defaults to true for new accounts

**State Transitions**:

- Created → Active (default state)
- Active → Inactive (admin action or account suspension)
- Inactive → Active (admin reactivation)

### Role Entity (New)

**Purpose**: Represents permission levels and access groups for users

**Attributes**:

- `id` (uint): Unique identifier, primary key, auto-generated
- `name` (string): Role name, required, unique, max 50 characters
- `description` (string): Human-readable role description, optional, max 255 characters
- `permissions` (JSON/text): Serialized permissions data, optional
- `created_at` (timestamp): Role creation time, auto-generated
- `updated_at` (timestamp): Last modification time, auto-updated

**Relationships**:

- Many-to-many with User through UserRole junction table

**Validation Rules**:

- Role name must be unique across system
- Role name must be 2-50 characters, lowercase letters and underscores only
- Standard roles: "admin", "user", "viewer" (predefined)
- Description maximum 255 characters if provided

**Predefined Roles**:

- `admin`: Full system access, user management, all features
- `user`: Standard user access, project participation, basic features
- `viewer`: Read-only access, view projects and results only

### UserRole Entity (New)

**Purpose**: Junction table managing many-to-many relationship between users and roles

**Attributes**:

- `id` (uint): Unique identifier, primary key, auto-generated
- `user_id` (uint): Foreign key to users table, required
- `role_id` (uint): Foreign key to roles table, required
- `assigned_at` (timestamp): When role was assigned, auto-generated
- `assigned_by` (uint): User ID who assigned the role, optional foreign key

**Relationships**:

- Belongs to User (via user_id)
- Belongs to Role (via role_id)
- Belongs to User (via assigned_by, optional)

**Validation Rules**:

- Composite unique constraint on (user_id, role_id) - no duplicate role assignments
- user_id must reference existing user
- role_id must reference existing role
- assigned_by must reference existing user if provided

### JWT Token Entity (Logical)

**Purpose**: Represents authentication token structure (not stored in database)

**Claims Structure**:

- `sub` (string): Subject - user ID
- `username` (string): Username for display purposes
- `roles` ([]string): Array of role names for authorization
- `iat` (int64): Issued at timestamp
- `exp` (int64): Expiration timestamp
- `iss` (string): Issuer identifier (application name)

**Validation Rules**:

- Token must be signed with HS256 algorithm
- Expiration must be within 1-24 hours of issuance
- Roles must be valid role names from roles table
- Subject must reference existing active user

## Database Schema

### GORM Auto-Migration (SQLite)

Since the application uses SQLite with GORM AutoMigrate, the database schema will be automatically managed by GORM based on the struct definitions. No manual SQL migrations are needed.

#### Schema Creation via GORM Models

```go
// GORM will automatically create these tables based on the struct definitions:
// - roles table from Role struct
// - user_roles table from UserRole struct
// - Add password_hash and is_active columns to existing users table

// Initial data seeding will be handled in the application startup:
func seedRoles(db *gorm.DB) error {
    roles := []domain.Role{
        {Name: "admin", Description: "Administrator with full system access"},
        {Name: "user", Description: "Standard user with project participation rights"},
        {Name: "viewer", Description: "Read-only access to view projects and results"},
    }

    for _, role := range roles {
        db.FirstOrCreate(&role, domain.Role{Name: role.Name})
    }
    return nil
}
```

#### GORM AutoMigrate Integration

The database schema will be automatically updated in `cmd/server/main.go`:

```go
// Auto-migrate all models including new auth entities
err = db.AutoMigrate(
    &domain.Project{},
    &domain.Attendee{},
    &domain.Feature{},
    // ... existing models
    &domain.User{},        // Extended with password_hash, is_active
    &domain.Role{},        // New model
    &domain.UserRole{},    // New junction table
)
```

### Domain Model Relationships

```
User (1) ←→ (N) UserRole (N) ←→ (1) Role
  ↓
AuditLog (existing)
```

## GORM Model Definitions

### Go Struct Definitions

```go
// User model (extended)
type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
    PasswordHash string    `gorm:"size:255" json:"-"` // Never serialize password
    Email        string    `gorm:"size:100" json:"email,omitempty"`
    IsActive     bool      `gorm:"default:true" json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    // Relationships
    Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// Role model (new)
type Role struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
    Description string    `gorm:"size:255" json:"description,omitempty"`
    Permissions string    `gorm:"type:text" json:"permissions,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // Relationships
    Users []User `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// UserRole junction model (new)
type UserRole struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    UserID     uint      `gorm:"not null" json:"user_id"`
    RoleID     uint      `gorm:"not null" json:"role_id"`
    AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
    AssignedBy *uint     `json:"assigned_by,omitempty"`

    // Relationships
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// JWT Claims structure (logical)
type JWTClaims struct {
    UserID   uint     `json:"sub"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}
```

## Data Access Patterns

### User Repository Methods

- `CreateUser(user *User) error`
- `GetUserByUsername(username string) (*User, error)`
- `GetUserByID(id uint) (*User, error)`
- `UpdateUser(user *User) error`
- `ListUsers(offset, limit int) ([]User, error)`
- `LoadUserRoles(user *User) error`

### Role Repository Methods

- `GetAllRoles() ([]Role, error)`
- `GetRoleByName(name string) (*Role, error)`
- `CreateRole(role *Role) error`
- `AssignRoleToUser(userID, roleID uint, assignedBy *uint) error`
- `RemoveRoleFromUser(userID, roleID uint) error`

### Query Optimization

- Index on users.username for login queries
- Index on user_roles(user_id) for role lookup
- Composite index on user_roles(user_id, role_id) for uniqueness
- Preload roles when fetching users to avoid N+1 queries

## Security Considerations

### Password Storage

- Never store plain text passwords
- Use bcrypt with cost factor 12 minimum
- Hash passwords before database storage
- Never include password_hash in JSON responses

### Token Security

- JWT tokens contain only necessary claims
- Short expiration times (1-24 hours)
- Signed with strong secret key
- Validate all claims on each request

### Data Validation

- Validate input at API boundaries
- Sanitize usernames and role names
- Check role assignments before token generation
- Audit all authentication and authorization events
