let dragState = {
  active: false,
  thumb: null,
  root: null,
};

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

    if (actualValue > hi || Math.abs(hi - actualValue) < Math.abs(lo - actualValue)) {
      thumb = thumbs[1];
    }
  }
  updateSlider(clientX, slider, thumb);
  startSliderDrag({ ele: thumb });
}

export function startSliderDrag({ ele }) {
  ele.classList.add("active");
  dragState.active = true;
  dragState.thumb = ele;
  dragState.root = ele.closest(".slider");
  document.addEventListener("pointermove", moveSliderDrag);
  document.addEventListener("pointerup", endSliderDrag);
}

function moveSliderDrag(evt) {
  if (!dragState.active) return;
  const clientX = evt.touches ? evt.touches[0].clientX : evt.clientX;
  updateSlider(clientX, dragState.root, dragState.thumb);
}

function endSliderDrag() {
  dragState.thumb.classList.remove("active");
  dragState.active = false;
  dragState.thumb = null;
  dragState.root = null;
  document.removeEventListener("pointermove", moveSliderDrag);
  document.removeEventListener("pointerup", endSliderDrag);
}

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
    let dragging;

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

    for (let i = 0; i < thumbs.length; i++) {
      if (thumbs[i].dataset.thumb === thumbName) {
        dragging = i === 0 ? "lo" : "hi";
      }
    }

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
