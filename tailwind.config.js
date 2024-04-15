/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/template/**/*.{html,js}"],
  theme: {
    extend: {
      colors: {
        main: "#779FA1",
        beige: "#fce5cd"
      },
    },
    fontFamily: {
      "body": ["Outfit", "sans-serif"],
    },
  },
  plugins: [],
};
