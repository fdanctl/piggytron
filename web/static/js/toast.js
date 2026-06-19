// const (
// 	Success toastLevel = "success"
// 	Warning toastLevel = "warning"
// 	Error   toastLevel = "error"
// 	Info    toastLevel = "info"
// )
function makeSVG(level, size = 24, color = "currentColor", className = "") {
  let children = "";
  switch (level) {
    case "success":
      children = `<path d="M20 6 9 17l-5-5"></path>`;

    case "warning":
      children = `<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"></path><path d="M12 9v4"></path><path d="M12 17h.01"></path>`;

    case "error":
      children = `<circle cx="12" cy="12" r="10"></circle><path d="m15 9-6 6"></path><path d="m9 9 6 6"></path>`;

    case "info":
      children = `<circle cx="12" cy="12" r="10"></circle><path d="M12 16v-4"></path><path d="M12 8h.01"></path>`;

    default:
      children = `<circle cx="12" cy="12" r="10"></circle><path d="M12 16v-4"></path><path d="M12 8h.01"></path>`;
  }

  return `
	<svg
		xmlns="http://www.w3.org/2000/svg"
		width=${size}
		height=${size}
		viewBox="0 0 24 24"
		fill="none"
		stroke=${color}
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
		class=${className}
	>
		${children}
	</svg>
  `;
}

export function showToast(level, msg) {
  const toast = document.createElement("div");
  toast.classList.add("toast", level);
  toast.dataset.animationend = "ui.element.remove";

  const inner = document.createElement("div");
  const p = document.createElement("p");
  p.innerText = msg;
  inner.appendChild(p);

  toast.innerHTML = makeSVG(level, 16);
  toast.append(inner);

  document.getElementById("toast-root").appendChild(toast);
}
