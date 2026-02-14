# Contributing to terraform-provider-opnsense

Thank you for your interest in contributing to the OPNsense Terraform provider!

## Development Environment Setup

### Prerequisites

- Go >= 1.21
- Terraform >= 1.5
- Git
- Access to an OPNsense instance for testing

### Setting Up Your Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/terraform-provider-opnsense.git
   cd terraform-provider-opnsense
   ```

3. Install dependencies:
   ```bash
   make deps
   ```

4. Build the provider:
   ```bash
   make build
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Make your changes

3. Format your code:
   ```bash
   make fmt
   ```

4. Run tests:
   ```bash
   make test
   ```

5. Build the provider to ensure it compiles:
   ```bash
   make build
   ```

### Testing

#### Unit Tests

Unit tests should be added for all new functionality:

```go
func TestMyNewFunction(t *testing.T) {
    // Test implementation
}
```

Run unit tests with:
```bash
make test
```

#### Manual Testing

1. Install the provider locally:
   ```bash
   make install
   ```

2. Create a test Terraform configuration in a separate directory
3. Run Terraform commands to test your changes

### Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Add comments for exported functions and types
- Keep functions small and focused

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, etc.)
- Reference issue numbers when applicable

Example:
```
Add support for DHCP static mappings

Implements #123
```

## Adding New Resources

When adding a new resource:

1. Create a new file in `internal/resources/`
2. Implement the resource interface with CRUD operations
3. Register the resource in `internal/provider/provider.go`
4. Add tests in a `*_test.go` file
5. Update the import tool in `tools/import/main.go` to discover the new resource
6. Add documentation and examples

## Pull Request Process

1. Ensure all tests pass
2. Update documentation as needed
3. Add examples for new features
4. Submit a pull request with a clear description of the changes
5. Wait for review and address any feedback

## Questions or Problems?

If you have questions or run into problems, please open an issue on GitHub.

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project.
