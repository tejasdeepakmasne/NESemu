# Contributing to NESemu

Thank you for considering contributing to this project! Here are some guidelines to help you get started.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Setting Up Your Development Environment](#setting-up-your-development-environment)
3. [Making Changes](#making-changes)
4. [Running Tests](#running-tests)
5. [Submitting Changes](#submitting-changes)
6. [Code of Conduct](#code-of-conduct)

## Getting Started

1. **Fork the repository**: Click the "Fork" button at the top right corner of this page to create a copy of the repository under your GitHub account.

2. **Clone your fork**: Clone your forked repository to your local machine using the command:
   ```bash
   git clone https://github.com/yourusername/NESemu.git
   cd NESemu

3. **Add the upstream repository**: Add the original repository as a remote to keep your fork up to date.
   ```bash
   git remote add upstream https://github.com/originaluser/NESemu.git
   ```

## Setting Up Your Development Environment

1. **Install Go**: Ensure you have Go installed on your system. Follow the instructions [here](https://golang.org/doc/install) if you need to install it.

2. **Install dependencies**: Navigate to your project directory and run:
   ```bash
   go mod tidy
   ```

3. **Create a feature branch**: Create a new branch for your work.
   ```bash
   git checkout -b your-feature-branch
   ```

## Making Changes

1. **Write your code**: Make your changes in the relevant files. Ensure your code follows the project's coding standards.

2. **Commit your changes**: Commit your changes with a clear and descriptive commit message.
   ```bash
   git add .
   git commit -m "Description of the changes"
   ```

3. **Keep your branch up to date**: Before submitting your changes, ensure your branch is up to date with the latest changes from the upstream repository.
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

## Running Tests

1. **Run the tests**: Ensure all tests pass before submitting your changes.
   ```bash
   go test ./...
   ```

2. **Add new tests**: If you added new features or fixed bugs, add corresponding tests to cover your changes.

## Submitting Changes

1. **Push your branch**: Push your feature branch to your forked repository.
   ```bash
   git push origin your-feature-branch
   ```

2. **Open a Pull Request**: Navigate to the original repository and open a Pull Request (PR) from your forked repository's branch. Provide a clear description of your changes and the problem they solve.

3. **Address review comments**: Be responsive to feedback and make necessary changes as requested by reviewers.

## Code of Conduct

Please adhere to the project's [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming and inclusive environment for everyone.

Thank you for contributing!

