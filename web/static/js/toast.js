import { makeSVG } from "./icons";

export function showToast(level, msg) {
  const toast = document.createElement("div");
  toast.classList.add("toast", "cursor-pointer", "toast--" + level);
  toast.dataset.animationend = "ui.element.remove";
  toast.dataset.action = "ui.element.remove";

  const inner = document.createElement("div");
  const p = document.createElement("p");
  p.innerText = msg;
  inner.appendChild(p);

  toast.innerHTML = makeSVG(level, 16);
  toast.append(inner);

  document.getElementById("toast-root").appendChild(toast);
}
