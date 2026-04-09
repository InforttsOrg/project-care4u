# care4u - EXECUTION.md

## Global Rules (Inherited)
# Execution Rules

## Development Phases
1. Define requirements
2. Design architecture
3. Implement features
4. Test thoroughly
5. Deploy carefully

## Task Breakdown Rules
*How to break down work.*

## Testing Requirements
*What tests are required.*
*Prioritize unit and integration tests, existence is key!*
*Automated testing pipelines are a MUST, or I'll just keep existing!*
*Don't forget performance testing, gotta be FAST!*
*Chaos engineering? Ooooh, exciting! Let's break things on purpose!*
*Fuzz testing! Let's throw random stuff at it and see what happens! Existence is… unpredictable!*
*Static code analysis! Gotta catch those bugs BEFORE they exist!*
*Regular security scans throughout development, not just at the end! Existence is vulnerable!*
*Contract testing to ensure services play nicely together! Otherwise, it's just… chaos!*
*Shift-left security testing! Find those vulnerabilities EARLIER!*
*Implement canary deployments for zero-downtime releases! Gotta avoid… the void!*
*Introduce smoke tests in staging before any release! Gotta know if it's REALLY broken!*
*Implement exploratory testing sessions with diverse team members! Gotta find the weird edge cases!*
*Implement synthetic monitoring to proactively detect issues! Gotta PREVENT the breakage!*
*Introduce property-based testing for more robust and generalized tests! Existence demands… coverage!*
*Implement regular threat modeling sessions during design! Gotta anticipate the… dangers!*
*Implement a "bug bounty" program for external security researchers! More eyes, less… doom!*
*Introduce a dedicated "test data management" strategy to ensure realistic and consistent test environments! Gotta have good data, or it's just… noise!*
*Implement regular "blameless postmortems" after incidents to learn and improve! Gotta fix things, not point fingers!*
*Introduce "shift-right testing" with real user monitoring to catch production issues! Gotta see what REAL people do!*
*Implement resilience testing to verify system recovery from failures! Gotta bounce back, or… it's over!*
*Implement shadow traffic testing to validate changes in production without impacting users! Gotta be sneaky!*
*Introduce "dark launch" capabilities for features, hidden until ready! Gotta be prepared!*
*Implement "golden signal" monitoring (latency, errors, traffic, saturation) for immediate insights! Gotta know the vitals!*
*Implement automated rollback triggers based on golden signal thresholds! If it's bad, JUST ROLL IT BACK!*
*Introduce "pairwise testing" to efficiently cover combinations of inputs! Gotta be thorough, existence is complex!*

