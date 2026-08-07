import {
  buildCalendar,
  generateYearLI,
  nextMonth,
  prevMonth,
} from "../calendar";
import {
  changeMaxAmountChip,
  changeMaxDateChip,
  changeMinAmountChip,
  changeMinDateChip,
  filterAccordionToggle,
  removeFilterPill,
  resetTransactionFiltersForm,
  toggleFilterPill,
} from "../filters-sheet";
import { disableBtns } from "../form";
import { sliderClick, startSliderDrag } from "../slider";
import {
  cashInputBlur,
  cashInputFocus,
  budgetInput,
  checkboxPillToggle,
  dateOnChange,
  openCalendar,
  openTimePopover,
  passwordToggle,
  sanitizeCashInput,
  selectDay,
  selectSelect,
  selectTime,
  centerSelected,
  timeOnChange,
  budgetInputScan,
  selectPillSelect,
} from "../input";
import {
  closeLastDialog,
  colapseSidebar,
  collapseAllSublinks,
  collapseSublinks,
  expandSublinks,
  handleSublinkFocus,
  navigate,
  openNavSheet,
  overlayClose,
  sidebarHidePopover,
  sidebarShowPopover,
  startDialogDrag,
  toggleSublinks,
} from "../navigation";
import { changeTab } from "../tabs";
import {
  generateMoreMonthsLI,
  nextBudgetMonth,
  prevBudgetMonth,
} from "../month-pill";
import { hidePopover, showPopover } from "../tooltip";

function queryNextSelector(ele, selector) {
  const results = ele.parentElement.querySelectorAll(selector);
  for (let i = 0; i < results.length; i++) {
    const elt = results[i];
    if (elt.compareDocumentPosition(ele) === Node.DOCUMENT_POSITION_PRECEDING) {
      return elt;
    }
  }
}

function queryPreviousSelector(ele, selector) {
  const results = ele.parentElement.querySelectorAll(selector);
  for (let i = results.length - 1; i >= 0; i--) {
    const elt = results[i];
    if (elt.compareDocumentPosition(ele) === Node.DOCUMENT_POSITION_FOLLOWING) {
      return elt;
    }
  }
}
function getTarget(ele, selector) {
  if (selector.indexOf("closest ") === 0) {
    return ele.closest(selector.slice(8));
  } else if (selector.indexOf("find ") === 0) {
    return ele.querySelector(selector.slice(5));
  } else if (selector === "next" || selector === "nextElementSibling") {
    return ele.nextElementSibling;
  } else if (selector.indexOf("next ") === 0) {
    return queryNextSelector(ele, selector.slice(5));
  } else if (selector === "previous" || selector === "previousElementSibling") {
    return ele.previousElementSibling;
  } else if (selector.indexOf("previous ") === 0) {
    return queryPreviousSelector(ele, selector.slice(9));
  } else {
    return document.querySelector(selector);
  }
}

function removeEle({ ele, data }) {
  let t = ele;
  if (data.target) {
    t = getTarget(ele, data.target);
  }
  t.remove();
}

function clickEle({ ele, data }) {
  let t = ele;
  if (data.target) {
    t = getTarget(ele, data.target);
  }
  t.click();
}

export const uiActions = {
  "ui.element.remove": removeEle,

  "ui.element.click": clickEle,

  "ui.calendar.prev-month": prevMonth,

  "ui.calendar.next-month": nextMonth,

  "ui.calendar.rebuild": buildCalendar,

  "ui.calendar.generate-year-options": generateYearLI,

  "ui.amount-chip.update-min": changeMinAmountChip,

  "ui.amount-chip.update-max": changeMaxAmountChip,

  "ui.date-chip.update-min": changeMinDateChip,

  "ui.date-chip.update-max": changeMaxDateChip,

  "ui.filters.remove": removeFilterPill,

  "ui.filters-accordion.toggle": filterAccordionToggle,

  "ui.filters.reset": resetTransactionFiltersForm,

  "ui.filters.toggle": toggleFilterPill,

  "ui.password.toggle": passwordToggle,

  "ui.checkbox-pill.toggle": checkboxPillToggle,

  "ui.select.option": selectSelect,

  "ui.select-pill.option": selectPillSelect,

  "ui.date-input.toggle": openCalendar,

  "ui.date-input.select": selectDay,

  "ui.time-input.toggle": openTimePopover,

  "ui.time-input.select": selectTime,

  "ui.select.center-selected": centerSelected,

  "ui.budget-input.keydown": budgetInput,

  "ui.budget-input.animate": budgetInputScan,

  "ui.cash-input.input": sanitizeCashInput,

  "ui.cash-input.focus": cashInputFocus,

  "ui.cash-input.blur": cashInputBlur,

  "ui.date-input.input": dateOnChange,

  "ui.time-input.input": timeOnChange,

  "ui.sidebar.colapse": colapseSidebar,

  "ui.nav-sheet.open": openNavSheet,

  "ui.dialog.close-last": closeLastDialog,

  "ui.dialog.start-drag": startDialogDrag,

  "ui.nav.navigate": navigate,

  "ui.form.disable-btns": disableBtns,

  "ui.tab.change": changeTab,

  "ui.slider.drag-start": startSliderDrag,

  "ui.slider.click": sliderClick,

  "ui.overlay.close": overlayClose,

  "ui.month-pill.generate-options": generateMoreMonthsLI,

  "ui.month-pill.prev-budget": prevBudgetMonth,

  "ui.month-pill.next-budget": nextBudgetMonth,

  "ui.sidebar.show-popover": sidebarShowPopover,

  "ui.sidebar.hide-popover": sidebarHidePopover,

  "ui.nav.toggle-sublinks": toggleSublinks,

  "ui.nav.expand-sublinks": expandSublinks,

  "ui.nav.collapse-sublinks": collapseSublinks,

  "ui.nav.collapse-all-sublinks": collapseAllSublinks,

  "ui.nav.sublinks-focus": handleSublinkFocus,

  "ui.popover.show": showPopover,

  "ui.popover.hide": hidePopover,
};
