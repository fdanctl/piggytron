import { handleInputOnBlur } from "../input";
import { showToast } from "../toast";
import { goalActions } from "./goal";
import { uiActions } from "./ui";

/** Observes data-intersect elements and dispatches their action on entry. */
const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      const ele = entry.target;
      dispatch(ele.dataset[eventAttributes.intersect], {
        ele,
        evt: entry,
        data: ele.dataset,
      });
    }
  });
});

// Re-observes data-intersect elements after each htmx response, since the
// DOM may have been swapped.
document.body.addEventListener("htmx:afterOnLoad", function () {
  observer.disconnect();
  const elements = document.querySelectorAll(
    `[data-${eventAttributes.intersect}]`,
  );
  elements.forEach((el) => observer.observe(el));
});

/**
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The triggering element.
 * @param {Event} param0.evt - The triggering event.
 * @param {DOMStringMap} param0.data - The dataset of the element.
 */
function log({ ele, evt, data }) {
  console.log(ele);
  console.log(evt);
  console.log(data);
}

/** Registry of every known action. ([k]: v, [action-name]: function) */
const actions = {
  "dev.log.this": log,
  ...uiActions,
  ...goalActions,
};

/**
 * Maps DOM event names to the data-* attribute that carries the action
 * name for that event (e.g. click -> data-action, input -> data-input).
 */
const eventAttributes = {
  click: "action",
  input: "input",
  change: "change",
  focusin: "focus",
  focusout: "blur",
  keydown: "keydown",
  pointerdown: "pointerdown",
  "htmx:beforeRequest": "beforerequest",
  "htmx:afterRequest": "afterrequest",
  "htmx:afterOnLoad": "afteronload",
  animationend: "animationend",
  submit: "submit",
  intersect: "intersect",
  mouseover: "mouseover",
  mouseout: "mouseout",
};

// Checks and marks inputs invalid on blur.
document.addEventListener("focusout", (evt) => {
  if (!(evt.target instanceof HTMLInputElement)) return;

  handleInputOnBlur(evt.target);
});

window.addEventListener("offline", () => {
  showToast("error", "You are offline");
});

window.addEventListener("online", () => {
  showToast("success", "Internet connection restored");
});

/**
 * Runs every space-separated action name in an action string.
 *
 * @param {string} actionName - Space-separated action names.
 * @param {Object} payload - The payload passed to each action.
 * @param {Element} payload.ele - The triggering element.
 * @param {Event} payload.evt - The triggering event.
 * @param {DOMStringMap} payload.data - The dataset of the element.
 */
function dispatch(actionName, payload) {
  const names = actionName.trim().split(/\s+/);

  for (const n of names) {
    const action = actions[n];
    if (!action) {
      console.warn(`Unknown action: ${n}`);

      continue;
    }

    action(payload);
  }
}

/**
 * Creates a delegated listener that dispatches the action stored in the
 * data-* attribute of the closest matching element.
 *
 * @param {string} eventName - The DOM event name.
 * @param {string} dataAttr - The data-* attribute name (without the prefix).
 */
function createListener(eventName, dataAttr) {
  document.addEventListener(eventName, (evt) => {
    const ele = evt.target.closest(`[data-${dataAttr}]`);

    if (!ele) return;

    dispatch(ele.dataset[dataAttr], {
      ele,
      evt,
      data: ele.dataset,
    });
  });
}

// Wires one listener per event type.
// For data-intersect, it starts 'observe' to all elements with data-intersect.
for (let [key, value] of Object.entries(eventAttributes)) {
  if (value === eventAttributes.intersect) {
    const elements = document.querySelectorAll(`[data-${value}]`);
    elements.forEach((el) => observer.observe(el));
  } else {
    createListener(key, value);
  }
}
