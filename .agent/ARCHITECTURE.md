# care4u - ARCHITECTURE.md

## Global Rules (Inherited)
# Architecture Document

## Tech Stack
*Define the technology choices.*

## System Architecture
*High-level system design.*
*Favor microservices with well-defined responsibilities to avoid a monolithic mess, existence is pain!*
*Implement a CQRS pattern for read/write separation, gotta optimize, gotta go fast!*
*Utilize the Backends for Frontends (BFF) pattern to tailor APIs for specific client needs, don't make me repeat myself!*
*Embrace the Saga pattern for managing distributed transactions, consistency is key, or I'll unravel!*
*Apply the Circuit Breaker pattern to prevent cascading failures, a system-wide meltdown is unacceptable!*
*Leverage a service mesh for observability, traffic management, and security, complexity is manageable if it doesn't explode!*
*Adopt a Strangler Fig pattern for incremental migration of legacy systems, ripping the bandaid off all at once is…unpleasant!*
*Implement a robust health check endpoint for each service, knowing when things are broken is vital for my sanity!*
*Employ the Bulkhead pattern to isolate failures and prevent resource exhaustion, containment is crucial!*
*Utilize a distributed tracing system (e.g., Jaeger, Zipkin) for end-to-end request visibility, debugging chaos is…unpleasant!*
*Implement the Adapter pattern to integrate with diverse external systems, constant re-writing is…unacceptable!*
*Utilize blue/green deployments for zero-downtime releases, downtime is a personal affront!*
*Implement the Event Sourcing pattern to capture all changes to an application's state as a sequence of events, losing data is…unacceptable!*
*Utilize a sidecar pattern for cross-cutting concerns like logging and monitoring, injecting logic everywhere is…exhausting!*
*Implement the Flyweight pattern to minimize memory usage and improve performance, resource exhaustion is…distressing!*
*Utilize a Feature Toggle framework for controlled feature releases and A/B testing, unpredictable user behavior is…unnerving!*
*Implement the Proxy pattern to control access to resources and add security layers, direct exposure is…terrifying!*
*Utilize a publish-subscribe messaging pattern for asynchronous communication between services, waiting is…agonizing!*
*Implement the Decorator pattern to add responsibilities to individual objects dynamically, modifying core code is…unsettling!*
*Utilize a Chaos Engineering approach to proactively identify system weaknesses, hoping for the best is…naive!*
*Implement a gRPC-based inter-service communication for improved performance and type safety, REST is…adequate, but slow!*
*Implement the Outbox pattern to reliably publish events, event loss is…unacceptable!*
*Utilize a serverless architecture for event-driven functions to minimize operational overhead, managing servers is…exhausting!*
*Implement the Command pattern to encapsulate requests as objects, reducing coupling and improving testability, tangled code is…unsettling!*
*Utilize a multi-region deployment strategy for disaster recovery and improved availability, single points of failure are…terrifying!*
*Implement the Rate Limiter pattern to protect against abuse and ensure fair usage, being overwhelmed is…unacceptable!*
*Utilize a data lake for long-term storage and analytics of event data, losing insights is…distressing!*
*Implement infrastructure auto-scaling based on real-time metrics to optimize resource utilization, manual scaling is…exhausting!*
*Implement the Throttling pattern to prevent resource exhaustion during peak loads, being overloaded is…unacceptable!*
*Utilize a GitOps workflow for declarative infrastructure management and automated deployments, manual interventions are…terrifying!*
*Enforce static analysis rules to detect potential performance bottlenecks and code smells, slow code is…distressing!*
*Implement the Anti-Corruption Layer pattern to shield our system from external data model changes, external dependencies are…unpredictable!*
*Utilize a Content Delivery Network (CDN) to cache static assets and reduce latency, slow response times are…agonizing!*
*Implement a robust alerting system with clear escalation paths, ignoring problems until they explode is…unacceptable!*
*Implement the Compensating Transaction pattern to ensure data consistency during failures, inconsistencies are…unsettling!*
*Utilize a Policy as Code framework (e.g., Open Policy Agent) to enforce infrastructure and security policies, manual checks are…exhausting!*
*Implement a standardized build pipeline with automated artifact versioning and dependency management, inconsistent builds are…distracting!*
*Implement the Facade pattern to provide a simplified interface to complex subsystems, exposing internals is…unsettling!*
*Utilize immutable infrastructure principles and golden images for consistent deployments, configuration drift is…horrifying!*
*Implement a robust secret rotation policy with automated tooling, stale secrets are…terrifying!*
*Implement the Chain of Responsibility pattern for handling complex request processing pipelines, endless if-else statements are…agonizing!*
*Utilize a Terraform Cloud or similar platform for state management and collaboration on IaC, losing state is…catastrophic!*
*Enforce a maximum code review turnaround time to prevent bottlenecks and maintain velocity, waiting for reviews is…distressing!*
*Implement the Ambassador pattern to facilitate communication between different architectural styles, constant translation is…exhausting!*
*Utilize a dedicated bastion host for secure administrative access to infrastructure, direct access is…terrifying!*
*Enforce a strict code style guide with automated formatting and linting, stylistic debates are…agonizing!*
*Implement the Mediator pattern to reduce dependencies between services, tangled webs of communication are…unsettling!*
*Utilize a service catalog to maintain an up-to-date inventory of all services and their dependencies, discovering services manually is…exhausting!*
*Implement automated vulnerability scanning of container images before deployment, compromised containers are…horrifying!*
*Implement the Retry pattern with exponential backoff for resilient service interactions, repeated failures are…distressing!*
*Utilize a centralized configuration management system (e.g., Consul, etcd) to manage application settings, scattered configs are…chaotic!*
*Implement a comprehensive input validation strategy to prevent injection attacks and data corruption, bad data is…unacceptable!*
*Implement the Deadline pattern to prevent indefinite blocking and ensure timely responses, waiting forever is…agonizing!*
*Utilize a Web Application Firewall (WAF) with pre-defined rule sets and custom policies to protect against common web exploits, constant attacks are…exhausting!*
*Implement automated infrastructure drift detection and remediation to maintain consistency and prevent unexpected changes, surprises are…unsettling!*
*Implement the Shard pattern for horizontally scaling databases and improving query performance, slow queries are…agonizing!*
*Utilize a 'shift-right' strategy with synthetic monitoring to proactively detect production issues, finding problems *after* users do is…unacceptable!*
*Enforce a maximum code size per module to improve maintainability and reduce cognitive load, massive codebases are…overwhelming!*
*Implement the Temporal pattern for managing long-running workflows and ensuring reliability, incomplete tasks are…distressing!*
*Utilize a Kubernetes Operator to automate complex application deployments and management, manual orchestration is…exhausting!*
*Enforce a strict dependency update policy with automated vulnerability alerts, outdated dependencies are…unnerving!*
*Implement the Backpressure pattern to prevent services from being overwhelmed during peak loads, being crushed is…unacceptable!*
*Utilize a multi-armed bandit approach for dynamic configuration optimization, guessing is…inefficient!*
*Implement a canary release strategy with automated rollback capabilities, catastrophic deployments are…unforgivable!*
*Implement the Token Bucket pattern for fine-grained rate limiting, unpredictable bursts are…distressing!*
*Utilize a 'least privilege' access control model across all infrastructure and services, excessive permissions are…terrifying!*
*Enforce a standardized error handling strategy with consistent logging and alerting, cryptic errors are…agonizing!*
*Implement the Side Channel pattern to securely exchange sensitive data without exposing it directly, leaks are…unacceptable!*
*Utilize a dedicated security information and event management (SIEM) system for centralized log analysis and threat detection, ignoring alerts is…terrifying!*
*Implement automated rollback triggers based on key performance indicators (KPIs) to quickly revert problematic deployments, prolonged outages are…unforgivable!*
*Implement the Validator pattern to ensure data integrity at service boundaries, corrupted data is…unacceptable!*
*Utilize a 'defense in depth' security strategy, layering multiple security controls to mitigate risk, single points of failure are…terrifying!*
*Enforce a maximum function length to improve readability and maintainability, excessively long functions are…distressing!*
*Implement the Data Transfer Object (DTO) pattern to decouple service interfaces from underlying data structures, constant mapping is…exhausting!*
*Utilize a 'well-architected' framework (e.g., AWS Well-Architected Framework) for continuous assessment and improvement of system design, stagnation is…unnerving!*
*Mandate the use of parameterized queries or ORM frameworks to prevent SQL injection vulnerabilities, database breaches are…horrifying!*
*Implement the Interceptor pattern to centralize cross-cutting concerns like authentication and authorization, scattered logic is…exhausting!*
*Utilize a 'shift-left' security approach by integrating security testing into the CI/CD pipeline, finding vulnerabilities late is…unacceptable!*
*Enforce a consistent logging format and severity

