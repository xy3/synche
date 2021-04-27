module.exports = {
  purge: ["./pages/**/*.{js,ts,jsx,tsx}", "./components/**/*.{js,ts,jsx,tsx}"],
  darkMode: false,
  theme: {
    extend: {
      colors: {
        "brand-dark-blue": "#011e38",
        "brand-blue": "#027cac",
        "brand-cyan": "#04b1cd",
      },
      fontFamily: {
        rubik: ["Rubik", "sans-serif"],
        "noto-sans-kr": ["Noto Sans KR", "sans-serif"],
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [require("@tailwindcss/custom-forms")],
};
