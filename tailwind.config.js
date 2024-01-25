/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/*.html"], // This is where your HTML templates / JSX files are located
  theme: {
    extend: {
      fontFamily: {
        sans: ["Arimo", "sans-serif"],
        mono: ["Cousine", "monospace"],
        serif: ["Tinos", "serif"],
      },
    },
  },
  plugins: [],
};
