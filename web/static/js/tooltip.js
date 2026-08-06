let tid = null;

export function showPopover({ data }) {
  if (tid) return;

  const delay = !Number.isNaN(Number(data.delay)) ? Number(data.delay) : 0;
  const popover = document.getElementById(data.name + "-popover");
  if (popover.matches(":popover-open")) return;

  tid = setTimeout(() => {
    popover.showPopover();
  }, delay);
}

export function hidePopover({ ele, data, evt }) {
  if (!tid) return;
  const popover = document.getElementById(data.name + "-popover");
  const toElement = evt.toElement;
  if (ele.contains(toElement) || popover.contains(toElement)) {
    return;
  }

  clearTimeout(tid);
  document.getElementById(data.name + "-popover").hidePopover();
  tid = null;
}
