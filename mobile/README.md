# BankFlow

BankFlow is a Flutter mobile application designed to validate and exercise a custom-built banking API.

Rather than a feature-driven product, this application acts as a controlled client environment to simulate real-world financial workflows and validate end-to-end system behavior.

## Purpose

This project is part of a broader effort to explore system design across mobile and backend layers, with a focus on:

- End-to-end validation of financial operations
- Consistency between client behavior and backend guarantees
- API contract design and integration boundaries
- Execution context and transaction safety

## Scope

The application includes:

- JWT-based authentication
- Account creation and lifecycle management
- Financial operations (deposit, withdraw, transfer)
- Transaction history with cursor-based pagination
- Integration with a strongly consistent backend

## Architectural Role

This application is not treated as an isolated frontend.

It is designed to:

- Validate backend assumptions through real usage flows
- Expose inconsistencies in API design and data contracts
- Ensure alignment between user interaction and backend behavior

The mobile layer is intentionally structured to reflect production concerns such as:

- Clear separation between UI, state, and business logic
- Predictable data flow across layers
- Explicit handling of asynchronous operations and failures

## Backend Context

The backend (Go) follows a layered architecture:

- Domain-driven design principles
- Application layer for use case orchestration
- Infrastructure layer (PostgreSQL)
- Delivery via REST APIs

## Motivation

This project was created to:

- Revisit Go development in a structured way
- Explore backend design decisions beyond prior production constraints
- Build and evolve a complete system (backend + mobile) with full control over architectural decisions

## Notes

- This is an engineering-focused project, not a production product
- Emphasis is placed on correctness, consistency, and system design
- Features are intentionally limited to maintain focus on architectural validation

## License

MIT License — see the LICENSE file for details