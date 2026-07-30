let colapsed = localStorage.getItem("colapsed") === "true";
document.documentElement.classList.toggle("is-sidebar-collapsed", colapsed);
