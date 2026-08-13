/** Pending timeout id for the delayed popover show. */
let tid = null;

/**
 * Shows the named popover after an optional delay (data.delay, ms).
 *
 * @param {Object} param0 - Action payload.
 * @param {DOMStringMap} param0.data - Dataset of the trigger; `data.name`
 *   identifies the popover and `data.delay` the delay in ms.
 */
export function showPopover({ data }) {
  if (tid) return;

  const delay = !Number.isNaN(Number(data.delay)) ? Number(data.delay) : 0;
  const popover = document.getElementById(data.name + "-popover");
  if (popover.matches(":popover-open")) return;

  tid = setTimeout(() => {
    popover.showPopover();
  }, delay);
}

/**
 * Hides the named popover, unless the pointer moved into the trigger or
 * the popover itself.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The trigger element.
 * @param {DOMStringMap} param0.data - Dataset of the trigger (`data.name`).
 * @param {Event} param0.evt - The mouseout event.
 */
export function hidePopover({ ele, data, evt }) {
  if (!tid) return;
  const popover = document.getElementById(data.name + "-popover");

  // .toElement it's the chrome way and
  // .explicitOriginalTarget it's the firefox way
  const toElement = evt.toElement || evt.explicitOriginalTarget;
  if (ele.contains(toElement) || popover.contains(toElement)) {
    return;
  }

  clearTimeout(tid);
  document.getElementById(data.name + "-popover").hidePopover();
  tid = null;
}
