/** State used for the slider drag */
let dragState = {
  active: false,
  thumb: null,
  root: null,
};

/**
 * Handles a click on the slider track: picks the thumb closest to the click
 * position (double sliders), snaps it there and starts dragging it.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The slider element.
 * @param {MouseEvent|TouchEvent} param0.evt - The pointer/touch event.
 */
export function sliderClick({ ele, evt }) {
  const clientX = evt.touches ? evt.touches[0].clientX : evt.clientX;
  const slider = ele.closest(".slider");
  const range = Number(slider.dataset.range);
  const sliderMin = Number(slider.dataset.min) || 0;
  const thumbs = slider.getElementsByClassName("slider__thumb");
  const isDouble = thumbs.length === 2;

  let thumb = thumbs[0];
  // get value and update input
  if (isDouble) {
    const rect = slider.getBoundingClientRect();
    const frac = Math.max(0, Math.min(1, (clientX - rect.left) / rect.width));
    const actualValue = sliderMin + Math.round(frac * range);

    // find the closest to value
    const loInput = slider.querySelector(`[name='${thumbs[0].dataset.thumb}']`);
    let lo;
    if (loInput.value === "") {
      lo = Number(loInput.dataset.default);
    } else {
      lo = Number(loInput.value);
    }

    const hiInput = slider.querySelector(`[name='${thumbs[1].dataset.thumb}']`);
    let hi;
    if (hiInput.value === "") {
      hi = Number(hiInput.dataset.default);
    } else {
      hi = Number(hiInput.value);
    }

    if (
      actualValue > hi ||
      Math.abs(hi - actualValue) < Math.abs(lo - actualValue)
    ) {
      thumb = thumbs[1];
    }
  }
  updateSlider(clientX, slider, thumb);
  startSliderDrag({ ele: thumb });
}

/**
 * Starts dragging a thumb: marks it active and wires the document
 * pointermove/pointerup listeners.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The thumb element.
 */
export function startSliderDrag({ ele }) {
  ele.classList.add("active");
  dragState.active = true;
  dragState.thumb = ele;
  dragState.root = ele.closest(".slider");
  document.addEventListener("pointermove", moveSliderDrag);
  document.addEventListener("pointerup", endSliderDrag);
}

/**
 * Moves the active thumb with the pointer.
 *
 * @param {MouseEvent|TouchEvent} evt - The pointer/touch event.
 */
function moveSliderDrag(evt) {
  if (!dragState.active) return;
  const clientX = evt.touches ? evt.touches[0].clientX : evt.clientX;
  updateSlider(clientX, dragState.root, dragState.thumb);
}

/**
 * Ends the drag, clears the drag state and unwires the document listeners.
 */
function endSliderDrag() {
  dragState.thumb.classList.remove("active");
  dragState.active = false;
  dragState.thumb = null;
  dragState.root = null;
  document.removeEventListener("pointermove", moveSliderDrag);
  document.removeEventListener("pointerup", endSliderDrag);
}

/**
 * Moves a thumb to the pointer position, clamps it against the other thumb
 * (double sliders), updates the hidden input and repaints the thumb/fill.
 *
 * @param {number} clientX - Pointer X coordinate.
 * @param {Element} slider - The slider root element.
 * @param {Element} thumb - The thumb being moved.
 */
function updateSlider(clientX, slider, thumb) {
  const range = Number(slider.dataset.range);
  const sliderMin = Number(slider.dataset.min) || 0;
  const thumbs = slider.getElementsByClassName("slider__thumb");
  const isDouble = thumbs.length === 2;
  const fill = slider.getElementsByClassName("slider__fill")[0];
  const thumbName = thumb.dataset.thumb;
  const input = slider.querySelector(`[name='${thumbName}']`);

  // get value and update input
  const rect = slider.getBoundingClientRect();
  const frac = Math.max(0, Math.min(1, (clientX - rect.left) / rect.width));
  let rawValue = Math.round(frac * range);
  let actualValue = sliderMin + rawValue;

  if (!isDouble) {
    input.value = actualValue;

    const pct = (rawValue / range) * 100;
    thumb.style.left = pct + "%";
    fill.style.width = pct + "%";
  } else {
    const loInput = slider.querySelector(`[name='${thumbs[0].dataset.thumb}']`);
    let lo;
    if (loInput.value === "") {
      lo = Number(loInput.dataset.default);
    } else {
      lo = Number(loInput.value);
    }

    const hiInput = slider.querySelector(`[name='${thumbs[1].dataset.thumb}']`);
    let hi;
    if (hiInput.value === "") {
      hi = Number(hiInput.dataset.default);
    } else {
      hi = Number(hiInput.value);
    }

    let dragging = thumbs[0].dataset.thumb === thumbName ? "lo" : "hi";
    let loPct = Number(thumbs[0].style.left.slice(0, -1));
    let hiPct = Number(thumbs[1].style.left.slice(0, -1));
    let pct;

    if (dragging === "lo") {
      actualValue = Math.min(actualValue, hi - 1);
      rawValue = actualValue - sliderMin;
      pct = (rawValue / range) * 100;
      loPct = pct;
    }
    if (dragging === "hi") {
      actualValue = Math.max(actualValue, lo + 1);
      rawValue = actualValue - sliderMin;
      pct = (rawValue / range) * 100;
      hiPct = pct;
    }
    if (Number(input.dataset.default) === actualValue) {
      actualValue = "";
    }
    input.value = actualValue;
    input.dispatchEvent(new Event("input", { bubbles: true }));

    thumb.style.left = pct + "%";

    fill.style.left = loPct + "%";
    fill.style.width = `${hiPct - loPct}%`;
  }
}

/**
 * Resets the slider visuals back to the full range.
 *
 * @param {Element} slider - The slider root element.
 */
export function resetSlider(slider) {
  const fill = slider.querySelector(".slider__fill");
  fill.style.left = "0%";
  fill.style.width = "100%";

  const thumbs = slider.querySelectorAll(".slider__thumb");
  if (thumbs.length == 2) {
    thumbs[0].style.left = "0%";
    thumbs[1].style.left = "100%";
  } else {
    thumbs[0].style.left = "100%";
  }
}
