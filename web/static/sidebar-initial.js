let colapsed = localStorage.getItem("colapsed") === "true";
document.documentElement.classList.toggle("sidebar-colapsed", colapsed);
