function accordionToggle(ele) {
  const div = ele.parentElement.parentElement.children[1];
  div.classList.toggle("flex-wrap");
  ele.children[0].classList.toggle("hidden");
  ele.children[1].classList.toggle("hidden");
}
