
# BankFlow

BankFlow is a mobile application built with Flutter designed to consume and validate the capabilities of a custom Bank API.

This project is part of a broader study focused on backend architecture, financial domain modeling, and secure transaction flows. The mobile app acts as a client layer to simulate real-world usage scenarios, enabling end-to-end validation of the system.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Purpose

The primary goals of this project are:

- Validate a modular Bank API built in Go
- Simulate real user interactions with financial operations
- Explore secure transaction flows and authorization strategies
- Evaluate architectural decisions across client and backend layers

## Scope

The application focuses on:

- User authentication (JWT-based)
- Account creation and management
- Financial operations (deposit, withdraw, transfer)
- Transaction history (statement with pagination)
- Integration with a strongly consistent backend

## Architecture Context

The backend follows a layered modular architecture:

- Domain-driven design principles
- Application layer for use cases and transaction control
- Infrastructure for persistence (PostgreSQL)
- Delivery via REST API

See architecture details: :contentReference[oaicite:0]{index=0}

API contract: :contentReference[oaicite:1]{index=1}

## Motivation

This project was created to:

- Revisit Go development after several years
- Explore backend design decisions not fully addressed in previous professional experiences
- Build a complete system (backend + mobile client) with controlled architectural evolution

## Notes

- This is a study project and not intended for production use
- Focus is on correctness, structure, and learning—not feature completeness