import { buildCalendar, nextMonth, prevMonth } from "../calendar";
import {
  filterAccordionToggle,
  removeFilterPill,
  resetTransactionFiltersForm,
  toggleFilterPill,
} from "../filters-sheet";
import { disableBtns } from "../form";
import {
  cashInputBlur,
  cashInputFocus,
  cashInputNav,
  checkboxPillToggle,
  dateOnChange,
  openCalendar,
  openTimePopover,
  passwordToggle,
  sanitizeCashInput,
  selectDay,
  selectSelect,
  selectTime,
  selectToggle,
  timeOnChange,
} from "../input";
import {
  closeLastDialog,
  colapseSidebar,
  navigate,
  openNavSheet,
} from "../navigation";

export const uiActions = {
  "ui.calendar.prev-month": prevMonth,

  "ui.calendar.next-month": nextMonth,

  "ui.calendar.rebuild": buildCalendar,

  "ui.filters.remove": removeFilterPill,

  "ui.filters-accordion.toggle": filterAccordionToggle,

  "ui.filters.reset": resetTransactionFiltersForm,

  "ui.filters.toggle": toggleFilterPill,

  "ui.password.toggle": passwordToggle,

  "ui.checkbox-pill.toggle": checkboxPillToggle,

  "ui.select.option": selectSelect,

  "ui.date-input.toggle": openCalendar,

  "ui.date-input.select": selectDay,

  "ui.time-input.toggle": openTimePopover,

  "ui.time-input.select": selectTime,

  "ui.select.toggle": selectToggle,

  "ui.cash-input.keydown": cashInputNav,

  "ui.cash-input.input": sanitizeCashInput,

  "ui.cash-input.focus": cashInputFocus,

  "ui.cash-input.blur": cashInputBlur,

  "ui.date-input.input": dateOnChange,

  "ui.time-input.input": timeOnChange,

  "ui.sidebar.colapse": colapseSidebar,

  "ui.nav-sheet.open": openNavSheet,

  "ui.dialog.close-last": closeLastDialog,

  "ui.nav.navigate": navigate,

  "ui.form.disable-btns": disableBtns,
};
