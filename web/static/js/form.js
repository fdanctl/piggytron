/**
 * Disables every .btn inside the given element (used while a form is
 * submitting, to prevent double submissions).
 * TODO: this could be replaced with hx-disable (v4), try.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The element whose buttons should be disabled.
 */
export function disableBtns({ ele }) {
  const btns = ele.querySelectorAll(".btn");
  for (let i = 0; i < btns.length; i++) {
    btns[i].disabled = true;
  }
}
