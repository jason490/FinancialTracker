## 1. Project Overview
A high-performance, user-centric financial tracking application inspired by Rocket Money. The core objective is to provide blazing-fast functionality and seamless expense management by directly syncing with bank accounts via Plaid. The system allows users to categorize, filter, and visualize their financial data through an easily customizable dashboard.

---

## 2. Technical Stack
The stack is chosen for maximum performance, mobile-first responsiveness, and a modern SPA (Single Page Application) experience, while maintaining the legacy system during transition.

* **Backend:** Golang (using `echo` framework for HTTP REST API and legacy HTML rendering).
* **Database:** SQLite (Embedded, high-performance relational database).
* **New Frontend:** SolidJS (Fine-grained reactive UI library for high efficiency).
* **Mobile Bridge:** Capacitor (Cross-platform bridge for Android and Web. **iOS excluded**).
* **Styling (New):** Vanilla CSS / CSS Modules (Targeting mobile-first design; **No TailwindCSS** for new development).
* **Legacy Stack:** Templ, HTMX, and TailwindCSS (To be phased out).
* **External APIs:** Plaid (for fetching account balances and transaction histories).

---

## 3. Core Functionalities (Target State)

### 3.1 Authentication & Account Management
* **SSO Logins:** Integration with standard OAuth providers (Google, Apple, etc.).
* **API-Based Auth:** Secure JWT or session management for the SolidJS app.
* **Account Controls:** Complete user lifecycle management (Creation, Deletion, Data Export).

### 3.2 Plaid Integration & Expense Syncing
* **Automated Fetching:** Scheduled or background synchronization via Plaid.
* **Native Plaid Link:** Integration of the Plaid Link SDK via Capacitor for a native Android experience.

### 3.3 Advanced Tagging & Categorization Engine
* **Auto-Tagging:** Rule-based tagging engine (shared by legacy and new frontends).
* **Smart Re-Tagging:** System-wide tag updates based on user corrections.

### 3.4 ViewPoint (Dynamic Dashboard)
* **High-Performance UI:** Leveraging SolidJS for instantaneous updates and smooth transitions on mobile.
* **Mobile-First Analytics:** Touch-optimized charts and drill-down views.

---

## 4. Transition & Folder Structure

The project maintains the existing structure while introducing the new frontend environment.

```text
.
├── app/                   # Legacy Full-stack (Go + Templ/HTMX)
│   ├── internal/          
│   │   ├── server/        # Echo routes (serving both HTML and JSON)
│   │   ├── storage/       
│   │   ├── services/      # Core business logic (shared)
│   │   └── models/        
│   └── web/               # Legacy Templ/HTMX assets
├── frontend/              # New SolidJS application
│   ├── src/               # SolidJS components, stores, and pages
│   ├── public/            # Static assets
│   ├── package.json       
│   └── vite.config.ts     
├── capacitor.config.ts    # Capacitor configuration
├── android/               # Native Android project
├── Makefile               # Commands for both legacy and new stacks
└── GEMINI.md              # Project instructions
```
