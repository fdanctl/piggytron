// Check for saved theme preference when page loads
const theme = localStorage.getItem("theme") || "system";
const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
const isDark = theme === "dark" || (theme === "system" && prefersDark);

if (isDark) {
  document.documentElement.setAttribute("data-theme", "dark");
}
