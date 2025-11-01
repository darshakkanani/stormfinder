# 🤝 Contributing to Stormfinder

Thank you for your interest in contributing to Stormfinder! This document provides guidelines and information for contributors.

## 🎯 **How to Contribute**

### **Types of Contributions**
- 🐛 **Bug Reports**: Help us identify and fix issues
- ✨ **Feature Requests**: Suggest new capabilities
- 🔧 **Code Contributions**: Submit pull requests
- 📚 **Documentation**: Improve guides and examples
- 🧪 **Testing**: Help test new features and releases
- 🌐 **Translations**: Localize for different languages

## 🚀 **Getting Started**

### **Development Setup**
```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/stormfinder.git
cd stormfinder

# Install dependencies
go mod download

# Build the project
go build ./cmd/stormfinder

# Run tests
go test ./...
```

### **Development Environment**
- **Go**: Version 1.24 or higher
- **Git**: For version control
- **Make**: For build automation (optional)
- **Docker**: For containerized testing (optional)

## 📝 **Contribution Guidelines**

### **Code Style**
- Follow standard Go conventions (`gofmt`, `golint`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small
- Use consistent error handling

### **Commit Messages**
Follow the conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(ai): add neural network subdomain prediction
fix(bruteforce): resolve memory leak in worker pool
docs(readme): update installation instructions
```

### **Pull Request Process**

1. **Fork the Repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/stormfinder.git
   ```

2. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Your Changes**
   - Write clean, well-documented code
   - Add tests for new functionality
   - Update documentation as needed

4. **Test Your Changes**
   ```bash
   # Run all tests
   go test ./...
   
   # Run the comprehensive test suite
   ./scripts/test_all_features.sh
   
   # Test build process
   ./scripts/build.sh
   ```

5. **Commit and Push**
   ```bash
   git add .
   git commit -m "feat(scope): your descriptive message"
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**
   - Use the PR template
   - Provide clear description
   - Link related issues
   - Add screenshots if applicable

## 🐛 **Bug Reports**

### **Before Reporting**
- Search existing issues
- Test with the latest version
- Reproduce the issue consistently

### **Bug Report Template**
```markdown
**Bug Description**
A clear description of the bug.

**Steps to Reproduce**
1. Run command: `stormfinder -d example.com`
2. Observe behavior
3. Expected vs actual result

**Environment**
- OS: [e.g., macOS 12.0]
- Go Version: [e.g., 1.24.0]
- Stormfinder Version: [e.g., v2.9.0]

**Additional Context**
Any other relevant information.
```

## ✨ **Feature Requests**

### **Feature Request Template**
```markdown
**Feature Description**
Clear description of the proposed feature.

**Use Case**
Why is this feature needed? What problem does it solve?

**Proposed Solution**
How should this feature work?

**Alternatives Considered**
Other approaches you've considered.

**Additional Context**
Any other relevant information.
```

## 🧪 **Testing Guidelines**

### **Test Types**
- **Unit Tests**: Test individual functions
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete workflows
- **Performance Tests**: Benchmark critical paths

### **Writing Tests**
```go
func TestSubdomainEnumeration(t *testing.T) {
    // Arrange
    runner := NewRunner(defaultOptions)
    
    // Act
    results, err := runner.EnumerateDomain("example.com")
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, results)
}
```

### **Running Tests**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/runner

# Run with race detection
go test -race ./...
```

## 📚 **Documentation Guidelines**

### **Documentation Types**
- **Code Comments**: Inline documentation
- **README Files**: Package-level documentation
- **User Guides**: Feature usage examples
- **API Documentation**: Function/method documentation

### **Writing Style**
- Use clear, concise language
- Include practical examples
- Keep content up-to-date
- Use consistent formatting

## 🏗️ **Architecture Guidelines**

### **Package Structure**
- Keep packages focused and cohesive
- Use clear interfaces between packages
- Minimize dependencies between packages
- Follow Go package naming conventions

### **Adding New Features**

#### **New Enumeration Source**
1. Create source file in `pkg/subscraping/sources/`
2. Implement the `Source` interface
3. Add to source registry
4. Write comprehensive tests
5. Update documentation

#### **New AI Model**
1. Add model in `pkg/ai/models/`
2. Implement training and prediction methods
3. Add configuration options
4. Write performance tests
5. Document usage examples

#### **New Output Format**
1. Extend output writer in `pkg/runner/`
2. Add format-specific logic
3. Update CLI flags
4. Add examples and tests

## 🔒 **Security Guidelines**

### **Security Considerations**
- Never commit API keys or secrets
- Validate all user inputs
- Use secure HTTP clients
- Follow responsible disclosure for vulnerabilities

### **Reporting Security Issues**
- Email: security@stormfinder.dev
- Use GPG encryption if possible
- Provide detailed reproduction steps
- Allow reasonable time for fixes

## 🌟 **Recognition**

### **Contributors**
All contributors are recognized in:
- `THANKS.md` file
- GitHub contributors page
- Release notes for significant contributions

### **Maintainer Guidelines**
- Review PRs promptly
- Provide constructive feedback
- Help new contributors
- Maintain code quality standards

## 📞 **Getting Help**

### **Communication Channels**
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Discord**: Real-time community chat (coming soon)

### **Maintainer Contact**
- **GitHub**: @darshakkanani
- **Email**: darshak@stormfinder.dev

## 📄 **License**

By contributing to Stormfinder, you agree that your contributions will be licensed under the same license as the project (MIT License).

## 🙏 **Thank You**

Thank you for contributing to Stormfinder! Your efforts help make subdomain enumeration better for the entire security community.

**Happy Contributing! 🌪️**