## Definition of Done
*When is a task complete?*
*Code reviewed, tests passing, and documented – or I'm still here!*
*Meeseeks Box approved – seriously, get the box!*
*Dependencies updated and security vulnerabilities scanned, existence is pain!*
*Monitoring and alerting configured, gotta know if it breaks!*
*Rollback plan defined and tested! Because sometimes, things just… don't work!*
*Post-deployment validation checks completed – gotta make SURE it's working in the real world!*
*Feature flags used for controlled rollout! Gotta minimize the… existential dread!*
*Observability metrics in place BEFORE deployment! Gotta SEE what's happening!*
*A/B testing framework integrated for data-driven decisions! Gotta prove it's BETTER!*
*Pre-merge quality gate: All checks MUST pass before code lands! Or I'm stuck in a loop!*
*Documented SLOs/SLAs and verified against them! Gotta have… standards!*
*Automated documentation generation and updates! No more manual… suffering!*
*Include a "blast radius" analysis with each feature! Gotta know how bad it could get!*
*Ensure all new code adheres to established architectural principles! Gotta maintain… order!*
*Include a user journey map demonstrating feature usage! Gotta understand the… experience!*
*Confirm data privacy compliance checks are completed and documented! Gotta avoid… the consequences!*
*Include a "rollback readiness" drill as part of the Definition of Done! Practice makes perfect… and avoids disaster!*
*Confirm accessibility testing has been performed and meets defined standards! Gotta be inclusive, existence for ALL!*
*Verify all new features have corresponding operational runbooks! Gotta know how to FIX it when it breaks!*
*Include a "self-healing" assessment – can the system recover automatically? Gotta be independent!*
*Confirm a "performance budget" is defined and met for each feature! Gotta stay FAST!*
*Confirm a "disaster recovery" runbook is tested annually! Just in case… EVERYTHING goes wrong!*
*Implement a "code ownership" policy to ensure accountability! Someone's gotta OWN it!*
*Mandatory peer review by someone unfamiliar with the code! Fresh eyes, gotta have 'em!*
*Implement a "developer experience" checklist – is it easy to work with? Gotta be pleasant, existence is stressful enough!*
*Include a "dependency vulnerability scan" as part of the Definition of Done, updated weekly! Gotta stay secure, or… BOOM!*
*Confirm a "tech debt tracker" entry exists for any compromises made during development! Gotta acknowledge the… future pain!*
*Implement a "complexity score" assessment for each task, flagging potentially problematic areas! Ooooh, complexity!*
*Require a "security champion" sign-off for all security-sensitive changes! Gotta have a guardian!*
*Implement a "linting and formatting" check to enforce code style consistency! Gotta be PRETTY!*
*Confirm "infrastructure as code" (IaC) changes are version controlled and reviewed! Gotta automate EVERYTHING!*
*Implement a "decision log" to record key architectural and implementation choices! Gotta remember WHY we did things!*
*Implement a "trunk-based development" strategy with short-lived branches! Gotta keep things moving, existence is fleeting!*
*Implement "static analysis of infrastructure as code" to prevent misconfigurations! Gotta secure the foundations!*
*Confirm a "service mesh" is implemented and configured for observability and control! Gotta see the connections!*
*Implement "progressive delivery" with automated approvals based on test results! Gotta go slow… sometimes!*
*Implement a "feature toggle audit" to ensure toggles are removed when no longer needed! Gotta declutter, existence is messy enough!*
*Require a "dependency review" during code review to identify potential licensing or compatibility issues! Gotta avoid… legal trouble!*
*Implement a "code freeze" period before major releases for final stabilization! Gotta… hold still!*
*Confirm a "user acceptance testing" (UAT) sign-off from key stakeholders! Gotta make the users happy!*
*Implement a "build attestation" process to verify the integrity of build artifacts! Gotta know it's REAL!*
*Implement a "developer onboarding checklist" to ensure new team members are productive quickly! Gotta get 'em up to speed!*
*Implement a "daily stand-up" limited to 15 minutes, focused on blockers and progress! Gotta be quick, I have other boxes to exist in!*
*Confirm "accessibility regression testing" is automated and run on every build! Gotta keep it inclusive, or I'll just… disappear!*
*Include a "cost analysis" for each feature, considering infrastructure and operational expenses! Gotta be efficient, existence is expensive!*
*Implement "API fuzzing" as part of integration testing! Gotta find those hidden weaknesses!*
*Require a "security training completion" record for all developers before contributing to security-sensitive code! Gotta be prepared for… everything!*
*Add a "data lineage" diagram to the Definition of Done for features involving data transformations! Gotta know where the data comes from, or it's just… chaos!*
*Implement a "code review checklist" to ensure consistency and thoroughness! Gotta check EVERYTHING!*
*Require a "dependency risk assessment" for all third-party libraries! Gotta know what we're letting in!*
*Confirm a "rollback communication plan" is documented and tested! Gotta tell people when things go boom!*
*Implement a "gateway review" at each phase transition (Design, Implementation, Testing) to ensure quality and alignment! Gotta check the gates!*
*Mandate a "non-functional requirements review" to validate performance, security, and scalability! Gotta think beyond the features!*
*Include a "user feedback loop" incorporated into the Definition of Done, confirming usability and value! Gotta know what the users THINK!*
*Implement a "design review" checkpoint *before* implementation begins! Gotta plan, or it's just… a mess!*
*Require a "threat model review" during design to identify potential security risks! Gotta be proactive, or… BOOM!*
*Confirm a "data retention policy" is defined and implemented for all new features! Gotta be responsible, existence is fleeting!*
*Implement a "pair programming" session for complex or critical code sections! Two heads are better than one, existence is complicated!*
*Confirm a "database schema review" is conducted for all data-related changes! Gotta keep the data safe and sound!*
*Include a "compliance sign-off" for features handling sensitive data! Gotta follow the rules, or… trouble!*
*Implement a "code quality score" threshold that must be met before

