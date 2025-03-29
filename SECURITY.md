Below is the updated security policy with a dependency version table that uses icons to indicate that versions from the specified minimum onward are supported:

---

# Security Policy

## Supported Dependency Versions

This section specifies the minimum required versions for the project's dependencies. Any version equal to or higher than these is considered supported with security updates.

| Dependency                                | Minimum Version | Supported          |
| ----------------------------------------- | --------------- | ------------------ |
| github.com/eclipse/paho.mqtt.golang       | v1.5.0          | :white_check_mark: |
| github.com/gin-gonic/gin                  | v1.10.0         | :white_check_mark: |
| github.com/gorilla/websocket              | v1.5.3          | :white_check_mark: |
| go.mongodb.org/mongo-driver               | v1.17.3         | :white_check_mark: |

## Reporting a Vulnerability

To report a vulnerability in our project, please follow the guidelines below:

- **Reporting Location:** Submit your report to [designated email or platform].
- **Required Information:**  
  - **Detailed Description:** Provide a comprehensive explanation of the issue, including the steps to reproduce the vulnerability and its potential impact.
  - **Mandatory Evidence:**  
    - *Postman Tests:* Attach the Postman collection or HTTP requests that demonstrate the issue.
    - *Unit Tests:* Include test cases that show the unexpected behavior or security flaw.
    - *Integration Tests:* Provide complete scenarios that illustrate how the vulnerability affects the interaction between system components.
- **Follow-Up Process:** After submitting your report, you will receive a notification regarding its review status and be kept informed on the progress of the resolution.
- **Confidentiality:** We ensure the confidentiality of your report during the verification and resolution process.

Thank you for your cooperation and commitment to maintaining the project's security.

---
