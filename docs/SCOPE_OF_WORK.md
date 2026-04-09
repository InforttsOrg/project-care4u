# Scope of Work & Functional Requirements (BRD/FRD)
**Project:** Care4uuu - On-Demand Healthcare Platform (MVP)
**Provider:** Growth Hikes
**Lead Developer:** Sahil Rathee
**Backend Developer:** Jitender Bachhraj
**Document Version:** 1.1

---

## 1. Executive Summary
This document delineates the comprehensive scope of work for the development of "Care4uuu," a dual-application platform designed to connect healthcare professionals (Physiotherapists, Nurses, etc.) with patients for on-demand and scheduled home services. This Scope of Work (SOW) serves as the primary reference for the "MVP" (Minimum Viable Product) deliverable executed by **Growth Hikes**.

**Objective:** To deliver a functional, scalable MVP that facilitates trust-based, efficient medical service bookings with secure payments and real-time communication.

---

## 2. Platform Architecture
The solution shall consist of three distinct interfaces:
1.  **Patient App** (iOS/Android): For customers to discover and book services.
2.  **Professional App** (iOS/Android): For medical staff to manage availability and requests.
3.  **Admin Dashboard** (Web): For Care4uuu internal teams to manage users, disputes, and settings.

---

## 3. Functional Requirements (In-Scope)

### A. Patient App (Demand Side)
*   **Authentication:** Secure login via Google, Apple, and Mobile OTP.
*   **Discovery:**
    *   Map-based view of nearby professionals.
    *   Filtering capabilities by Service Type (Physiotherapy, Nursing, etc.).
*   **Booking System:**
    *   **On-Demand:** "Request Now" feature for immediate emergency assistance.
    *   **Scheduled:** Capabilities to book slots based on the professional's real-time calendar.
*   **Virtual Wallet:**
    *   Top-up functionality (UPI, Cards, Netbanking).
    *   Transaction history view.
    *   Auto-deduction mechanisms for bookings.
*   **Communication:**
    *   In-app Chat with the professional (active only during active booking windows).
    *   AI-Powered FAQ Chatbot for general support queries.
*   **Trust & Safety:** SOS button integration for emergencies.

### B. Professional App (Supply Side)
*   **Profile Management:** Functionality to upload credentials, set service rates, and manage bio.
*   **Availability Manager:** Options to toggle "Online/Offline" status and manage calendar slots.
*   **Request Management:** Interface to Accept/Reject booking requests with countdown timers.
*   **Navigation:** Integrated map navigation to the patient's location.
*   **Earnings:** Dashboard to view daily/weekly booking revenue (pre-platform commission deduction).

### C. Admin Dashboard (Internal)
*   **User Management:** Tools to verify and approve Professional profiles (KYC check).
*   **Financials:** View Escrow holdings, process payouts, and manage platform commission rates.
*   **Dispute Resolution:** Access to chat logs (for compliance audit) and dispute ticket management.

---

## 4. Technical Specifications
*   **Mobile Apps:** Flutter (Single codebase compiled for iOS & Android).
*   **Backend:** Golang (High-concurrency microservices architecture).
*   **Database:** PostgreSQL (Transactional data), Redis (Caching/Geo-location).
*   **Infrastructure:** Cloud-native deployment (AWS/GCP) with Docker containerization.
*   **Payments:** Razorpay Escrow (Split payment logic between Platform & Professional).

---

## 5. Exclusions (Out of Scope)
The following features are **explicitly excluded** from this MVP agreement. These may be developed as separate "Change Orders" at an additional cost:
*   [ ] Full-scale Pharmacy Marketplace & Fulfillment.
*   [ ] Integration with IoT Medical Devices.
*   [ ] Voice/Audio Interaction (beyond standard phone calls).
*   [ ] ERP/CRM integrations for large hospital networks.
*   [ ] Offline-first mode (App requires active internet connection).
*   [ ] Manual Data Migration from legacy systems.

---

## 6. Development Timeline
The estimated timeline for MVP delivery is **12 Weeks (approx. 60 Working Days)** from the Project Start Date (Day 0).

| Phase | Description | Estimated Duration |
|-------|-------------|--------------------|
| **Phase 1** | System Design, UI/UX Finalization, Infra Setup | Weeks 1-2 |
| **Phase 2** | Core Auth, User Profiles, Admin Dashboard | Weeks 3-5 |
| **Phase 3** | Booking Engine, Calendar, Wallet Logic | Weeks 6-8 |
| **Phase 4** | Chat, AI Bot, Notifications, Geolocation | Weeks 9-10 |
| **Phase 5** | QA Testing, Security Audit, UAT, Deployment | Weeks 11-12 |

*> **Note:** "Working Days" excludes weekends and public holidays. The timeline relies on timely client feedback (within 24 hours).*

---

## 7. Deliverables
1.  Source Code Repository (GitHub/GitLab).
2.  Compiled APK/IPA files for testing.
3.  Admin Dashboard URL & Credentials.
4.  Technical Documentation (API Swagger Docs, Architecture Diagram).
5.  30 Days of Post-Launch Bug Fix Support (Hypercare).

---

## 8. Compliance & Legal Disclaimer
*   **Data Privacy:** The platform shall be designed to ensure compliance with local data protection laws (e.g., DPDP Act/GDPR) via encryption at rest and in transit.
*   **Medical Liability:** Growth Hikes acts solely as a technology partner, not a healthcare provider. Growth Hikes shall not be held liable for the medical advice or services rendered by professionals on the platform.
