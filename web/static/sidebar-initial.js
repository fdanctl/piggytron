let colapsed = localStorage.getItem("colapsed") === "true";
const sidebar = document.getElementById("sidebar");
sidebar.classList.toggle("colapsed", colapsed);
