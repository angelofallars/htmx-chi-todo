/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/**/*.{html,js,ts,go,templ}"],
  theme: {
    extend: {},
  },
  plugins: [],
  variants: {
    extend: {
        display: ["group-hover"],
    },
    width: {
        "128": "32rem",
        "192": "48rem",
    },
  },
}
