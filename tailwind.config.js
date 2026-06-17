/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/web/templ/**/*.templ",
    "./app/web/templ/**/*.go",
    "./app/internal/utils/**/*.go",
  ],
  darkMode: 'class',
  theme: {
    extend: {},
  },
  plugins: [],
}
