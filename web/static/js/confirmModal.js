import { makeSVG } from "./icons";
import { closeDialog } from "./navigation";

/**
 * Opens a confirmation dialog and resolves with the user's choice.
 *
 * @param {Object} [options] - Dialog options.
 * @param {string} [options.title] - Dialog title (default "Confirm").
 * @param {string} [options.message] - Message body (default "Are you sure?").
 * @param {string} [options.acceptText] - Accept button label (default "Yes").
 * @param {string} [options.refuseText] - Refuse button label (default "No").
 * @returns {Promise<boolean>} Resolves true if accepted, false if refused.
 */
export function confirmModal({
  title = "Confirm",
  message = "Are you sure?",
  acceptText = "Yes",
  refuseText = "No",
} = {}) {
  return new Promise((resolve) => {
    const overlay = document.createElement("div");
    overlay.classList.add("dialog__overlay");
    overlay.tabIndex = "-1";
    overlay.dataset.pointerdown = "ui.dialog.start-drag";

    const modal = document.createElement("div");
    modal.tabIndex = "-1";
    modal.classList.add("dialog", "dialog--float");

    modal.innerHTML = `
		<div class="dialog__bar">
			<div></div>
		</div>
		<button type="button" class="reset-btn dialog__x" data-action="ui.dialog.close-last">
      ${makeSVG("x", 26)}
		</button>
    <div>
      <h4 class="mb-md">${title}</h4>
      <p class="text-subtitle">${message}</p>
      <div class="flex gap-xs justify-end items-center mt-sm">
        <button class="btn btn--outline refuse">${refuseText}</button>
        <button class="btn btn--outline accept">${acceptText}</button>
      </div>
`;

    overlay.appendChild(modal);
    const root = document.getElementById("dialog-root");
    root.appendChild(overlay);
    root.classList.add("is-dialog-open");
    overlay.focus();

    modal.querySelector(".accept").onclick = () => {
      closeDialog(modal);
      resolve(true);
    };

    modal.querySelector(".refuse").onclick = () => {
      closeDialog(modal);
      resolve(false);
    };
  });
}
