# Copilot Instructions for PWVC

## Project Overview

This is a web application that facilitates the use of the Pairwise Weighted Value Comparison (PWVC) technique. The PWVC technique is a decision-making method that helps users compare and prioritize options through pairwise comparisons with weighted values.

## Technology Stack

- **Language**: Go
- **Type**: Web Application
- **License**: MIT

## Development Guidelines

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Follow the Go Code Review Comments guide
- Write clear, self-documenting code with meaningful variable and function names

### Project Structure

- Keep the codebase organized and modular
- Separate concerns between business logic, handlers, and data models
- Use appropriate Go project layout (e.g., cmd/, internal/, pkg/ directories as the project grows)

### Testing

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Aim for good test coverage of critical paths
- Tests should be runnable with `go test ./...`

### Documentation

- Document exported functions, types, and packages with Go doc comments
- Keep the README.md up to date with setup instructions and usage examples
- Document any environment variables or configuration needed

### Dependencies

- Minimize external dependencies when possible
- Use Go modules for dependency management
- Keep dependencies up to date and review them for security issues

### Git Workflow

- Write clear, descriptive commit messages
- Keep commits focused and atomic
- Reference issue numbers in commits when applicable

## Specific Patterns to Follow

### Error Handling

- Always handle errors appropriately
- Don't ignore errors silently
- Return errors to callers rather than panicking in most cases
- Use error wrapping with `fmt.Errorf` and `%w` for context

### Web Development

- Use proper HTTP status codes
- Implement appropriate middleware for logging, recovery, etc.
- Validate user input thoroughly
- Follow REST principles for API design if building a REST API

### Security

- Sanitize all user inputs
- Use HTTPS in production
- Implement proper authentication and authorization as needed
- Don't commit secrets or credentials to the repository (use .env files)

## Build and Run

When the project includes source code:
- Build: `go build`
- Test: `go test ./...`
- Run: Follow instructions in README.md

## Special Considerations

- The PWVC technique involves comparing pairs of options, so data structures should support efficient pairwise comparisons
- The UI should make it easy for users to input options and perform comparisons
- Consider storing comparison results in a structured format for later analysis
