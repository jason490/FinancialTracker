# FinancialTracker Frontend Design Context

The goal is to have a functional UI that is efficient, matching colors themes, and no werid UI interactions.

## 1. Dynamic Dashboard

A fully customizable dashboard that the user can setup or use the defaults.

*   **Grid:** Each grid box is dynamic in size and should be large for easy snaping.
*   **Components:** The user can select/remove components on and off from the dashboard. If the item of the component is a list of items, it should be scrollable.
*   **Sizing:** Each component will have minimum and maximum size it can take.
*   **Snaping:** Bringing in different components and allowing for snaping to a grid box.
*   **Interactions:** Mobile and Website/webapp have different interactions, i.e. mobile should click and hold a component to activate edit mode, while the website/webapp just clicks customize.

## 2. Themes

Allowing the user to adjust the themes with pre-built color schemes. The themes should be defined globally.

*   **Default (Light):** Use the default looks of a normal financial tracker i.e. Rocket Money.
*   **Dark:** Change the color scheme to dark.
*   **Tokyo Night:** Change the color scheme to Tokyo Night.
*   **Coffee:** Change the color scheme to Coffee.

## 3. Dynamic Sizing
*   **Any Device:** All items should fit within the parent element
*   **Mobile:** Ensure the UI is mobile friendly with not too much clutter and use Capacitor when needed
*   **Desktop:** Ensure the UI is desktop friendly

