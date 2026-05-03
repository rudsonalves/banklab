# GitHub Project Organization

The project board is structured to support clear decision-making, collaboration, and execution. Each issue is classified using three key dimensions: **Type**, **Area**, and **Priority**.

---

## **Type**

Defines the nature of the work:

* **Feature**
  Introduces new functionality or capabilities to the system.

* **Improvement**
  Enhances or refines existing behavior without changing core functionality.

* **Bug**
  Fixes incorrect or unexpected system behavior.

* **Research**
  Used for exploration, design decisions, or evaluation of approaches before implementation.

---

## **Area**

Represents the part of the system affected by the work, aligned with the system architecture:

* **Auth** — authentication, JWT, sessions
* **Account** — account lifecycle and balance operations
* **Ledger** — transactions, consistency, financial core
* **Customer** — customer data and identity model
* **Mobile** — Flutter client application
* **Admin Web** — internal operational interface
* **Client Web** — customer-facing web interface
* **Security** — TOTP, device trust, Zero Trust evolution
* **Infra** — database, configuration, and infrastructure

Area represents the primary ownership of the issue. Other impacted areas should be described in the issue body.

---

## **Priority**

Indicates the importance and urgency of the work:

* **High**
  Critical for system evolution or blocks other work

* **Medium**
  Important but not blocking

* **Low**
  Optional improvements or future enhancements

---

## **Usage Guidelines**

* Every issue must define **Type**, **Area**, and **Priority**
* **Research** issues are used to align decisions before implementation
* Implementation work should only begin after design is clearly defined
* The board reflects both **technical structure** and **product direction**
